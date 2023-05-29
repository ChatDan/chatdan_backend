package apis

import (
	. "chatdan_backend/models"
	. "chatdan_backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/juju/errors"
	"github.com/oleiade/reflections"
	"gorm.io/gorm"
	"time"
)

// common

type ModifyRequest interface {
	IsEmpty() bool
}

func ModifyRequestLevelValidation(sl validator.StructLevel) {
	if sl.Current().Interface().(ModifyRequest).IsEmpty() {
		sl.ReportError(sl.Current().Interface(), "ModifyRequest", "ModifyRequest", "modify", "should not empty")
	}
}

func init() {
	Validate.RegisterStructValidation(ModifyRequestLevelValidation, UserModifyRequest{})
	Validate.RegisterStructValidation(ModifyRequestLevelValidation, BoxModifyRequest{})
	Validate.RegisterStructValidation(ModifyRequestLevelValidation, PostModifyRequest{})
	Validate.RegisterStructValidation(ModifyRequestLevelValidation, DivisionModifyRequest{})
	Validate.RegisterStructValidation(ModifyRequestLevelValidation, TopicModifyRequest{})
	Validate.RegisterStructValidation(ModifyRequestLevelValidation, CommentModifyRequest{})
	Validate.RegisterStructValidation(ModifyRequestLevelValidation, TagModifyRequest{})
}

// User

type RegisterRequest struct {
	LoginRequest
	Email        *string `json:"email" validate:"omitempty,email"`
	Avatar       *string `json:"avatar" validate:"omitempty"`
	Introduction *string `json:"introduction" validate:"omitempty,min=2"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"min=2"`
	Password string `json:"password" validate:"min=8"`
}

type UserResponse struct {
	ID                  int     `json:"id"`
	Username            string  `json:"username"`
	IsAdmin             bool    `json:"is_admin"`
	Email               *string `json:"email,omitempty" extensions:"x-nullable"`        // 邮箱，可选登录
	Avatar              *string `json:"avatar,omitempty" extensions:"x-nullable"`       // 头像链接
	Introduction        *string `json:"introduction,omitempty" extensions:"x-nullable"` // 个人简介/个性签名
	Banned              bool    `json:"banned"`                                         // 是否被封禁
	TopicCount          int     `json:"topic_count"`                                    // 发表的话题数
	CommentCount        int     `json:"comment_count"`                                  // 发表的评论数
	FavoriteTopicsCount int     `json:"favorite_topics_count"`                          // 收藏的话题数
	FollowersCount      int     `json:"followers_count"`                                // 被关注数
	FollowingUsersCount int     `json:"following_users_count"`                          // 关注数
}

type LoginResponse struct {
	UserResponse
	AccessToken string `json:"access_token"`
}

type ResetRequest struct {
	OldPassword string `json:"old_password" validate:"min=8"`
	NewPassword string `json:"new_password" validate:"min=8"`
}

type UserListRequest struct {
	PageRequest
}

type UserListResponse struct {
	Users   []UserResponse `json:"users"`
	Version int            `json:"version"`
	Total   int            `json:"total"` // User 总数，便于前端分页
}

type UserModifyRequest struct {
	Username     *string `json:"username" validate:"omitempty,min=2"`
	Email        *string `json:"email" validate:"omitempty,email"`
	Avatar       *string `json:"avatar" validate:"omitempty,url"`
	Introduction *string `json:"introduction" validate:"omitempty,min=2"`
}

func (u UserModifyRequest) IsEmpty() bool {
	return u.Username == nil && u.Email == nil && u.Avatar == nil && u.Introduction == nil
}

func (u UserModifyRequest) Fields() []string {
	fields, _ := reflections.Fields(u)
	return fields
}

/* Box 提问箱 */

type BoxCommonResponse struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	OwnerID   int       `json:"owner_id"`
	Title     string    `json:"title"`
	PostCount int       `json:"post_count"`
	ViewCount int       `json:"view_count"`
}

