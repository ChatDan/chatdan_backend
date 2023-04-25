package api

import "ChatDanBackend/common/schemax"

type Post struct {
	ID         string `json:"id"`
	PosterID   string `json:"poster_id"`
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

type WallRequest struct {
	schemax.PageRequest
}

type WallResponse struct {
	Posts []Post `json:"posts"`
}
