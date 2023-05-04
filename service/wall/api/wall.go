package api

import (
	"ChatDanBackend/common"
	"ChatDanBackend/common/gormx"
	"ChatDanBackend/common/schemax"
	"ChatDanBackend/service/wall/model"
	"github.com/gofiber/fiber/v2"
)

// ListWalls
// @Summary 获取今日表白墙
// @Tags Wall
// @Router /api/wall [get]
// @Produce json
// @Param json query WallRequest true "query"

func ListWalls(c *fiber.Ctx) error {
	query, err := common.ValidateQuery[WallRequest](c)
	if err != nil {
		return err
	}
	querySet := gormx.DB.Offset((query.PageNum - 1) * query.PageSize).Limit(query.PageSize)
	var walls []model.Wall

	err = querySet.Find(&walls).Error
	if err != nil {
		return err
	}
	return c.JSON(schemax.Response{
		Data: WallResponse{Posts: common.MustConvert[[]Post](walls)},
	})
}
