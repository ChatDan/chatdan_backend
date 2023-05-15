package models

import (
	"gorm.io/gorm"
	"time"
)

type Wall struct {
	ID          int            `json:"id"`
	CreatedAt   time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"index"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Content     string         `json:"content"`
	Visibility  string         `json:"visibility"`
	IsAnonymous bool           `json:"is_anonymous"`

	// 关联数据
	PosterID int   `json:"poster_id"`
	Poster   *User `json:"-" gorm:"foreignKey:PosterID"`
}

func (w *Wall) IsPublic() bool {
	return w.Visibility == Public
}
