package api

import (
	"ChatDanBackend/common"
	"ChatDanBackend/common/gormx"
	"ChatDanBackend/common/schemax"
	"ChatDanBackend/service/user/config"
	"ChatDanBackend/service/user/model"
	"ChatDanBackend/service/user/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

// Login godoc
// @Summary Login
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/login [post]
// @Param json body LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} common.Response{data=UserResponse}
// @Failure 401 {object} common.Response "用户名或密码错误"
// @Failure 500 {object} common.Response "Internal Server Error"
func Login(c *fiber.Ctx) error {
	body, err := common.ValidateBody[LoginRequest](c)
	if err != nil {
		return err
	}

	var user model.User
	if err := gormx.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return schemax.Unauthorized("用户名或密码错误")
		}
		return err
	}

	if !utils.CheckPassword(body.Password, user.HashedPassword) {
		return schemax.Unauthorized("用户名或密码错误")
	}

	token, err := utils.CreateJWT(&user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Path:    "/",
		Domain:  config.CustomConfig.Hostname,
		Expires: time.Now().Add(time.Hour),
	})

	return c.JSON(schemax.Success(common.MustConvert[UserResponse](user)))
}

// Register godoc
// @Summary Register
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/register [post]
// @Param json body LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} common.Response{data=UserResponse}
// @Failure 400 {object} common.Response "Bad Request"
// @Failure 500 {object} common.Response "Internal Server Error"
func Register(c *fiber.Ctx) error {
	body, err := common.ValidateBody[LoginRequest](c)
	if err != nil {
		return err
	}

	user := model.User{
		Username:       body.Username,
		HashedPassword: utils.MakePassword(body.Password),
	}

	err = gormx.DB.Transaction(func(tx *gorm.DB) error {
		var exists model.User
		if err := tx.Where("username = ?", user.Username).First(&exists).Error; err == nil {
			return schemax.BadRequest("用户已存在")
		}

		return tx.Create(&user).Error
	})
	if err != nil {
		return err
	}

	token, err := utils.CreateJWT(&user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Path:    "/",
		Domain:  config.CustomConfig.Hostname,
		Expires: time.Now().Add(time.Hour),
	})

	return c.JSON(schemax.Success(common.MustConvert[UserResponse](user)))
}

// Reset godoc
// @Summary Reset Password
// @Tags Account
// @Accept json
// @Produce json
// @Router /api/user/reset [post]
// @Param json body ResetRequest true
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "Bad Request"
// @Failure 401 {object} common.Response "Invalid JWT Token"
// @Failure 500 {object} common.Response "Internal Server Error"
func Reset(c *fiber.Ctx) error {
	body, err := common.ValidateBody[ResetRequest](c)
	if err != nil {
		return err
	}

	user := c.Locals("user").(*model.User)

	if !utils.CheckPassword(body.OldPassword, user.HashedPassword) {
		return schemax.Unauthorized("旧密码错误")
	}

	user.HashedPassword = utils.MakePassword(body.NewPassword)

	return gormx.DB.Save(&user).Error
}
