package repository

import (
	"live-chat-gorilla/config"
	"live-chat-gorilla/entity"
	"net/http"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{
		db,
	}
}

func (repo *CommentRepository) Insert(e entity.Comment) entity.Comment {
	repo.db.Save(&e)
	repo.db.Model(&e).Where("comment= ?",e.Comment).Preload("Presence").First(&e)
	return e
}

func (repo *CommentRepository) List(r *http.Request) []entity.Comment {
	var comments []entity.Comment
	repo.db.Scopes(config.Paginate(r)).Preload("Presence").Find(&comments)
	return comments
}