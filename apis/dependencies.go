package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"encoding/base64"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strings"
)

func GetCurrentUser(c *fiber.Ctx, user *User) error {
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

	// parse jwt
	var claims UserClaims
	if err := parseJwt(accessToken, &claims); err != nil {
		Logger.Error("failed to parse jwt", zap.Error(err), zap.String("token", accessToken))
		return err
	}

	// convert to user
	user.ID = claims.UserID
	user.IsAdmin = claims.IsAdmin
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
