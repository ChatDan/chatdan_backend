package api

import "github.com/gofiber/fiber/v2"

// Login godoc
// @Summary Login
// @Description Login
// @Tags Account
// @Accept  json
// @Produce  json
// @Param json body LoginRequest true "json"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} utils.MessageResponse
// @Failure 404 {object} utils.MessageResponse "User Not Found"
// @Failure 500 {object} utils.MessageResponse
func Login(c *fiber.Ctx) error {
	return c.JSON(nil)
}
