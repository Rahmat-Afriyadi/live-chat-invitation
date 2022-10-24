package controller

import (
	"encoding/json"
	"fmt"
	"live-chat-gorilla/dto"
	"live-chat-gorilla/helper"
	"live-chat-gorilla/service"
	"net/http"
)

type CommentController struct {
	*service.CommentService
}

func NewCommentController(s *service.CommentService) *CommentController {
	return &CommentController{
		s,
	}
}

func (c *CommentController) Store(w http.ResponseWriter, r *http.Request) {
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

	final := map[string]interface{}{
		"status":  status,
		"comment": result,
	}

	helper.RespondWithJSON(w, http.StatusCreated, final)
	return
}

func (c *CommentController) Index(w http.ResponseWriter, r *http.Request) {
	result := c.CommentService.All(r)
	helper.RespondWithJSON(w, http.StatusAccepted, result)
	return
}

func (c *CommentController) Status(w http.ResponseWriter, r *http.Request) {
	status := c.CommentService.GetStatus()
	helper.RespondWithJSON(w, http.StatusAccepted, status)
	return
}
