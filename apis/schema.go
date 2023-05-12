package apis

import (
	. "ChatDanBackend/models"
	. "ChatDanBackend/utils"
	"time"
)

// User

type LoginRequest struct {
	Username string `json:"username" validate:"min=2"`
	Password string `json:"password" validate:"min=8"`
}

type UserResponse struct {
	ID                  int     `json:"id"`
	Username            string  `json:"username"`
	IsAdmin             bool    `json:"is_admin"`
	Email               *string `json:"email" extensions:"x-nullable"`        // 邮箱，可选登录
	Avatar              *string `json:"avatar" extensions:"x-nullable"`       // 头像链接
	Introduction        *string `json:"introduction" extensions:"x-nullable"` // 个人简介/个性签名
	Banned              bool    `json:"banned"`                               // 是否被封禁
	TopicCount          int     `json:"topic_count"`                          // 发表的话题数
	CommentCount        int     `json:"comment_count"`                        // 发表的评论数
	FavoriteTopicsCount int     `json:"favorite_topics_count"`                // 收藏的话题数
	FollowedUsersCount  int     `json:"followed_users_count"`                 // 关注的用户数
	FollowingUsersCount int     `json:"following_users_count"`                // 粉丝数
}

type ResetRequest struct {
	OldPassword string `json:"old_password" validate:"min=8"`
	NewPassword string `json:"new_password" validate:"min=8"`
}

/* Box 提问箱 */

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

/* Post 帖子、提问 */

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

/* Channel 频道、回复 */

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

/* 表白墙 */

type WallListRequest struct {
	PageRequest
}

type WallListResponse struct {
	Posts []PostCommonResponse `json:"posts"`
}

/* Division */

type DivisionCommonResponse struct {
	ID           int                   `json:"id"`
	Name         string                `json:"name"`
	Description  *string               `json:"description" extensions:"x-nullable"`
	PinnedTopics []TopicCommonResponse `json:"pinned_topics" extensions:"x-nullable"`
}

type DivisionListRequest struct {
	PageRequest
}

type DivisionListResponse struct {
	Divisions []DivisionCommonResponse `json:"divisions"`
}

type DivisionCreateRequest struct {
	Name           string  `json:"name" validate:"required,min=1,max=20"`
	Description    *string `json:"description" validate:"omitempty,min=1,max=200"`
	PinnedTopicIDs []int   `json:"pinned_topic_ids" validate:"omitempty,dive,min=1"`
}

type DivisionModifyRequest struct {
	Name           *string `json:"name" validate:"omitempty,min=1,max=20"`
	Description    *string `json:"description" validate:"omitempty,min=1,max=200"`
	PinnedTopicIDs []int   `json:"pinned_topic_ids" validate:"omitempty,dive,min=1"`
}

func (d *DivisionModifyRequest) IsEmpty() bool {
	return d.Name == nil && d.Description == nil && len(d.PinnedTopicIDs) == 0
}

/* Topic */

type TopicCommonResponse struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"time_created"`
	UpdatedAt    time.Time `json:"time_updated"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	IsHidden     bool      `json:"is_hidden"`
	IsOwner      bool      `json:"is_owner"`
	IsAnonymous  bool      `json:"is_anonymous"`
	Anonyname    *string   `json:"anonyname,omitempty" extensions:"x-nullable"`
	PosterID     int       `json:"poster_id,omitempty"`
	Poster       *User     `json:"poster,omitempty"`
	DivisionID   int       `json:"division_id"`
	Tags         []string  `json:"tags"`
	LastComment  *Comment  `json:"last_comment,omitempty" extensions:"x-nullable"` // 按照时间排序最后一条评论或者按照点赞数排序最高赞的评论，创建之后为空
	ViewCount    int       `json:"view_count"`                                     // 浏览数
	LikeCount    int       `json:"like_count"`                                     // 点赞数
	DislikeCount int       `json:"dislike_count"`                                  // 点踩数
	CommentCount int       `json:"comment_count"`                                  // 评论数
	FavorCount   int       `json:"favorite_count"`                                 // 收藏数
}

