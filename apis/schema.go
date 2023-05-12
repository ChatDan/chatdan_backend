package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
)

// User

type LoginRequest struct {
	Username string `json:"username" validate:"min=2"`
	Password string `json:"password" validate:"min=8"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

type ResetRequest struct {
	OldPassword string `json:"old_password" validate:"min=8"`
	NewPassword string `json:"new_password" validate:"min=8"`
}

/* Box */

type BoxCommonResponse struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Title   string `json:"title"`
}

type BoxCreateRequest struct {
	Title string `json:"title" query:"title" validate:"required"`
}

type BoxListRequest struct {
	PageRequest
	Title string `json:"title" query:"title"`
	Owner int    `json:"owner" query:"owner" validate:"omitempty,min=0"`
}

type BoxListResponse struct {
	MessageBoxes []BoxCommonResponse `json:"messageBoxes"`
}

type BoxGetResponse struct {
	BoxCommonResponse
	PostsContent []string `json:"posts" copier:"-"`
}

type BoxModifyRequest struct {
	Title *string `json:"title" query:"title"`
}

func (b *BoxModifyRequest) IsEmpty() bool {
	return b.Title == nil
}

/* Post */

type PostCommonResponse struct {
	ID         int    `json:"id"`
	PosterID   int    `json:"poster_id"`
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
	IsOwner    bool   `json:"is_owner"`
}

func (p *PostCreateRequest) IsPublic() bool {
	return p.Visibility == Public
}

type PostListRequest struct {
	PageRequest
	BoxID int `json:"message_box_id" query:"message_box_id" validate:"required"`
}

type PostListResponse struct {
	Posts []PostCommonResponse
}

type PostGetResponse struct {
	PostCommonResponse
	Channels []string `json:"channels"`
}

type PostCreateRequest struct {
	BoxID      int    `json:"message_box_id" validate:"required,min=1"`
	Content    string `json:"content" validate:"required,min=1,max=2000"` // 限制长度
	Visibility string `json:"visibility" validate:"omitempty,oneof=public private" default:"public"`
}

type PostModifyRequest struct {
	Content    *string `json:"content" validate:"omitempty,min=1,max=2000"`
	Visibility *string `json:"visibility" validate:"omitempty,oneof=public private"`
}

func (p *PostModifyRequest) IsEmpty() bool {
	return p.Content == nil && p.Visibility == nil
}

func (p *PostModifyRequest) IsPublic() *bool {
	if p.Visibility == nil {
		return nil
	}
	isPublic := *p.Visibility == Public
	return &isPublic
}

type PostModifyResponse struct {
	PostCommonResponse
}

/* Channel */

type ChannelCommonResponse struct {
	ID      int    `json:"id"`
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	IsOwner bool   `json:"is_owner"`
}

type ChannelCreateRequest struct {
	PostID  int    `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type ChannelListRequest struct {
	PageRequest
	PostID int `json:"post_id" query:"post_id" validate:"required,min=1"`
}

type ChannelListResponse struct {
	Channels []ChannelCommonResponse `json:"channels"`
}

type ChannelModifyRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

// Wall

type WallListRequest struct {
	PageRequest
}

type WallListResponse struct {
	Posts []PostCommonResponse `json:"posts"`
}
