package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-backend/pkg/database"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func main() {
	database.ConnectDb()
	app := fiber.New()

	app.Get("/", welcome)

	log.Fatal(app.Listen(":3000"))
}
