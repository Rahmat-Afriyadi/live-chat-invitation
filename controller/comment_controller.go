package controller

import (
	"encoding/json"
	"fmt"
	"live-chat-gorilla/dto"
	"live-chat-gorilla/helper"
	"live-chat-gorilla/service"
	"net/http"

	"github.com/pusher/pusher-http-go"
)

type CommentController struct {
	*service.CommentService	
}

func NewCommentController(s *service.CommentService) *CommentController {
	return &CommentController{
		s,
	}
}

func (c *CommentController) Store(w http.ResponseWriter,r *http.Request) {
	pusherClient := pusher.Client{
		AppID:   "1492360",
		Key:     "bb38092e1dac94cc76b4",
		Secret:  "9ddf1bbd4494afee3e86",
		Cluster: "ap1",
		Secure:  true,
	}	

	var dto dto.CommentInsertDTO
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&dto); err != nil {
        helper.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
	fmt.Println("test")			

	result := c.CommentService.Create(dto)
	status := c.CommentService.GetStatus()

	// var data map[string]interface{}
	newComment, _ := json.Marshal(result)
	currentStatus, _ := json.Marshal(status)
    // json.Unmarshal(inrec, &data)
	pusherClient.Trigger("live_chat_invitation", "message", string(newComment))
	pusherClient.Trigger("live_chat_invitation", "status", string(currentStatus))

	helper.RespondWithJSON(w, http.StatusCreated, result)
	return	
}

func (c *CommentController) Index (w http.ResponseWriter,r *http.Request) {
	result := c.CommentService.All(r)
	helper.RespondWithJSON(w, http.StatusAccepted, result)
	return
}