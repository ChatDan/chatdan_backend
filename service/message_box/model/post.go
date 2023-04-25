package model

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID         int            `json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`
	BoxID      int            `json:"box_id"`
	Content    string         `json:"content"`
	Visibility VisibilityType `json:"visibility"` // private or public
}

type VisibilityType = string

const (
	Private VisibilityType = "private"
	Public  VisibilityType = "public"
)
