package service

import (
	"live-chat-gorilla/dto"
	"live-chat-gorilla/entity"
	"live-chat-gorilla/repository"
	"net/http"
)

type CommentService struct {
	*repository.CommentRepository
	*repository.PresenceRepository
}

func NewCommentService(r1 *repository.CommentRepository,r2 *repository.PresenceRepository) *CommentService {
	return &CommentService{
		r1,
		r2,
	}
}

func (s *CommentService) Create(dto dto.CommentInsertDTO) entity.Comment {
	presenceEntity := entity.Presence{
				Name:   dto.Name,
				Status: dto.Status,
			}
	presenceResult := s.PresenceRepository.Insert(presenceEntity)

	commentEntity := entity.Comment{
		PresenceId: presenceResult.Id,
		Comment:    dto.Comment,
	}

	commentResult := s.CommentRepository.Insert(commentEntity)

	return commentResult
}

func (s *CommentService) All(r *http.Request) []entity.Comment {
	return s.CommentRepository.List(r)
}

func (s *CommentService) GetStatus() dto.Status {
	return s.PresenceRepository.GetStatus()
}
