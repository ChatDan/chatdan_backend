package utilsx

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserID  int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func CreateJwtToken(claims UserClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
