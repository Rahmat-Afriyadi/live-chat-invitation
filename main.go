package main

import (
	"live-chat-gorilla/config"
	"live-chat-gorilla/helper"

	"live-chat-gorilla/controller"
	"live-chat-gorilla/repository"
	"live-chat-gorilla/service"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type PresenceComment struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Status  int `json:"status"`
}

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	commentRepository *repository.CommentRepository = repository.NewCommentRepository(db)
	presenceRepository *repository.PresenceRepository = repository.NewPresenceRepository(db)
	commentService    *service.CommentService       = service.NewCommentService(commentRepository, presenceRepository)
	commentController *controller.CommentController = controller.NewCommentController(commentService)

)


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
	router.HandleFunc("/store", commentController.Store).Methods("POST")
	router.HandleFunc("/paginate", commentController.Index)
	router.HandleFunc("/status", commentController.Status)

	c := cors.New(cors.Options{
        AllowedOrigins: []string{"https://syafiq-ina.vercel.app"},
        AllowCredentials: true,
    })

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe( ":" + os.Getenv("PORT") , handler))

}
