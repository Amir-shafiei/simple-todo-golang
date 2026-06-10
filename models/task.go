package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
	UserID uint   `json:"user_id"`
}
