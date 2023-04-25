package api

import "ChatDanBackend/common/schemax"

// Box

type BoxCommonResponse struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	Title   string `json:"title"`
}

type BoxCreateRequest struct {
	Title string `json:"title" query:"title" validate:"required"`
}

type BoxCreateResponse struct {
	BoxCommonResponse
}

type BoxListRequest struct {
	schemax.PageRequest
	Title string `json:"title" query:"title"`
	Owner int    `json:"owner" query:"owner" validate:"min=0"`
}

type BoxListResponse struct {
	MessageBoxes []BoxCommonResponse `json:"messageBoxes"`
}

type BoxGetResponse struct {
	BoxCommonResponse
	Posts []string `json:"posts" copier:"-"`
}

type BoxModifyRequest struct {
	Title *string `json:"title" query:"title"`
}

type BoxModifyResponse struct {
	BoxCommonResponse
}

type BoxDeleteResponse struct {
	Message string `json:"message"`
}

// Post

type PostCommonResponse struct {
	ID         string `json:"id"`
	PosterID   string `json:"poster_id"`
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
}

type PostCreateRequest struct {
	MessageBoxID string `json:"message_box_id" query:"message_box_id" validate:"required"`
	Content      string `json:"content" query:"content" validate:"required"`
	Visibility   string `json:"visibility" query:"visibility" validate:"required,oneof=public private"`
}

type PostCreateResponse struct {
	PostCommonResponse
}

type PostListRequest struct {
	schemax.PageRequest
	MessageBoxID string `json:"message_box_id" query:"message_box_id" validate:"required"`
}

type PostListResponse struct {
	PostCommonResponse
}

type PostGetResponse struct {
	PostCommonResponse
	Channels []string `json:"channels"`
}

type PostModifyRequest struct {
	Content    *string `json:"content"`
	Visibility *string `json:"visibility"`
}

type PostModifyResponse struct {
	PostCommonResponse
}

type PostDeleteResponse struct {
	Message string `json:"message"`
}
