package api

import "github.com/gofiber/fiber/v2"

// ListWalls
// @Summary 获取今日表白墙
// @Tags Wall
// @Router /api/wall [get]
// @Produce json
// @Param json query WallRequest true "query"
func ListWalls(c *fiber.Ctx) error {
	return c.JSON(nil)
}
