package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

// ListWalls
// @Summary 获取今日表白墙
// @Tags Wall Module
// @Router /wall [get]
// @Produce json
// @Param json query WallListRequest true "query"
// @Success 200 {object} Response{data=WallListResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response "服务器错误"
func ListWalls(c *fiber.Ctx) (err error) {
	// get and validate query
	var query WallListRequest
	if err = ValidateQuery(c, &query); err != nil {
		return err
	}

	// construct querySet and load walls from database
	var walls []Wall
	if err = query.QuerySet(DB).Find(&walls).Error; err != nil {
		return err
	}

	// construct response
	var response WallListResponse
	if err = copier.CopyWithOption(&response.Posts, &walls, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	return Success(c, response)
}
