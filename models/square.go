package models

import (
	"chatdan_backend/utils"
	"github.com/juju/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

/* 广场 */

// Division 分区
type Division struct {
	ID             int            `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name           string         `json:"name" gorm:"not null;unique"`
	Description    *string        `json:"description"`
	PinnedTopicIDs []int          `json:"pinned_topic_ids" gorm:"serializer:json;not null;default:\"[]\""`
}

// Topic 帖子
type Topic struct {
	// 元数据
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Title       string         `json:"title" gorm:"not null"`                      // 帖子标题，必须有
	Content     string         `json:"content" gorm:"not null"`                    // 帖子内容，必须有
	IsAnonymous bool           `json:"is_anonymous" gorm:"not null;default:false"` // 是否匿名
	Anonyname   *string        `json:"anonyname"`                                  // 匿名时的昵称
	IsHidden    bool           `json:"is_hidden" gorm:"not null;default:false"`    // 是否隐藏，隐藏的帖子不会出现在列表中

	// 关联数据
	PosterID       int       `json:"poster_id" gorm:"not null"`                                  // 发帖人ID
	Poster         *User     `json:"poster" gorm:"foreignKey:PosterID"`                          // 发帖人
	DivisionID     int       `json:"division_id" gorm:"not null"`                                // 所属分区ID
	Division       *Division `json:"division" gorm:"foreignKey:DivisionID"`                      // 所属分区
	Tags           []*Tag    `json:"tags" gorm:"many2many:topic_tags"`                           // 帖子标签
	ViewedUsers    []*User   `json:"viewed_user" gorm:"many2many:topic_user_views"`              // 浏览过帖子的用户
	FavoredUsers   []*User   `json:"favored_user" gorm:"many2many:topic_user_favorites"`         // 收藏帖子的用户
	LikedUsers     []*User   `json:"liked_user" gorm:"many2many:topic_user_likes"`               // 点赞或点踩帖子的用户
	AnonymousUsers []*User   `json:"anonyname_mapping" gorm:"many2many:topic_anonyname_mapping"` // 匿名昵称映射

	// 统计数据
	ViewCount    int `json:"view_count" gorm:"not null;default:0"`     // 浏览数
	LikeCount    int `json:"like_count" gorm:"not null;default:0"`     // 点赞数
	DislikeCount int `json:"dislike_count" gorm:"not null;default:0"`  // 点踩数
	CommentCount int `json:"comment_count" gorm:"not null;default:0"`  // 评论数
	FavorCount   int `json:"favorite_count" gorm:"not null;default:0"` // 收藏数
}

func (t Topic) GetID() int {
	return t.ID
}

func (t Topic) TableName() string {
	return "topic"
}

func (t *Topic) FindOrCreateTags(tx *gorm.DB, tagNames []string) (err error) {
	// batch find and create if not exists
	var tags []*Tag
	err = tx.Where("name IN ?", tagNames).Find(&tags).Error
	if err != nil {
		return
	}

	var newTagNames []string
	for _, tagName := range tagNames {
		var found bool
		for _, tag := range tags {
			if tag.Name == tagName {
				found = true
				break
			}
		}
		if !found {
			newTagNames = append(newTagNames, tagName)
		}
	}

	var newTags []*Tag
	for _, tagName := range newTagNames {
		newTags = append(newTags, &Tag{
			Name: tagName,
		})
	}
	err = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&newTags).Error
	if err != nil {
		return
	}
	return nil
}

func (t *Topic) TagContents() []string {
	var contents []string
	for _, tag := range t.Tags {
		contents = append(contents, tag.Name)
	}
	return contents
}

func (t *Topic) AfterCreate(tx *gorm.DB) (err error) {
	if !t.IsAnonymous {
		t.Poster = new(User)
		err = LoadModel(tx.Where("id = ?", t.PosterID), &t.Poster)
	}
	return
}

func (t *Topic) AfterFind(tx *gorm.DB) (err error) {
	if !t.IsAnonymous {
		t.Poster = new(User)
		err = LoadModel(tx.Where("id = ?", t.PosterID), &t.Poster)
	}
	return
}

func (t Topic) ToSearchModel() TopicSearchModel {
	return TopicSearchModel{
		ID:         t.ID,
		Title:      t.Title,
		Content:    t.Content,
		CreatedAt:  int(t.CreatedAt.UnixMicro()),
		UpdatedAt:  int(t.UpdatedAt.UnixMicro()),
		PosterID:   t.PosterID,
		DivisionID: t.DivisionID,
	}
}

type TopicSearchModel struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  int    `json:"created_at"`
	UpdatedAt  int    `json:"updated_at"`
	PosterID   int    `json:"poster_id"`
	DivisionID int    `json:"division_id"`
}

func (t TopicSearchModel) GetID() int {
	return t.ID
}

func (TopicSearchModel) IndexName() string {
	return "topic"
}

func (TopicSearchModel) PrimaryKey() string {
	return "id"
}

func (TopicSearchModel) FilterableAttributes() []string {
	return []string{"poster_id", "division_id"}
}

func (t TopicSearchModel) SearchableAttributes() []string {
	return []string{"title", "content"}
}

func (t TopicSearchModel) SortableAttributes() []string {
	return []string{"id", "created_at", "updated_at"}
}

func (t TopicSearchModel) RankingRules() []string {
	return []string{"words", "attribute", "sort", "exactness"}
}

var _ SearchModel = TopicSearchModel{}

// Comment 评论
type Comment struct {
	// 元数据
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Content     string         `json:"content" gorm:"not null"`                    // 评论内容
	IsAnonymous bool           `json:"is_anonymous" gorm:"not null;default:false"` // 是否匿名
	Anonyname   *string        `json:"anonyname"`                                  // 匿名时的昵称，可选
	Ranking     int            `json:"ranking" gorm:"not null;default:0"`          // 评论的楼层
	ReplyToID   *int           `json:"reply_to_id"`                                // 回复的评论ID，必须是本帖子的评论
	IsHidden    bool           `json:"is_hidden" gorm:"not null;default:false"`    // 是否被隐藏，被隐藏的评论不会显示在帖子中

	// 关联数据
	PosterID   int     `json:"poster_id" gorm:"not null"`
	Poster     *User   `json:"poster" gorm:"foreignKey:PosterID"`
	TopicID    int     `json:"topic_id" gorm:"not null"`
	Topic      *Topic  `json:"topic" gorm:"foreignKey:TopicID"`
	LikedUsers []*User `json:"liked_user" gorm:"many2many:comment_user_likes"` // 点赞或点踩评论的用户

	// 统计数据
	LikeCount    int `json:"like_count" gorm:"not null;default:0"`    // 点赞数
	DislikeCount int `json:"dislike_count" gorm:"not null;default:0"` // 点踩数
}

func (c Comment) GetID() int {
	return c.ID
}

func (Comment) TableName() string {
	return "comment"
}

func (c *Comment) AfterFind(tx *gorm.DB) (err error) {
	if !c.IsAnonymous {
		c.Poster = new(User)
		err = LoadModel(tx.Where("id = ?", c.PosterID), &c.Poster)
	}
	return nil
}

func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	if !c.IsAnonymous {
		c.Poster = new(User)
		err = LoadModel(tx.Where("id = ?", c.PosterID), &c.Poster)
	}
	return nil
}

// CommentSearchModel 评论搜索模型
type CommentSearchModel struct {
	ID        int    `json:"id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	TopicID   int    `json:"topic_id"`
	Content   string `json:"content"`
	PosterID  int    `json:"poster_id"`
}

func (c CommentSearchModel) GetID() int {
	return c.ID
}

func (c CommentSearchModel) IndexName() string {
	return "comment"
}

func (c CommentSearchModel) PrimaryKey() string {
	return "id"
}

func (c CommentSearchModel) FilterableAttributes() []string {
	return []string{"topic_id", "poster_id"}
}

func (c CommentSearchModel) SearchableAttributes() []string {
	return []string{"content"}
}

func (c CommentSearchModel) SortableAttributes() []string {
	return []string{"id", "created_at", "updated_at"}
}

func (c CommentSearchModel) RankingRules() []string {
	return []string{"words", "attribute", "sort", "exactness"}
}

var _ SearchModel = CommentSearchModel{}

// Tag 标签
type Tag struct {
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name        string         `json:"name" gorm:"not null;unique"`                 // 标签名，大小写不敏感
	Temperature int            `json:"temperature" gorm:"not null;default:0;index"` // 热度，表示有多少帖子使用了这个标签，越高表示越热门，用于排序

	// 关联数据
	Topics []*Topic `json:"topics" gorm:"many2many:topic_tags"`
}

func (t Tag) GetID() int {
	return t.ID
}

func (t Tag) TableName() string {
	return "tag"
}

type TagSearchModel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Temperature int    `json:"temperature"`
}

