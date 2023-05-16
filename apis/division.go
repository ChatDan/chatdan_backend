package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
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
	var divisions []*Division
	result := DB.Find(&divisions)
	if result.Error != nil {
		return result.Error
	}

	var response DivisionListResponse
	if err = copier.CopyWithOption(&response, &divisions, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	return Success(c, response)
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
	divisionID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var division Division
	result := DB.First(&division, divisionID)
	if result.Error != nil {
		return result.Error
	}
	var response DivisionCommonResponse
	if err = copier.CopyWithOption(&response, &division, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	return Success(c, response)
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
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return Forbidden()
	}
	var body DivisionCreateRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}
	division := Division{
		Name:        body.Name,
		Description: body.Description,
	}
	result := DB.FirstOrCreate(&division, Division{Name: body.Name})
	if result.Error != nil {
		return result.Error
	}
	var response DivisionCommonResponse
	if err = copier.CopyWithOption(&response, &division, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	return Created(c, response)
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
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return Forbidden()
	}
	var body DivisionModifyRequest
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	id, err = c.ParamsInt("id")
	if err != nil {
		return err
	}
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}
	if body.IsEmpty() {
		return BadRequest()
	}
	var division Division
	result := DB.First(&division, id)
	if result.Error != nil {
		return result.Error
	}

	if err = copier.CopyWithOption(&division, &body, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}

	result = DB.Model(&division).Updates(division)
	if result.Error != nil {
		return result.Error
	}
	var response DivisionCommonResponse
	if err = copier.CopyWithOption(&response, &division, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	return Success(c, response)
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
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return Forbidden()
	}

	var body DivisionDeleteRequest
	err = ValidateBody(c, &body)
	if err != nil {
		return err
	}

	var division Division
	result := DB.First(&division, id)

	if result.RowsAffected == 0 {
		return NotFound()
	}

	err = DB.Exec("UPDATE Topic SET division_id = ? WHERE division_id = ?", body.To, id).Error
	if err != nil {
		return err
	}
	err = DB.Delete(&Division{ID: id}).Error
	if err != nil {
		return err
	}
	return Success(c, EmptyStruct{})
}
