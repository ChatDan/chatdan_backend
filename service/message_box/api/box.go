package api

import "github.com/gofiber/fiber/v2"

// CreateABox godoc
// @Summary Create a message box
// @Description Create a message box
// @Tags Box
// @Accept json
// @Produce json
// @Param box query BoxCreateRequest true "box"
// @Success 200 {object} common.Response{data=BoxCreateResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /messageBox [post]
func CreateABox(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// ListBoxes godoc
// @Summary List message boxes
// @Description List message boxes
// @Tags Box
// @Accept json
// @Produce json
// @Param body query BoxListRequest true "page"
// @Success 200 {object} common.Response{data=BoxListResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /messageBoxes [get]
func ListBoxes(c *fiber.Ctx) error {
	return c.JSON(nil)
}

// GetABox godoc
// @Summary Get a message box
// @Description Get a message box
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
// @Description Modify a message box, owner only
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
// @Description Delete a message box, owner only
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
