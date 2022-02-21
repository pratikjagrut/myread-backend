package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pratikjagrut/myreads-backend/pkg/database"
	"github.com/pratikjagrut/myreads-backend/pkg/routes"
)

func main() {
	database.ConnectDb()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)
	log.Fatal(app.Listen(":3000"))
}
