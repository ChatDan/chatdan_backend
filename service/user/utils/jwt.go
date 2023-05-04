package utils

import (
	"ChatDanBackend/service/user/model"
)

func CreateJWT(user *model.User) {
	claim := jwt.MapClaim{
		"uid": user.ID,
	}

}