func (t TagSearchModel) GetID() int {
	return t.ID
}

func (TagSearchModel) IndexName() string {
	return "tag"
}

func (TagSearchModel) PrimaryKey() string {
	return "id"
}

func (TagSearchModel) FilterableAttributes() []string {
	return []string{}
}

func (TagSearchModel) SearchableAttributes() []string {
	return []string{"name"}
}

func (TagSearchModel) SortableAttributes() []string {
	return []string{"temperature"}
}

func (TagSearchModel) RankingRules() []string {
	return []string{"words", "attribute", "sort", "exactness"}
}

var _ SearchModel = TagSearchModel{}

// TopicUserLikes 用户点赞或点踩帖子
type TopicUserLikes struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	TopicID   int       `json:"topic_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	LikeData  int       `json:"like_data" gorm:"not null;default:0"` // 1表示点赞，-1表示点踩
}

// TopicUserFavorites 用户收藏帖子
// 默认按照创建时间倒序返回
type TopicUserFavorites struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	TopicID   int       `json:"topic_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

// TopicUserViews 用户浏览过的帖子
// 默认按照更新时间倒序返回
type TopicUserViews struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	TopicID   int       `json:"topic_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Count     int       `json:"count" gorm:"not null;default:0"` // 浏览次数
}

type TopicAnonynameMapping struct {
	TopicID   int    `json:"topic_id" gorm:"primaryKey"`
	UserID    int    `json:"user_id" gorm:"primaryKey"`
	Anonyname string `json:"anonyname" gorm:"not null"`
}

func NewAnonyname(tx *gorm.DB, topicID, userID int) (string, error) {
	name := utils.NewRandName()
	return name, tx.Create(&TopicAnonynameMapping{
		TopicID:   topicID,
		UserID:    userID,
		Anonyname: name,
	}).Error
}

func FindOrGenerateAnonyname(tx *gorm.DB, topicID, userID int) (string, error) {
	var anonyname string
	err := tx.
		Model(&TopicAnonynameMapping{}).
		Select("anonyname").
		Where("hole_id = ?", topicID).
		Where("user_id = ?", userID).
		Take(&anonyname).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var names []string
			err = tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Model(&TopicAnonynameMapping{}).
				Where("hole_id = ?", topicID).
				Order("anonyname asc").
				Pluck("anonyname", &names).Error
			if err != nil {
				return "", err
			}

			anonyname = utils.GenerateName(names)
			err = tx.Create(&TopicAnonynameMapping{
				TopicID:   topicID,
				UserID:    userID,
				Anonyname: anonyname,
			}).Error
			if err != nil {
				return anonyname, err
			}
		} else {
			return "", err
		}
	}
	return anonyname, nil
}

type CommentUserLikes struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	CommentID int       `json:"comment_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	LikeData  int       `json:"like_data" gorm:"not null;default:0"` // 1表示点赞，-1表示点踩
}
