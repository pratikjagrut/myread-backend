package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Status      string `gorm:"not null" json:"status"`
	Image       string `json:"image"`
	Author      string `gorm:"not null" json:"author"`
	Description string `json:"description"`
	UserID      uint
}