type TopicListRequest struct {
	DivisionID     *int       `json:"division_id" query:"division_id" validate:"omitempty,min=1"`
	OrderBy        string     `json:"order_by" query:"order_by" validate:"omitempty,oneof=time_created time_updated" default:"time_updated"`
	CommentOrderBy string     `json:"comment_order_by" query:"comment_order_by" validate:"omitempty,oneof=id like" default:"id"`
	PageSize       int        `json:"page_size" query:"page_size" validate:"omitempty,min=1,max=100" default:"10"`
	StartTime      *time.Time `json:"start_time" query:"start_time" validate:"omitempty"`
}

type TopicListResponse struct {
	Topics []TopicCommonResponse `json:"topics"`
}

type TopicCreateRequest struct {
	Title       string   `json:"title" validate:"required,min=1,max=50"`
	Content     string   `json:"content" validate:"required,min=1,max=2000"`
	DivisionID  int      `json:"division_id" validate:"required,min=1"`
	IsAnonymous bool     `json:"is_anonymous"`
	Tags        []string `json:"tags"`
}

type TopicModifyRequest struct {
	Title      *string  `json:"title" validate:"omitempty,min=1,max=50"`
	Content    *string  `json:"content" validate:"omitempty,min=1,max=2000"`
	DivisionID *int     `json:"division_id" validate:"omitempty,min=1"` // admin only
	IsHidden   *bool    `json:"is_hidden"`                              // admin only
	Tags       []string `json:"tags"`                                   // owner or admin
}

func (t *TopicModifyRequest) IsEmpty() bool {
	return t.Title == nil && t.Content == nil && t.DivisionID == nil && t.IsHidden == nil && len(t.Tags) == 0
}

/* Comment */

type CommentCommonResponse struct {
	ID           int       `json:"id"`
	CreatedAt    time.Time `json:"time_created"`
	UpdatedAt    time.Time `json:"time_updated"`
	Content      string    `json:"content"`
	IsOwner      bool      `json:"is_owner"`
	IsHidden     bool      `json:"is_hidden"`
	TopicID      int       `json:"topic_id"`
	ReplyToID    *int      `json:"reply_to_id,omitempty" extensions:"x-nullable"`
	IsAnonymous  bool      `json:"is_anonymous"`
	Anonyname    *string   `json:"anonyname,omitempty" extensions:"x-nullable"`
	PosterID     int       `json:"poster_id,omitempty"`
	Poster       *User     `json:"poster,omitempty"`
	Ranking      int       `json:"ranking"`
	LikeCount    int       `json:"like_count"`    // 点赞数
	DislikeCount int       `json:"dislike_count"` // 点踩数
}

type CommentListRequest struct {
	PageRequest
	TopicID int    `json:"topic_id" query:"topic_id" validate:"required,min=1"`
	OrderBy string `json:"order_by" query:"order_by" validate:"omitempty,oneof=id like" default:"id"`
}

type CommentListResponse struct {
	Comments []CommentCommonResponse `json:"comments"`
}

type CommentCreateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type CommentModifyRequest struct {
	Content  *string `json:"content" validate:"omitempty,min=1,max=2000"`
	IsHidden *bool   `json:"is_hidden"` // admin only
}

func (c *CommentModifyRequest) IsEmpty() bool {
	return c.Content == nil && c.IsHidden == nil
}

/* Tag */

type TagCommonResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Temperature int    `json:"temperature"`
}

type TagListRequest struct {
	PageRequest
	OrderBy string `json:"order_by" query:"order_by" validate:"omitempty,oneof=id temperature" default:"id"`
	Search  string `json:"search" query:"search" validate:"omitempty,min=1,max=20"` // 搜索标签名
}

type TagListResponse struct {
	Tags []TagCommonResponse `json:"tags"`
}

type TagCreateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=20"`
}

type TagModifyRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=20"`
	Temperature *int    `json:"temperature" validate:"omitempty,min=1,max=100"`
}

func (t *TagModifyRequest) IsEmpty() bool {
	return t.Name == nil && t.Temperature == nil
}