type BoxCreateRequest struct {
	Title string `json:"title" query:"title" validate:"required"`
}

type BoxListRequest struct {
	PageRequest
	Title   string `json:"title" query:"title"`
	Owner   int    `json:"owner" query:"owner" validate:"omitempty,min=0"`
	OrderBy string `json:"order_by" query:"order_by" validate:"oneof='id asc' 'updated_at desc'" default:"id asc"`
}

type BoxListResponse struct {
	MessageBoxes []BoxCommonResponse `json:"messageBoxes"`
	Version      int                 `json:"version"`
	Total        int                 `json:"total"` // Box 总数，便于前端分页
}

type BoxGetResponse struct {
	BoxCommonResponse
	PostsContent []string `json:"posts" copier:"-"`
}

type BoxModifyRequest struct {
	Title *string `json:"title" query:"title"`
}

func (b BoxModifyRequest) IsEmpty() bool {
	return b.Title == nil
}

/* Post 帖子、提问 */

type PostCommonResponse struct {
	ID           int           `json:"id"`
	PosterID     int           `json:"poster_id"`
	Poster       *UserResponse `json:"poster,omitempty"`
	Content      string        `json:"content"`
	Visibility   string        `json:"visibility"` // public private
	IsOwner      bool          `json:"is_owner"`
	IsAnonymous  bool          `json:"is_anonymous"`
	Anonyname    string        `json:"anonyname"`
	ChannelCount int           `json:"channel_count"`
	ViewCount    int           `json:"view_count"`
}

type PostListRequest struct {
	PageRequest
	BoxID int `json:"message_box_id" query:"message_box_id" validate:"required"`
}

type PostListResponse struct {
	Posts   []PostCommonResponse `json:"posts"`
	Version int                  `json:"version"`
	Total   int                  `json:"total"` // Post 总数，便于前端分页
}

type PostGetResponse struct {
	PostCommonResponse
	Channels []string `json:"channels"`
}

type PostCreateRequest struct {
	BoxID       int    `json:"message_box_id" validate:"required,min=1"`
	Content     string `json:"content" validate:"required,min=1,max=2000"` // 限制长度
	Visibility  string `json:"visibility" validate:"omitempty,oneof=public private" default:"public"`
	IsAnonymous *bool  `json:"is_anonymous" validate:"omitempty"`
}

func (p *PostCreateRequest) IsPublic() bool {
	return p.Visibility == Public
}

func (p *PostCreateRequest) SetDefaults() {
	if p.IsAnonymous == nil {
		p.IsAnonymous = new(bool)
		*p.IsAnonymous = true
	}
}

type PostModifyRequest struct {
	Content    *string `json:"content" validate:"omitempty,min=1,max=2000"`
	Visibility *string `json:"visibility" validate:"omitempty,oneof=public private"`
}

func (p PostModifyRequest) IsEmpty() bool {
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
	Version  int                     `json:"version"`
	Total    int                     `json:"total"` // Channel 总数，便于前端分页
}

type ChannelModifyRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

/* 表白墙 */

