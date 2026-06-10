package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);unique" binding:"required"`
	Password string `gorm:"not null" binding:"required"`
	Tasks    []Task
}
