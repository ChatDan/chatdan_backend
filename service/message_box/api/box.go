package api

import (
	"ChatDanBackend/common"
	"ChatDanBackend/common/gormx"
	"ChatDanBackend/service/message_box/model"
	"github.com/gofiber/fiber/v2"
)

// CreateABox godoc
// @Summary Create a message box
// @Tags Box
// @Accept json
// @Produce json
// @Param box body BoxCreateRequest true "box"
// @Success 200 {object} common.Response{data=BoxCreateResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /messageBox [post]
func CreateABox(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// ListBoxes godoc
// @Summary List message boxes
// @Tags Box
// @Accept json
// @Produce json
// @Param body query BoxListRequest true "page"
// @Success 200 {object} common.Response{data=BoxListResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /messageBoxes [get]
func ListBoxes(c *fiber.Ctx) error {
	// validate body
	query, err := common.ValidateBody[BoxListRequest](c)
	if err != nil {
		return err
	}

	querySet := gormx.DB.Offset((query.PageNum - 1) * query.PageSize).Limit(query.PageSize)

	var boxes []model.Box
	err = querySet.Find(&boxes).Error
	if err != nil {
		return err
	}

	return c.JSON(common.MustConvert[BoxListResponse](boxes))
}

// GetABox godoc
// @Summary Get a message box
// @Tags Box
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} common.Response{data=BoxGetResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /messageBox/{id} [get]
func GetABox(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// ModifyABox godoc
// @Summary Modify a message box
// @Tags Box
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param body query BoxModifyRequest true "box"
// @Success 200 {object} common.Response{data=BoxModifyResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /messageBox/{id} [put]
func ModifyABox(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// DeleteABox godoc
// @Summary Delete a message box
// @Tags Box
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} common.Response{data=BoxDeleteResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /messageBox/{id} [delete]
func DeleteABox(c *fiber.Ctx) error {
	return c.JSON(nil)
}
