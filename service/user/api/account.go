package api

import "github.com/gofiber/fiber/v2"

// Login godoc
// @Summary Login
// @Description Login
// @Tags Account
// @Accept  json
// @Produce  json
// @Param json body LoginRequest true "json"
// @Success 200 {object} common.Response{data=UserResponse}
// @Failure 400 {object} common.Response "Bad Request"
// @Failure 404 {object} common.Response "User Not Found"
// @Failure 500 {object} common.Response "Internal Server Error"
func Login(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func Register(c *fiber.Ctx) error {
	return c.JSON(nil)
}
