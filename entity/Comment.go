package entity

import "time"

type Comment struct {
	Id         uint     `gorm:"primaryKey:autoIncreament" json:"id"`
	Comment    string   `json:"comment"`
	PresenceId uint     `json:"-"`
	Presence   Presence `gorm:"foreign:PresenceId;constraint:onDelete:CASCADE,onUpdate:CASCADE" json:"presence"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
