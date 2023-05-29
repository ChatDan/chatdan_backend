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
	group.Post("/user/logout", Logout)

	// User Info
	group.Get("/users", ListUsers) // admin only
	group.Get("/user/me", GetUserMe)
	group.Put("/user/me", ModifyUserMe)
	group.Delete("/user/me", DeleteUserMe)
	group.Get("/user/:id", GetAUser)
	group.Put("/user/:id", ModifyAUser)    // admin or owner only
	group.Delete("/user/:id", DeleteAUser) // admin only

	// Box
	group.Get("/messageBoxes", ListBoxes)
	group.Get("/messageBox/:id", GetABox)
	group.Post("/messageBox", CreateABox)
	group.Put("/messageBox/:id", ModifyABox)
	group.Delete("/messageBox/:id", DeleteABox)

	// Post
	group.Get("/posts", ListPosts)
	group.Get("/post/:id", GetAPost)
	group.Post("/post", CreateAPost)
	group.Put("/post/:id", ModifyAPost)
	group.Delete("/post/:id", DeleteAPost)

	// Channel
	group.Get("/channels", ListChannels)
	group.Get("/channel/:id", GetAChannel)
	group.Post("/channel", CreateAChannel)
	group.Put("/channel/:id", ModifyAChannel)
	group.Delete("/channel/:id", DeleteAChannel)

	// Wall
	group.Get("/wall", ListWalls)
	group.Get("/wall/:id", GetAWall)
	group.Post("/wall", CreateAWall)

	// Division
	group.Get("/divisions", ListDivisions)
	group.Get("/division/:id", GetADivision)
	group.Post("/division", CreateADivision)
	group.Put("/division/:id", ModifyADivision)
	group.Delete("/division/:id", DeleteADivision)

	// Topic
	group.Get("/topics", ListTopics)
	group.Get("/topic/:id", GetATopic)
	group.Post("/topic", CreateATopic)
	group.Put("/topic/:id", ModifyATopic)
	group.Delete("/topic/:id", DeleteATopic)
	group.Put("/topic/:id/_like/:data", LikeOrDislikeATopic)
	group.Put("/topic/:id/_view", ViewATopic)
	group.Put("/topic/:id/_favor", FavorATopic)
	group.Delete("/topic/:id/_unfavor", UnfavorATopic)
	group.Get("/topics/_favor", ListFavoriteTopics)
	group.Get("/topics/_user/:id", ListTopicsByUser)
	group.Get("/topics/_tag/:tag", ListTopicsByTag)
	group.Get("/topics/_search", SearchTopics)

	// Comment
	group.Get("/comments", ListComments)
	group.Get("/comment/:id", GetAComment)
	group.Post("/comment", CreateAComment)
	group.Put("/comment/:id", ModifyAComment)
	group.Delete("/comment/:id", DeleteAComment)
	group.Put("/comment/:id/_like/:data", LikeOrDislikeAComment)
	group.Get("/comments/_user/:id", ListCommentsByUser)
	group.Get("/comments/_search", SearchComments)

	// Tag
	group.Get("/tags", ListTags)
	group.Get("/tag/:id", GetATag)
	group.Post("/tag", CreateATag)
	group.Put("/tag/:id", ModifyATag) // admin only
	group.Delete("/tag/:id", DeleteATag)

	// Chat and Message
	group.Get("/chats", ListChats)
	group.Get("/messages", ListMessages)
	group.Post("/messages", CreateMessage)
}
