package models

import (
	"gorm.io/gorm"
	"time"
)

// Box 提问箱，MessageBox的简称
// 一个用户可以有多个提问箱，一个提问箱可以有多个提问
type Box struct {
	// 元数据
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Title     string         `json:"title"`

	// 关联数据
	OwnerID int    `json:"owner_id"`
	Owner   *User  `json:"owner" gorm:"foreignKey:OwnerID"`
	Posts   []Post `json:"posts"` // 一个提问箱可以有多个提问

	// 统计数据
	PostCount int `json:"post_count"` // 回复数
	ViewCount int `json:"view_count"` // 浏览数
}

func (Box) TableName() string {
	return "box"
}

func (b Box) GetID() int {
	return b.ID
}

type BoxSearchModel struct {
	ID        int    `json:"id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	Title     string `json:"title"`
	OwnerID   int    `json:"owner_id"`
}

func (BoxSearchModel) TableName() string {
	return "box"
}

func (BoxSearchModel) IndexName() string {
	return "box"
}

func (b BoxSearchModel) GetID() int {
	return b.ID
}

func (BoxSearchModel) PrimaryKey() string {
	return "id"
}

func (BoxSearchModel) FilterableAttributes() []string {
	return []string{"owner_id"}
}

func (BoxSearchModel) SearchableAttributes() []string {
	return []string{"title"}
}

func (BoxSearchModel) SortableAttributes() []string {
	return []string{"id", "updated_at", "created_at"}
}

func (BoxSearchModel) RankingRules() []string {
	return []string{"words", "attribute", "sort", "exactness"}
}

func (BoxSearchModel) ReloadModel() error {
	var boxes []Box
	return DB.FindInBatches(&boxes, 100, func(tx *gorm.DB, batch int) error {
		var boxSearchModels []BoxSearchModel
		for _, box := range boxes {
			boxSearchModels = append(boxSearchModels, box.ToBoxSearchModel())
		}
		return SearchAddOrReplaceInBatch(boxSearchModels)
	}).Error
}

func (b *Box) ToBoxSearchModel() BoxSearchModel {
	return BoxSearchModel{
		ID:        b.ID,
		CreatedAt: int(b.CreatedAt.UnixMicro()),
		UpdatedAt: int(b.UpdatedAt.UnixMicro()),
		Title:     b.Title,
		OwnerID:   b.OwnerID,
	}
}

// Post 提问、帖子
// 一个提问箱可以有多个提问，一个提问包含一个回复 Thread，Thread 里的元素是追问追答的 Channel
// 被提问者的回答和提问者的追问都是 Channel
type Post struct {
	// 元数据
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	Content     string         `json:"content"`
	IsPublic    bool           `json:"is_public"`    // true if the post is public
	IsAnonymous bool           `json:"is_anonymous"` // true if the post is anonymous

	// 关联数据
	PosterID int       `json:"poster_id"`
	Poster   *User     `json:"poster" gorm:"foreignKey:PosterID"`
	BoxID    int       `json:"message_box_id"`
	Box      *Box      `json:"box" gorm:"foreignKey:BoxID"`
	Channel  []Channel `json:"channel"` // 一个提问包含一个回复 Thread，Thread 里的元素是追问追答的 Channel

	// 统计数据
	ChannelCount int `json:"channel_count"` // 回复数
	ViewCount    int `json:"view_count"`    // 浏览数
}

func (p *Post) Visibility() string {
	if p.IsPublic {
		return Public
	}
	return Private
}

const (
	Private string = "private"
	Public         = "public"
)

// Channel 回复、追问、追答
// 一个提问包含一个回复 Thread，Thread 里的元素是追问追答的 Channel，即 Thread = []Channel
type Channel struct {
	// 元数据
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Content   string         `json:"content"`

	// 关联数据
	OwnerID int   `json:"owner_id"`
	Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`
	PostID  int   `json:"post_id"`
	Post    *Post `json:"post" gorm:"foreignKey:PostID"`
}