type WallCommonResponse struct {
	ID          int           `json:"id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	IsAnonymous bool          `json:"is_anonymous"`
	PosterID    int           `json:"poster_id"`        // 匿名时为 0
	Poster      *UserResponse `json:"poster,omitempty"` // 匿名时为 null
	Content     string        `json:"content"`
	Visibility  string        `json:"visibility"`
	IsShown     bool          `json:"is_shown"` // 是否显示在表白墙页面
}

type WallListRequest struct {
	PageRequest
	Date *time.Time `json:"date" query:"date" validate:"omitempty"` // 日期（所对应的时间，只解析时间），不填默认当天（即昨天发送的表白墙）
}

type WallListResponse struct {
	Posts []WallCommonResponse `json:"posts"`
	Total int                  `json:"total"`                     // Post 总数，便于前端分页
	Date  time.Time            `json:"date" swaggertype:"string"` // 日期
}

type WallCreateRequest struct {
	Content     string `json:"content" validate:"required,min=1,max=2000"`
	IsAnonymous *bool  `json:"is_anonymous" validate:"omitempty"` // 是否匿名，不填默认匿名
}

func (w *WallCreateRequest) SetDefaults() {
	if w.IsAnonymous == nil {
		w.IsAnonymous = new(bool)
		*w.IsAnonymous = true
	}
}

type WallModifyRequest struct {
	Visibility *string `json:"visibility" validate:"omitempty,oneof=public private"` // 管理员修改可见性
}

func (w WallModifyRequest) IsEmpty() bool {
	return w.Visibility == nil
}

func (w *WallModifyRequest) IsPublic() *bool {
	if w.Visibility == nil {
		return nil
	}
	isPublic := *w.Visibility == Public
	return &isPublic
}

/* Division */

type DivisionCommonResponse struct {
	ID           int                   `json:"id"`
	Name         string                `json:"name"`
	Description  *string               `json:"description" extensions:"x-nullable"`
	PinnedTopics []TopicCommonResponse `json:"pinned_topics" extensions:"x-nullable"`
}

type DivisionListResponse struct {
	Divisions []DivisionCommonResponse `json:"divisions"`
}

type DivisionCreateRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=20"`
	Description *string `json:"description" validate:"omitempty,min=1,max=200"`
}

type DivisionModifyRequest struct {
	Name           *string `json:"name" validate:"omitempty,min=1,max=20"`
	Description    *string `json:"description" validate:"omitempty,min=1,max=200"`
	PinnedTopicIDs []int   `json:"pinned_topic_ids" validate:"omitempty,dive,min=1"`
}

type DivisionDeleteRequest struct {
	To int `json:"to" validate:"min=0" default:"1"`
}

func (d DivisionModifyRequest) IsEmpty() bool {
	return d.Name == nil && d.Description == nil && len(d.PinnedTopicIDs) == 0
}

/* Topic */

