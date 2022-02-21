package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"primaryKey;not null;unique" json:"email"`
	Password string `json:"-"`
	Books    []Book
}
