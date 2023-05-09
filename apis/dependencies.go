package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func GetCurrentUser(c *fiber.Ctx, user *User) error {
	accessToken := c.Cookies("jwt")
	if accessToken == "" {
		accessToken = c.Get("Authorization")
		if accessToken == "" {
			return Unauthorized()
		}
		if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
			accessToken = accessToken[7:]
		}
	}
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, nil, jwt.WithoutClaimsValidation())
	if err != nil {
		Logger.Error("invalid jwt token", zap.String("token", accessToken), zap.Error(err))
		return Unauthorized("invalid jwt token")
	}

	if userClaims, ok := token.Claims.(*UserClaims); ok {
		user.ID = userClaims.UserID
		user.IsAdmin = userClaims.IsAdmin
		c.Locals("user_id", user.ID)
		return nil
	} else {
		Logger.Error("invalid jwt token", zap.String("token", accessToken))
		return Unauthorized("invalid jwt token")
	}
}
