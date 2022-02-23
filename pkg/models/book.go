package models

import (
	"errors"
	"html"
	"strings"

	"github.com/badoux/checkmail"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `gorm:"type:varchar(191);primaryKey;not null" json:"name"`
	Status      string `gorm:"type:varchar(191);not null" json:"status"`
	Image       []byte `gorm:"type:blob" json:"image"`
	Author      string `gorm:"type:varchar(191);not null" json:"author"`
	Description string `json:"description"`
	UserEmail   string `gorm:"type:varchar(191);primaryKey;not null" json:"user_email"`
}

type BookStatus string

var (
	Reading  BookStatus = "reading"
	Finished BookStatus = "finished"
	Wishlist BookStatus = "wishlist"
	Status              = map[BookStatus]bool{
		Reading:  true,
		Finished: true,
		Wishlist: true,
	}
)

func (b *Book) SaveBook(db *gorm.DB) error {
	sanitizeBook(b)
	return db.Create(&b).Error
}

func ListAllBooks(db *gorm.DB, email string) (books *[]Book, err error) {
	err = db.Where("user_email = ?", email).Find(&books).Error
	return books, err
}

func ListBooksByStatus(db *gorm.DB, status, email string) (books *[]Book, err error) {
	err = db.Where("user_email = ? and status = ?", email, status).Find(&books).Error
	return books, err
}

func GetBookByName(db *gorm.DB, name, email string) (book *Book, err error) {
	err = db.Where("user_email = ? and name = ?", email, name).First(&book).Error
	return book, err
}

func UpdateBook(db *gorm.DB, name, status, email string) error {
	var book Book
	return db.Model(&book).Where("user_email = ? and name = ?", email, name).Update("status", status).Error
}

func DeleteBook(db *gorm.DB, name, email string) error {
	var book Book
	// delete matched records permanently with Unscoped()
	return db.Where("user_email = ? and name = ?", email, name).Unscoped().Delete(&book).Error
}

func sanitizeBook(b *Book) {
	b.Name = html.EscapeString(strings.TrimSpace(b.Name))
	b.Status = html.EscapeString(strings.TrimSpace(b.Status))
	b.Author = html.EscapeString(strings.TrimSpace(b.Author))
	b.Description = html.EscapeString(strings.TrimSpace(b.Description))
}

func ValidateBook(b *Book, op string) error {
	switch op {
	case "create":
		if b.Name == "" {
			return errors.New("Required Book Name")
		}
		if b.Author == "" {
			return errors.New("Required Author Name")
		}
		if b.Status == "" {
			return errors.New("Required on of the book status(reading, finished, wishlist)")
		} else if !Status[BookStatus(b.Status)] {
			return errors.New("book status must be from (reading, finished, wishlist)")
		}
		if err := checkmail.ValidateFormat(b.UserEmail); err != nil {
			return errors.New("Invalid Email")
		}
	case "update":
		if b.Name == "" {
			return errors.New("Required Book Name")
		}
		if b.Status == "" {
			return errors.New("Required on of the book status(reading, finished, wishlist)")
		} else if !Status[BookStatus(b.Status)] {
			return errors.New("book status must be from (reading, finished, wishlist)")
		}
	case "delete":
		if b.Name == "" {
			return errors.New("Required Book Name")
		}
	}
	return nil
}
