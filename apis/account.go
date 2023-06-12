package apis

import (
	. "chatdan_backend/config"
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"time"
)

// Login godoc
// @Summary Login
// @Tags User Module
// @Accept json
// @Produce json
// @Router /user/login [post]
// @Param json body LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} RespForSwagger{data=LoginResponse}
// @Failure 401 {object} RespForSwagger "用户名或密码错误"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
func Login(c *fiber.Ctx) (err error) {
	// parse and validate body
	var body LoginRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	// get user from database
	var user User
	if err = DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Unauthorized("用户名或密码错误")
		}
		return err
	}

	// check password
	if !CheckPassword(body.Password, user.HashedPassword) {
		return Unauthorized("用户名或密码错误")
	}

	token, err := CreateJwtToken(&user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Path:    "/",
		Domain:  Config.Hostname,
		Expires: time.Now().Add(time.Hour),
	})

	// construct response
	var response LoginResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}
	response.AccessToken = token

	return Success(c, &response)
}

// Register godoc
// @Summary Register
// @Tags User Module
// @Accept json
// @Produce json
// @Router /user/register [post]
// @Param json body LoginRequest true "The two fields are required, you can also add other fields(e.g. email)."
// @Success 200 {object} RespForSwagger{data=LoginResponse}
// @Failure 400 {object} RespForSwagger "Bad Request"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
func Register(c *fiber.Ctx) (err error) {
	// parse and validate body
	var body LoginRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	// create user or fail if username exists
	var user User
	result := DB.Where(User{Username: body.Username}).
		Attrs(User{HashedPassword: MakePassword(body.Password)}).FirstOrCreate(&user)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return BadRequest("用户名已存在")
	}

	// insert into meilisearch
	if err = SearchAddOrReplace(user.ToSearchModel()); err != nil {
		return err
	}

	// create jwt token
	token, err := CreateJwtToken(&user)
	if err != nil {
		return
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   token,
		Path:    "/",
		Domain:  Config.Hostname,
		Expires: time.Now().Add(time.Hour),
	})

	// construct response
	var response LoginResponse
	if err = copier.CopyWithOption(&response, &user, CopyOption); err != nil {
		return
	}
	response.AccessToken = token

	return Success(c, &response)
}

// Reset godoc
// @Summary Reset Password
// @Tags User Module
// @Accept json
// @Produce json
// @Router /user/reset [post]
// @Param json body ResetRequest true "json"
// @Success 200 {object} RespForSwagger{data=UserResponse}
// @Failure 400 {object} RespForSwagger "Bad Request"
// @Failure 401 {object} RespForSwagger "Invalid JWT Token"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
func Reset(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}
	if err = DB.Take(&user, user.ID).Error; err != nil {
		return err
	}

	// parse and validate body
	var body ResetRequest
	if err = ValidateBody(c, &body); err != nil {
		return err
	}

	// check old password
	if !CheckPassword(body.OldPassword, user.HashedPassword) {
		return Unauthorized("原密码错误")
	}

	// update password
	if err = DB.Model(&user).Update("hashed_password", MakePassword(body.NewPassword)).Error; err != nil {
		return err
	}

	// construct response
	var response UserResponse
	if err = copier.Copy(&response, &user); err != nil {
		return
	}

	return Success(c, &response)
}

// Logout godoc
// @Summary Logout
// @Tags User Module
// @Produce json
// @Router /user/logout [post]
// @Success 200 {object} RespForSwagger{data=UserResponse}
// @Failure 401 {object} RespForSwagger "Invalid JWT Token"
// @Failure 500 {object} RespForSwagger "Internal Server Error"
func Logout(c *fiber.Ctx) (err error) {
	// get current user
	var user User
	if err = GetCurrentUser(c, &user); err != nil {
		return err
	}

	// delete jwt token
	if err = DeleteJwtToken(&user); err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "jwt",
		Value:   "",
		Path:    "/",
		Domain:  Config.Hostname,
		Expires: time.Now().Add(-time.Hour),
	})

	return Success(c, &EmptyStruct{})
}
