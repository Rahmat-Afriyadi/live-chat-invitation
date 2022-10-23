package dto

type CommentInsertDTO struct {
	Name string `json:"name"`
	Comment string `json:"comment"`
	Status int `json:"status"`
}