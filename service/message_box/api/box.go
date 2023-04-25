package api

import (
	"ChatDanBackend/common"
	"ChatDanBackend/common/gormx"
	"ChatDanBackend/common/schemax"
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
	query, err := common.ValidateQuery[BoxListRequest](c)
	if err != nil {
		return err
	}

	querySet := gormx.DB.Offset((query.PageNum - 1) * query.PageSize).Limit(query.PageSize)
	if query.Title != "" {
		querySet = querySet.Where("title = ?", query.Title)
	}
	if query.Owner != 0 {
		querySet = querySet.Where("owner_id = ?", query.Owner)
	}

	var boxes []model.Box
	err = querySet.Find(&boxes).Error
	if err != nil {
		return err
	}

	return c.JSON(schemax.Response{Data: BoxListResponse{MessageBoxes: common.MustConvert[[]BoxCommonResponse](boxes)}})
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
	// get box id
	boxID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	// load box from database
	var box model.Box
	err = gormx.DB.Model(model.Box{}).Preload("Posts").First(&box, boxID).Error
	if err != nil {
		return err
	}

	// convert to response
	response := common.MustConvert[BoxGetResponse](box)
	for _, post := range box.Posts {
		response.Posts = append(response.Posts, post.Content)
	}

	return c.JSON(schemax.Response{Data: response})
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
