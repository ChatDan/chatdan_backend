package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
)

// ListDivisions godoc
// @Summary List all divisions
// @Tags Division Module
// @Produce json
// @Router /divisions [get]
// @Param body query DivisionListRequest true "page"
// @Success 200 {object} Response{data=DivisionListResponse}
// @Failure 400 {object} Response{data=ErrorDetail}
// @Failure 500 {object} Response
func ListDivisions(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// GetADivision godoc
// @Summary Get a division
// @Tags Division Module
// @Produce json
// @Router /division/{id} [get]
// @Param id path int true "division id"
// @Success 200 {object} Response{data=DivisionCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func GetADivision(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// CreateADivision godoc
// @Summary Create a division, admin only
// @Tags Division Module
// @Accept json
// @Produce json
// @Router /division [post]
// @Param json body DivisionCreateRequest true "division"
// @Success 201 {object} Response{data=DivisionCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func CreateADivision(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// ModifyADivision godoc
// @Summary Modify a division, admin only
// @Tags Division Module
// @Accept json
// @Produce json
// @Router /division/{id} [put]
// @Param id path int true "division id"
// @Param json body DivisionModifyRequest true "division"
// @Success 200 {object} Response{data=DivisionCommonResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func ModifyADivision(c *fiber.Ctx) (err error) {
	return Success(c, nil)
}

// DeleteADivision godoc
// @Summary Delete a division, admin only
// @Tags Division Module
// @Produce json
// @Router /division/{id} [delete]
// @Param id path int true "division id"
// @Success 200 {object} Response{data=EmptyStruct}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
func DeleteADivision(c *fiber.Ctx) (err error) {
	return Success(c, EmptyStruct{})
}
