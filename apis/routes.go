package apis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", swagger.HandlerDefault)

	group := app.Group("/api")

	// User
	group.Post("/user/login", Login)
	group.Post("/user/register", Register)
	group.Post("/user/reset", Reset)

	// Box
	app.Get("/messageBoxes", ListBoxes)
	app.Get("/messageBox/:id", GetABox)
	app.Post("/messageBox", CreateABox)
	app.Put("/messageBox/:id", ModifyABox)
	app.Delete("/messageBox/:id", DeleteABox)

	// Post
	app.Get("/posts", ListPosts)
	app.Get("/post/:id", GetAPost)
	app.Post("/post", CreateAPost)
	app.Put("/post/:id", ModifyAPost)
	app.Delete("/post/:id", DeleteAPost)

	// Wall
	app.Get("/wall", ListWalls)
}
