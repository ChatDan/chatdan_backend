package apis

import (
	"ChatDanBackend/config"
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"encoding/base64"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/juju/errors"
	"go.uber.org/zap"
	"strings"
)

func GetCurrentUser(c *fiber.Ctx, user *User) (err error) {
	if config.Config.Mode == "dev" {
		user.ID = 1
		user.IsAdmin = true
		return nil
	}

	// get access token from cookie "jwt"
	accessToken := c.Cookies("jwt")
	if accessToken == "" {
		// get access token from header "Authorization"
		accessToken = c.Get("Authorization")
		if accessToken == "" {
			return Unauthorized()
		}
		if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
			accessToken = accessToken[7:]
		}
	}

	if config.Config.Standalone {
		token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if userClaims, ok := token.Claims.(*UserClaims); !ok {
				return nil, errors.New("invalid jwt token")
			} else {
				var userJwtSecret UserJwtSecret
				if err = DB.Take(&userJwtSecret, userClaims.ID).Error; err != nil {
					return nil, err
				}
				return []byte(userJwtSecret.Secret), nil
			}
		})
		if err != nil {
			Logger.Error("failed to parse jwt", zap.Error(err), zap.String("token", accessToken))
			return Unauthorized()
		}

		if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			user.ID = userClaims.UserID
			user.IsAdmin = userClaims.IsAdmin
			c.Locals("user_id", user.ID)
			return nil
		} else {
			Logger.Error("invalid jwt token", zap.String("token", accessToken))
			return Unauthorized()
		}
	} else {
		// parse jwt
		var claims UserClaims
		if err = parseJwt(accessToken, &claims); err != nil {
			Logger.Error("failed to parse jwt", zap.Error(err), zap.String("token", accessToken))
			return err
		}

		// convert to user
		user.ID = claims.UserID
		user.IsAdmin = claims.IsAdmin
	}
	return nil
}

func parseJwt(token string, claims *UserClaims) (err error) {
	// split token into 3 parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Unauthorized()
	}

	// jwt encoding ignores padding, so RawStdEncoding should be used instead of StdEncoding
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		Logger.Error("failed to decode jwt", zap.Error(err), zap.String("token", token))
		return Unauthorized()
	}

	// decode payload
	return json.Unmarshal(data, claims)
}
