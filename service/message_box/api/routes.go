package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	// Box
	app.Post("/messageBox", CreateABox)
	app.Get("/messageBoxes", ListBoxes)
	app.Get("/messageBox/:id", GetABox)
	app.Put("/messageBox/:id", ModifyABox)
	app.Delete("/messageBox/:id", DeleteABox)

	// Post
	app.Post("/post", CreateAPost)
	app.Get("/posts", ListPosts)
	app.Get("/post/:id", GetAPost)
	app.Put("/post/:id", ModifyAPost)
	app.Delete("/post/:id", DeleteAPost)
}
