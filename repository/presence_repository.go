package repository

import (
	"fmt"
	"live-chat-gorilla/dto"
	"live-chat-gorilla/entity"

	"gorm.io/gorm"
)

type PresenceRepository struct {
	db *gorm.DB
}

func NewPresenceRepository(db *gorm.DB) *PresenceRepository {
	return &PresenceRepository{
		db,
	}
}

func (repo *PresenceRepository) Insert(e entity.Presence) entity.Presence {
	fmt.Println("berhasil  masuk sini")
	if repo.db.Model(&e).Where("name = ?", e.Name).Updates(&e).RowsAffected == 0 {
		repo.db.Save(&e)
	}else {
		repo.db.Model(&e).Where("name = ?", e.Name).First(&e)
	}
	return e
}

func (repo *PresenceRepository) GetStatus() dto.Status {
	var status dto.Status 
	repo.db.Model(&entity.Presence{}).Where("status = ?", 99).Count(&status.NPresent)
	repo.db.Model(&entity.Presence{}).Where("status = ?", 1).Count(&status.Present)
	repo.db.Model(&entity.Presence{}).Where("status = ?", 2).Count(&status.Hesitant)

	return status
}