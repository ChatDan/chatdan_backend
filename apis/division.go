package apis

import (
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ListDivisions godoc
// @Summary List all divisions
// @Tags Division Module
// @Produce json
// @Router /divisions [get]
// @Success 200 {object} RespForSwagger{data=DivisionListResponse}
// @Failure 400 {object} RespForSwagger{data=ErrorDetail}
// @Failure 500 {object} RespForSwagger
func ListDivisions(c *fiber.Ctx) (err error) {
	var divisions []Division
	result := DB.Find(&divisions)
	if result.Error != nil {
		return result.Error
	}

	var response DivisionListResponse
	if err = copier.Copy(&response.Divisions, &divisions); err != nil {
		return err
	}

	return Success(c, &response)
}

// GetADivision godoc
// @Summary Get a division
// @Tags Division Module
// @Produce json
// @Router /division/{id} [get]
// @Param id path int true "division id"
// @Success 200 {object} RespForSwagger{data=DivisionCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func GetADivision(c *fiber.Ctx) (err error) {
	divisionID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var division Division
	if err = DB.First(&division, divisionID).Error; err != nil {
		return err
	}

	var response DivisionCommonResponse
	if err = copier.Copy(&response, &division); err != nil {
		return err
	}

	return Success(c, &response)
}

// CreateADivision godoc
// @Summary Create a division, admin only
// @Tags Division Module
// @Accept json
// @Produce json
// @Router /division [post]
// @Param json body DivisionCreateRequest true "division"
// @Success 201 {object} RespForSwagger{data=DivisionCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func CreateADivision(c *fiber.Ctx) (err error) {
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}
	if !user.IsAdmin {
		return Forbidden()
	}

	var body DivisionCreateRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	var division Division
	// see https://gorm.io/zh_CN/docs/advanced_query.html#FirstOrCreate
	result := DB.Where(Division{Name: body.Name}).Attrs(Division{Description: body.Description}).FirstOrCreate(&division)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return BadRequest("division already exists")
	}

	var response DivisionCommonResponse
	if err = copier.Copy(&response, &division); err != nil {
		return err
	}
	return Created(c, &response)
}

// ModifyADivision godoc
// @Summary Modify a division, admin only
// @Tags Division Module
// @Accept json
// @Produce json
// @Router /division/{id} [put]
// @Param id path int true "division id"
// @Param json body DivisionModifyRequest true "division"
// @Success 200 {object} RespForSwagger{data=DivisionCommonResponse}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func ModifyADivision(c *fiber.Ctx) (err error) {
	var user User
	err = GetCurrentUser(c, &user)
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return Forbidden()
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var body DivisionModifyRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	var division Division
	if err = DB.Transaction(func(tx *gorm.DB) error {
		// load division with lock
		if err = tx.Clauses(LockClause).First(&division, id).Error; err != nil {
			return err
		}

		// copy body to division
		if err = copier.CopyWithOption(&division, &body, copier.Option{IgnoreEmpty: true}); err != nil {
			return err
		}

		// update division
		return tx.Model(&division).Updates(&division).Error
	}); err != nil {
		return err
	}

	var response DivisionCommonResponse
	if err = copier.Copy(&response, &division); err != nil {
		return err
	}
	return Success(c, &response)
}

// DeleteADivision godoc
// @Summary Delete a division, admin only
// @Tags Division Module
// @Produce json
// @Router /division/{id} [delete]
// @Param id path int true "division id"
// @Success 200 {object} RespForSwagger{data=EmptyStruct}
// @Failure 400 {object} RespForSwagger
// @Failure 500 {object} RespForSwagger
func DeleteADivision(c *fiber.Ctx) (err error) {
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return Forbidden()
	}

	var body DivisionDeleteRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	var division Division
	if err = DB.Transaction(func(tx *gorm.DB) error {
		// load division with lock
		if err = tx.Clauses(LockClause).First(&division, id).Error; err != nil {
			return err
		}

		err = tx.Exec("UPDATE Topic SET division_id = ? WHERE division_id = ?", body.To, id).Error
		if err != nil {
			return err
		}

		return DB.Delete(&division).Error
	}); err != nil {
		return err
	}
	return Success(c, &EmptyStruct{})
}
