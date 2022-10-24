package main

import (
	"encoding/json"
	"fmt"
	"live-chat-gorilla/config"
	"live-chat-gorilla/helper"

	"live-chat-gorilla/controller"
	"live-chat-gorilla/entity"
	"live-chat-gorilla/repository"
	"live-chat-gorilla/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"os"
	"github.com/joho/godotenv"
)

type PresenceComment struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Status  int `json:"status"`
}

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsConn *websocket.Conn

	commentRepository *repository.CommentRepository = repository.NewCommentRepository(db)
	presenceRepository *repository.PresenceRepository = repository.NewPresenceRepository(db)
	commentService    *service.CommentService       = service.NewCommentService(commentRepository, presenceRepository)
	commentController *controller.CommentController = controller.NewCommentController(commentService)

)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {

	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		// check the http.Request
		// make sure it's OK to access
		return true
	}
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}

	defer wsConn.Close()

	// event loop
	for {
		var msg PresenceComment
		
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error reading JSON: %s\n", err.Error())
			break
		}
		fmt.Printf("ini datanya %s\n", msg.Name)

		presenceEntity := entity.Presence{
			Name:   msg.Name,
			Status: msg.Status,
		}
		presenceResult := presenceRepository.Insert(presenceEntity)

		commentEntity := entity.Comment{
			PresenceId: presenceResult.Id,
			Comment:    msg.Comment,
		}

		commentResult := commentRepository.Insert(commentEntity)

		b, err := json.Marshal(commentResult)
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println(string(b))

		// fmt.Printf("Message Received: %s\n", msg.Comment)
		err = wsConn.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			fmt.Printf("error sending message: %s\n", err.Error())
		}
	}
}


func main() {
	defer config.CloseDatabaseConnection(db)

	err := godotenv.Load()

    if err != nil {
        log.Fatal("Error loading .env file")
    }

	router := mux.NewRouter()

	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./assets/")))
	router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
		helper.RespondWithJSON(w, http.StatusAccepted, os.Getenv("PORT"))
	})
	router.HandleFunc("/socket", WsEndpoint)
	router.HandleFunc("/store", commentController.Store).Methods("POST")
	router.HandleFunc("/paginate", commentController.Index)

	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowCredentials: true,
    })

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe( ":" + os.Getenv("PORT") , handler))

}
