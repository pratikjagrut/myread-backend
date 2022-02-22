package models

import (
	"errors"
	"html"
	"strings"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"primaryKey;not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Books    []Book
}

func (u *User) SaveUser(db *gorm.DB) error {
	b, err := hashedPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(b)
	sanitize(u)
	return db.Create(&u).Error
}

func FindUserByEmail(db *gorm.DB, email string) (user *User, err error) {
	err = db.Where("email = ?", email).First(&user).Error
	return user, err
}

func hashedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func sanitize(u *User) {
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func Validate(action string, u *User) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "create":
		if u.Name == "" {
			return errors.New("Required Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	default:
		return errors.New("Wrong operation")
	}
}

// func VerifyPassword(hashedPassword, password string) error {
// 	return
// }