type TopicCommonResponse struct {
	ID          int                    `json:"id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	IsHidden    bool                   `json:"is_hidden"`
	IsAnonymous bool                   `json:"is_anonymous"`
	Anonyname   *string                `json:"anonyname,omitempty" extensions:"x-nullable"`
	PosterID    int                    `json:"poster_id,omitempty"`
	Poster      *UserResponse          `json:"poster,omitempty"`
	DivisionID  int                    `json:"division_id"`
	Tags        []TagCommonResponse    `json:"tags"`
	LastComment *CommentCommonResponse `json:"last_comment,omitempty" extensions:"x-nullable"` // 按照时间排序最后一条评论或者按照点赞数排序最高赞的评论，创建之后为空

	// 统计数据
	ViewCount    int `json:"view_count"`     // 浏览数
	LikeCount    int `json:"like_count"`     // 点赞数
	DislikeCount int `json:"dislike_count"`  // 点踩数
	CommentCount int `json:"comment_count"`  // 评论数
	FavorCount   int `json:"favorite_count"` // 收藏数

	// 动态生成的字段
	IsOwner  bool `json:"is_owner"`
	Liked    bool `json:"liked"`
	Disliked bool `json:"disliked"`
	Favored  bool `json:"favored"`
}

func (t *TopicCommonResponse) Postprocess(c *fiber.Ctx) (err error) {
	userID := c.Locals("user_id").(int)
	if userID == t.PosterID {
		t.IsOwner = true
	}

	// load last comment
	var comment Comment
	err = DB.Last(&comment, "topic_id = ?", t.ID).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	} else {
		var commentResponse CommentCommonResponse
		err = copier.Copy(&commentResponse, &comment)
		if err != nil {
			return
		}

		t.LastComment = &commentResponse
	}

	// load like
	var like TopicUserLikes
	err = DB.Where("topic_id = ? AND user_id = ?", t.ID, userID).First(&like).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	} else {
		switch like.LikeData {
		case 1:
			t.Liked = true
		case -1:
			t.Disliked = true
		}
	}

	// load favorite
	var favorite TopicUserFavorites
	err = DB.Where("topic_id = ? AND user_id = ?", t.ID, userID).First(&favorite).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	} else {
		t.Favored = true
	}

	// clear user info
	if t.IsAnonymous {
		t.Poster = nil
		t.PosterID = 0
	}

	return nil
}

type TopicListRequest struct {
	DivisionID     *int       `json:"division_id" query:"division_id" validate:"omitempty,min=1"`
	OrderBy        string     `json:"order_by" query:"order_by" validate:"omitempty,oneof=created_at updated_at" default:"created_at"` // 排序方式，created_at 按照收藏的时间排序，updated_at 按照主题帖更新的时间排序
	CommentOrderBy string     `json:"comment_order_by" query:"comment_order_by" validate:"omitempty,oneof=id like" default:"id"`
	PageSize       int        `json:"page_size" query:"page_size" validate:"omitempty,min=1,max=100" default:"10"`
	StartTime      *time.Time `json:"start_time" query:"start_time" validate:"omitempty"`
}

type TopicSearchRequest struct {
	PageRequest
	Search string `json:"search" query:"search" validate:"omitempty,min=1,max=100"`
}

type TopicListResponse struct {
	Topics []TopicCommonResponse `json:"topics"`
}

func (t *TopicListResponse) Postprocess(c *fiber.Ctx) (err error) {
	userID := c.Locals("user_id").(int)
	for i := range t.Topics {
		if userID == t.Topics[i].PosterID {
			t.Topics[i].IsOwner = true
		}
	}

	// batch load last comment
	var comments []Comment
	var topicIDs []int
	for _, topic := range t.Topics {
		topicIDs = append(topicIDs, topic.ID)
	}
	err = DB.Raw(
		`select * from comments 
			where topic_id in (?) and id in (
				select max(id) from comments group by topic_id
			)`, topicIDs).Scan(&comments).Error
	if err != nil {
		return err
	}
	var commentResponses []CommentCommonResponse
	err = copier.Copy(&commentResponses, &comments)
	if err != nil {
		return
	}

	for i := range t.Topics {
		for j := range commentResponses {
			if t.Topics[i].ID == commentResponses[j].TopicID {
				t.Topics[i].LastComment = &commentResponses[j]
				err = commentResponses[j].Postprocess(c)
				if err != nil {
					return
				}
				break
			}
		}
	}

	// batch load like
	var likes []TopicUserLikes
	err = DB.Where("topic_id in (?) AND user_id = ?", topicIDs, userID).Find(&likes).Error
	if err != nil {
		return
	}

	for i := range t.Topics {
		for j := range likes {
			if t.Topics[i].ID == likes[j].TopicID {
				switch likes[j].LikeData {
				case 1:
					t.Topics[i].Liked = true
				case -1:
					t.Topics[i].Disliked = true
				}
				break
			}
		}
	}

	// batch load favorite
	var favorites []TopicUserFavorites
	err = DB.Where("topic_id in (?) AND user_id = ?", topicIDs, userID).Find(&favorites).Error
	if err != nil {
		return
	}

	for i := range t.Topics {
		for j := range favorites {
			if t.Topics[i].ID == favorites[j].TopicID {
				t.Topics[i].Favored = true
				break
			}
		}
	}

	// clear user info
	for i := range t.Topics {
		if t.Topics[i].IsAnonymous {
			t.Topics[i].Poster = nil
			t.Topics[i].PosterID = 0
		}
	}

	return
}

type TopicCreateRequest struct {
	Title       string             `json:"title" validate:"required,min=1,max=50"`
	Content     string             `json:"content" validate:"required,min=1,max=2000"`
	DivisionID  int                `json:"division_id" validate:"required,min=1"`
	IsAnonymous bool               `json:"is_anonymous"` // 默认不传为 false
	Tags        []TagCreateRequest `json:"tags" validate:"omitempty,dive,min=1,max=10,dive"`
}

type TopicModifyRequest struct {
	Title      *string            `json:"title" validate:"omitempty,min=1,max=50"`
	Content    *string            `json:"content" validate:"omitempty,min=1,max=2000"`
	DivisionID *int               `json:"division_id" validate:"omitempty,min=1"`           // admin only
	IsHidden   *bool              `json:"is_hidden"`                                        // admin only
	Tags       []TagCreateRequest `json:"tags" validate:"omitempty,dive,min=1,max=10,dive"` // owner or admin
}

func (t TopicModifyRequest) IsEmpty() bool {
	return t.Title == nil && t.Content == nil && t.DivisionID == nil && t.IsHidden == nil && len(t.Tags) == 0
}

func (t *TopicModifyRequest) Fields() []string {
	var fields []string
	if t.Title != nil {
		fields = append(fields, "Title")
	}
	if t.Content != nil {
		fields = append(fields, "Content")
	}
	if t.DivisionID != nil {
		fields = append(fields, "DivisionID")
	}
	if t.IsHidden != nil {
		fields = append(fields, "IsHidden")
	}
	return fields
}

/* Comment */

type CommentCommonResponse struct {
	ID          int           `json:"id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Content     string        `json:"content"`
	IsHidden    bool          `json:"is_hidden"`
	TopicID     int           `json:"topic_id"`
	ReplyToID   *int          `json:"reply_to_id,omitempty" extensions:"x-nullable"`
	IsAnonymous bool          `json:"is_anonymous"`
	Anonyname   *string       `json:"anonyname,omitempty" extensions:"x-nullable"`
	PosterID    int           `json:"poster_id,omitempty"`
	Poster      *UserResponse `json:"poster,omitempty"`
	Ranking     int           `json:"ranking"`

	// 统计数据
	LikeCount    int `json:"like_count"`    // 点赞数
	DislikeCount int `json:"dislike_count"` // 点踩数

	// 动态生成的字段
	IsOwner  bool `json:"is_owner"`
	Liked    bool `json:"liked"`
	Disliked bool `json:"disliked"`
}

