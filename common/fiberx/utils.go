package fiberx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func GetUserIDFromJWT(c *fiber.Ctx) (int, error) {
	tokenString := c.Cookies("access")
	if tokenString == "" {
		return 0, fiber.ErrUnauthorized
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return 0, fiber.ErrUnauthorized
	}
	if claims, ok := token.Claims.(UserClaim); ok {
		return claims.UserID, nil
	} else {
		return 0, fiber.ErrUnauthorized
	}
}
