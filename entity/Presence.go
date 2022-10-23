package entity

import "time"

type Presence struct {
	Id       uint      `gorm:"primaryKey:autoIncreament" json:"id"`
	Name     string    `json:"name"`
	Status   int       `json:"status"`
	Comments []Comment `json:"comments"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