func (comment *CommentCommonResponse) Postprocess(c *fiber.Ctx) (err error) {
	userID := c.Locals("user_id").(int)
	// set owner
	if comment.PosterID == userID {
		comment.IsOwner = true
	}

	// load like
	var like CommentUserLikes
	err = DB.Where("comment_id = ? AND user_id = ?", comment.ID, userID).First(&like).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	} else {
		switch like.LikeData {
		case 1:
			comment.Liked = true
		case -1:
			comment.Disliked = true
		}
	}

	// clear user info
	if comment.IsAnonymous {
		comment.Poster = nil
		comment.PosterID = 0
	}
	return
}

type CommentListRequest struct {
	PageRequest
	TopicID int    `json:"topic_id" query:"topic_id" validate:"required,min=1"`
	OrderBy string `json:"order_by" query:"order_by" validate:"omitempty,oneof=id like" default:"id"`
}

type CommentListByUserRequest struct {
	PageRequest
	OrderBy string `json:"order_by" query:"order_by" validate:"omitempty,oneof=id like" default:"id"`
}

type CommentSearchRequest struct {
	PageRequest
	Search string `json:"search" query:"search" validate:"required,min=1"`
}

type CommentListResponse struct {
	Comments []CommentCommonResponse `json:"comments"`
}

