package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-backend/pkg/database"
	"github.com/pratikjagrut/myreads-backend/pkg/models"
)

func AddBook(c *fiber.Ctx) error {
	issuer, err := getIssuer(c)
	if err != nil {
		e := fmt.Sprintf("BookEntry: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  fiber.StatusUnauthorized,
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println("ERROR: BookEntry: Form: ", err)
		return err
	}

	if err := c.BodyParser(form); err != nil {
		e := fmt.Sprintf("AddBook: BodyParser: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	book := &models.Book{
		Author:      form.Value["author"][0],
		Name:        form.Value["name"][0],
		Status:      form.Value["status"][0],
		Description: form.Value["description"][0],
	}

	book.UserEmail = issuer

	if err := models.ValidateBook(book, "create"); err != nil {
		e := fmt.Sprintf("AddBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
			"status":  fiber.StatusBadRequest,
		})
	}

	foundBook, err := models.GetBookByName(database.Database.Db, book.Name, issuer)
	if err != nil && err.Error() != "record not found" {
		e := fmt.Sprintf("AddBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	if foundBook != nil && foundBook.Name == book.Name {
		c.Status(fiber.StatusFound)
		return c.JSON(fiber.Map{
			"message": "This book is present in your bookshelf",
			"status":  fiber.StatusFound,
		})
	}

	if err := book.SaveBook(database.Database.Db); err != nil {
		var mysqlErr *mysql.MySQLError
		m := ""
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			m = "This book is present in your bookshelf"
		} else {
			m = "Failed to add book"
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
		"message": fmt.Sprintf(`Book "%s" is successfully added to "%s"`, book.Name, book.Status),
		"status":  fiber.StatusOK,
	})
}

func GetBoooks(c *fiber.Ctx, which models.BookStatus) error {
	issuer, err := getIssuer(c)
	if err != nil {
		e := fmt.Sprintf("BookEntry: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  fiber.StatusUnauthorized,
		})
	}

	var books *[]models.Book

	switch status := string(which); status {
	case "":
		books, err = models.ListAllBooks(database.Database.Db, issuer)
		if err != nil {
			e := fmt.Sprintf("BookEntry: %s", err)
			database.Database.Db.Logger.Error(context.Background(), e)
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "Unauthorized",
				"status":  fiber.StatusUnauthorized,
			})
		}
	default:
		books, err = models.ListBooksByStatus(database.Database.Db, status, issuer)
		if err != nil {
			e := fmt.Sprintf("BookEntry: %s", err)
			database.Database.Db.Logger.Error(context.Background(), e)
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "Unauthorized",
				"status":  fiber.StatusUnauthorized,
			})
		}
	}

	return c.JSON(books)
}

func UpdateStatus(c *fiber.Ctx) error {
	issuer, err := getIssuer(c)
	if err != nil {
		e := fmt.Sprintf("UpdateBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  fiber.StatusUnauthorized,
		})
	}

	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		e := fmt.Sprintf("UpdateBook: BodyParser: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := models.ValidateBook(&book, "update"); err != nil {
		e := fmt.Sprintf("UpdateBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := models.UpdateBook(database.Database.Db, book.Name, book.Status, issuer); err != nil {
		e := fmt.Sprintf("UpdateBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Update status failed",
			"status":  fiber.StatusBadRequest,
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf(`Book moved to "%s" list`, book.Status),
		"status":  fiber.StatusOK,
	})
}

func DeleteBook(c *fiber.Ctx) error {
	issuer, err := getIssuer(c)
	if err != nil {
		e := fmt.Sprintf("UpdateBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  fiber.StatusUnauthorized,
		})
	}

	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		e := fmt.Sprintf("DeleteBook: BodyParser: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := models.ValidateBook(&book, "delete"); err != nil {
		e := fmt.Sprintf("DeleteBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Bad Request",
			"status":  fiber.StatusBadRequest,
		})
	}

	if err := models.DeleteBook(database.Database.Db, book.Name, issuer); err != nil {
		e := fmt.Sprintf("DeleteBook: %s", err)
		database.Database.Db.Logger.Error(context.Background(), e)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Failed to delete the book",
			"status":  fiber.StatusBadRequest,
		})
	}

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Book \"%s\" removed from your bookshelf.", book.Name),
		"status":  fiber.StatusOK,
	})
}
