package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-backend/pkg/controllers"
)

func Setup(app *fiber.App) {
	// User apis
	app.Post("/api/register", controllers.CreateUser)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
}
