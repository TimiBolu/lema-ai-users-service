package models

import "time"

type Post struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"not null;index" json:"userId"`
	Title     string    `gorm:"not null" json:"title"`
	Body      string    `gorm:"not null" json:"body"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
