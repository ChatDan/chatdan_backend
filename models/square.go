package models

import (
	"gorm.io/gorm"
	"time"
)

/* 广场 */

// Division 分区
type Division struct {
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"time_created"`
	UpdatedAt   time.Time      `json:"time_updated"`
	DeletedAt   gorm.DeletedAt `json:"time_deleted" gorm:"index"`
	Name        string         `json:"name" gorm:"not null;unique"`
	Description string         `json:"description"`
}

// Topic 帖子
type Topic struct {
	// 元数据
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"time_created"`
	UpdatedAt   time.Time      `json:"time_updated"`
	DeletedAt   gorm.DeletedAt `json:"time_deleted" gorm:"index"`
	Title       string         `json:"title" gorm:"not null"`                      // 帖子标题，必须有
	Content     string         `json:"content" gorm:"not null"`                    // 帖子内容，必须有
	IsAnonymous bool           `json:"is_anonymous" gorm:"not null;default:false"` // 是否匿名
	Anonyname   *string        `json:"anonyname"`                                  // 匿名时的昵称
	IsHidden    bool           `json:"is_hidden" gorm:"not null;default:false"`    // 是否隐藏，隐藏的帖子不会出现在列表中

	// 关联数据
	PosterID     int       `json:"poster_id" gorm:"not null"`                    // 发帖人ID
	Poster       *User     `json:"poster" gorm:"foreignKey:PosterID"`            // 发帖人
	DivisionID   int       `json:"division_id" gorm:"not null"`                  // 所属分区ID
	Division     *Division `json:"division" gorm:"foreignKey:DivisionID"`        // 所属分区
	Tags         []*Tag    `json:"tags" gorm:"many2many:topic_tags"`             // 帖子标签
	FavoredUsers []*User   `json:"favored_user" gorm:"many2many:user_favorites"` // 收藏帖子的用户
	LikedUsers   []*User   `json:"liked_user" gorm:"many2many:topic_user_likes"` // 点赞或点踩帖子的用户

	// 统计数据
	LikeCount    int `json:"like_count" gorm:"not null;default:0"`
	DislikeCount int `json:"dislike_count" gorm:"not null;default:0"`
	CommentCount int `json:"comment_count" gorm:"not null;default:0"`
	FavorCount   int `json:"favorite_count" gorm:"not null;default:0"`
}

// Comment 评论
type Comment struct {
	// 元数据
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"time_created"`
	UpdatedAt   time.Time      `json:"time_updated"`
	DeletedAt   gorm.DeletedAt `json:"time_deleted" gorm:"index"`
	Content     string         `json:"content" gorm:"not null"`                    // 评论内容
	IsAnonymous bool           `json:"is_anonymous" gorm:"not null;default:false"` // 是否匿名
	Anonyname   *string        `json:"anonyname"`                                  // 匿名时的昵称，可选
	Ranking     int            `json:"ranking" gorm:"not null;default:0"`          // 评论的楼层
	ReplyTo     *int           `json:"reply_to"`                                   // 回复的评论ID，必须是本帖子的评论
	IsHidden    bool           `json:"is_hidden" gorm:"not null;default:false"`    // 是否被隐藏，被隐藏的评论不会显示在帖子中

	// 关联数据
	PosterID int    `json:"poster_id" gorm:"not null"`
	Poster   *User  `json:"poster" gorm:"foreignKey:PosterID"`
	TopicID  int    `json:"topic_id" gorm:"not null"`
	Topic    *Topic `json:"topic" gorm:"foreignKey:TopicID"`

	// 统计数据
	LikeCount    int `json:"like_count" gorm:"not null;default:0"`    // 点赞数
	DislikeCount int `json:"dislike_count" gorm:"not null;default:0"` // 点踩数
}

// Tag 标签
type Tag struct {
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"time_created"`
	UpdatedAt   time.Time      `json:"time_updated"`
	DeletedAt   gorm.DeletedAt `json:"time_deleted" gorm:"index"`
	Name        string         `json:"name" gorm:"not null;unique"`                 // 标签名，大小写不敏感
	Temperature int            `json:"temperature" gorm:"not null;default:0;index"` // 热度，表示有多少帖子使用了这个标签，越高表示越热门，用于排序

	// 关联数据
	Topics []*Topic `json:"topics" gorm:"many2many:topic_tags"`
}

// TopicUserLikes 用户点赞或点踩帖子
type TopicUserLikes struct {
	UserID   int `json:"user_id" gorm:"primaryKey"`
	TopicID  int `json:"topic_id" gorm:"primaryKey"`
	LikeData int `json:"like_data" gorm:"not null;default:0"` // 1表示点赞，-1表示点踩
}
