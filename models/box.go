package models

import (
	"gorm.io/gorm"
	"time"
)

type Box struct {
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"time_created"`
	UpdatedAt time.Time      `json:"time_updated"`
	DeletedAt gorm.DeletedAt `json:"time_deleted"`
	OwnerID   int            `json:"owner_id"`
	Owner     *User          `json:"owner" gorm:"foreignKey:OwnerID"`
	Title     string         `json:"title"`
	Posts     []Post         `json:"posts"`
}

type Post struct {
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"time_created"`
	UpdatedAt time.Time      `json:"time_updated"`
	DeletedAt gorm.DeletedAt `json:"time_deleted"`
	PosterID  int            `json:"poster_id"`
	Poster    *User          `json:"poster" gorm:"foreignKey:PosterID"`
	BoxID     int            `json:"message_box_id"`
	Box       *Box           `json:"box" gorm:"foreignKey:BoxID"`
	Content   string         `json:"content"`
	IsPublic  bool           `json:"is_public"` // true if the post is public
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

type Channel struct {
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"time_created"`
	UpdatedAt time.Time      `json:"time_updated"`
	DeletedAt gorm.DeletedAt `json:"time_deleted"`
	OwnerID   int            `json:"poster_id"`
	Owner     *User          `json:"poster" gorm:"foreignKey:OwnerID"`
	PostID    int            `json:"post_id"`
	Post      *Post          `json:"post" gorm:"foreignKey:PostID"`
	Content   string         `json:"content"`
}
