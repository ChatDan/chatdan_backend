package model

import "ChatDanBackend/service/user/model"

type Wall struct {
	ID       int `json:"id"`
	PosterID int `json:"poster_id"`
	Poster   *model.User
}
