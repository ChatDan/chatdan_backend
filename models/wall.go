package models

import (
	"time"
)

type Wall struct {
	ID         int        `json:"id"`
	CreatedAt  time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"index"`
	DeletedAt  *time.Time `json:"-" gorm:"index"`
	PosterID   int        `json:"poster_id"`
	Poster     *User      `json:"-" gorm:"foreignKey:PosterID"`
	Content    string     `json:"content"`
	Visibility string     `json:"visibility"`
}