func (comments *CommentListResponse) Postprocess(c *fiber.Ctx) (err error) {
	if len(comments.Comments) == 0 {
		return
	}
	userID := c.Locals("user_id").(int)

	// set owner
	for i := range comments.Comments {
		if comments.Comments[i].PosterID == userID {
			comments.Comments[i].IsOwner = true
		}
	}

	// batch load like
	var likes []CommentUserLikes
	commentIDs := make([]int, len(comments.Comments))
	for i := range comments.Comments {
		commentIDs[i] = comments.Comments[i].ID
	}
	err = DB.Where("comment_id in (?) AND user_id = ?", commentIDs, userID).Find(&likes).Error
	if err != nil {
		return
	}
	for i := range comments.Comments {
		for j := range likes {
			if comments.Comments[i].ID == likes[j].CommentID {
				switch likes[j].LikeData {
				case 1:
					comments.Comments[i].Liked = true
				case -1:
					comments.Comments[i].Disliked = true
				}
				break
			}
		}
	}

	// clear user info
	for i := range comments.Comments {
		if comments.Comments[i].IsAnonymous {
			comments.Comments[i].Poster = nil
			comments.Comments[i].PosterID = 0
		}
	}

	return
}

type CommentCreateRequest struct {
	TopicID     int    `json:"topic_id" validate:"required,min=1"`
	ReplyToID   *int   `json:"reply_to_id"`
	Content     string `json:"content" validate:"required,min=1,max=2000"`
	IsAnonymous bool   `json:"is_anonymous"` // 默认实名
}

type CommentModifyRequest struct {
	Content  *string `json:"content" validate:"omitempty,min=1,max=2000"`
	IsHidden *bool   `json:"is_hidden"` // admin only
}

func (c CommentModifyRequest) IsEmpty() bool {
	return c.Content == nil && c.IsHidden == nil
}

/* Tag */

type TagCommonResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Temperature int    `json:"temperature"`
}

func (t TagCommonResponse) GetName() string {
	return t.Name
}

type TagListRequest struct {
	PageRequest
	OrderBy string `json:"order_by" query:"order_by" validate:"omitempty,oneof='id asc' 'temperature desc'" default:"id asc"`
	Search  string `json:"search" query:"search" validate:"omitempty,min=1,max=20"` // 搜索标签名
}

type TagListResponse struct {
	Tags    []TagCommonResponse `json:"tags"`
	Version int                 `json:"version"`
	Total   int                 `json:"total"`
}

type TagCreateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=20"`
}

func (t TagCreateRequest) GetName() string {
	return t.Name
}

type Tagger interface {
	GetName() string
}

func ToTagNames[T Tagger](tags []T) (names []string) {
	names = make([]string, len(tags))
	for i := range tags {
		names[i] = tags[i].GetName()
	}
	return
}

type TagModifyRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=20"`
	Temperature *int    `json:"temperature" validate:"omitempty,min=1,max=100"`
}

func (t TagModifyRequest) IsEmpty() bool {
	return t.Name == nil && t.Temperature == nil
}

/* Chat */

type ChatCommonResponse struct {
	ID                 int           `json:"id"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	OneUserID          int           `json:"one_user_id"`
	OneUser            *UserResponse `json:"one_user"`
	AnotherUserID      int           `json:"another_user_id"`
	AnotherUser        *UserResponse `json:"another_user"`
	LastMessageContent string        `json:"last_message_content"`
	MessageCount       int           `json:"message_count"`
}

type ChatListResponse struct {
	Chats []ChatCommonResponse `json:"chats"` // 返回时按照 UpdatedAt 降序排列
}

/* Message */

type MessageCommonResponse struct {
	ID         int       `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	Content    string    `json:"content"`
	FromUserID int       `json:"from_user_id"`
	ToUserID   int       `json:"to_user_id"`
	IsOwner    bool      `json:"is_me"`
}

type MessageCreateRequest struct {
	Content  string `json:"content" validate:"required,min=1,max=2000"`
	ToUserID int    `json:"to_user_id" validate:"required,min=1"`
}

type MessageListRequest struct {
	PageSize  int        `json:"page_size" query:"page_size" validate:"omitempty,min=1,max=100" default:"10"`
	ToUserID  int        `json:"to_user_id" query:"to_user_id" validate:"required,min=1"`
	StartTime *time.Time `json:"start_time" query:"start_time" validate:"omitempty"` // 不填默认为当前时间
}

type MessageListResponse struct {
	Messages []MessageCommonResponse `json:"messages"` // 按照 CreatedAt 倒序排列
}
