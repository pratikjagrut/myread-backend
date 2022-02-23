package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-backend/pkg/controllers"
	"github.com/pratikjagrut/myreads-backend/pkg/models"
)

func Setup(app *fiber.App) {
	// User apis
	app.Post("/api/register", controllers.CreateUser)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.GetUser)
	app.Post("/api/logout", controllers.Logout)

	// Book apis
	app.Post("/api/books/add", controllers.AddBook)
	app.Get("/api/books/all", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, "")
	})
	app.Get("/api/books/reading", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, models.Reading)
	})
	app.Get("/api/books/finished", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, models.Finished)
	})
	app.Get("/api/books/wishlist", func(c *fiber.Ctx) error {
		return controllers.GetBoooks(c, models.Wishlist)
	})
	app.Post("/api/books/updatestatus", controllers.UpdateStatus)

	app.Post("/api/books/deletebook", controllers.DeleteBook)
}
