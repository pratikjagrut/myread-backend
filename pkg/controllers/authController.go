package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-backend/pkg/database"
	"github.com/pratikjagrut/myreads-backend/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const SecretKey = "secret"

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		e := fmt.Sprintf("CreateUser: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := models.Validate("create", &user); err != nil {
		e := fmt.Sprintf("CreateUser: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := user.SaveUser(database.Database.Db); err != nil {
		var mysqlErr *mysql.MySQLError
		m := ""
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			m = "This email id is already registered"
		} else {
			m = "User registration failed"
		}

		database.Database.Db.Logger.Error(context.Background(), err.Error())
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": m,
			"status":  fiber.StatusBadRequest,
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "User registration successful",
		"status":  fiber.StatusOK,
	})
}

type ResponseUser struct {
	Name  string
	Email string
}

func getResponseUser(user *models.User) *ResponseUser {
	return &ResponseUser{
		Name:  user.Name,
		Email: user.Email,
	}
}

func Login(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		e := fmt.Sprintf("Login: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := models.Validate("login", &user); err != nil {
		e := fmt.Sprintf("Login: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	foundUser, err := models.FindUserByEmail(database.Database.Db, user.Email)
	if err != nil {
		m := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			m = "User Not Found"
		} else {
			m = "something went wrong"
		}
		e := fmt.Sprintf("Login: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": m,
			"status":  fiber.StatusNotFound,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		log.Println("Login: ", err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Incorrect credentials",
			"status":  fiber.StatusUnauthorized,
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    foundUser.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		log.Println("Login: ", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  fiber.StatusInternalServerError,
		})
	}

	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(cookie)
	return c.JSON(fiber.Map{
		"message": "Login Success",
		"status":  fiber.StatusOK,
		"user":    getResponseUser(foundUser),
	})
}

func getIssuer(c *fiber.Ctx) (string, error) {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		c.Status(fiber.StatusUnauthorized)
		return "", fmt.Errorf("getIssuer: empty cookie")
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return "", fmt.Errorf("getIssuer: %v", err)
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer, nil

}

func GetUser(c *fiber.Ctx) error {
	issuer, err := getIssuer(c)
	if err != nil {
		e := fmt.Sprintf("GetUser: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	foundUser, err := models.FindUserByEmail(database.Database.Db, issuer)
	if err != nil {
		m := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			m = "User Not Found"
		} else {
			m = "something went wrong"
		}
		e := fmt.Sprintf("GetUser: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": m,
			"status":  fiber.StatusNotFound,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Fetch user successful",
		"status":  fiber.StatusOK,
		"user":    getResponseUser(foundUser),
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Logout Success.",
		"status":  fiber.StatusOK,
	})
}
