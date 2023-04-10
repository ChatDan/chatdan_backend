package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	group := app.Group("/api")

	group.Post("/user/login", Login)
	group.Post("/user/register", Register)
}
