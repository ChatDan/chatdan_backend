package api

import (
	_ "ChatDanBackend/common"
	"github.com/gofiber/fiber/v2"
)

// Login godoc
// @Summary Login
// @Description Login
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/login [post]
// @Param body query LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} common.Response{data=UserResponse}
// @Failure 400 {object} common.Response "Bad Request"
// @Failure 404 {object} common.Response "User Not Found"
// @Failure 500 {object} common.Response "Internal Server Error"
func Login(c *fiber.Ctx) error {
	// todo: set cookie
	return c.JSON(nil)
}

// Register godoc
// @Summary Register
// @Description Register
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/register [post]
// @Param body query LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} common.Response{data=UserResponse}
// @Failure 400 {object} common.Response "Bad Request"
// @Failure 500 {object} common.Response "Internal Server Error"
func Register(c *fiber.Ctx) error {
	// todo: set cookie
	return c.JSON(nil)
}
