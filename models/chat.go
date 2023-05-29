package models

import (
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"index"`

	// 关联数据
	OneUserID          int           `json:"one_user_id" gorm:"index;index:idx_chat_one_another,priority:1"` // one_user_id < another_user_id
	OneUser            *User         `json:"-" gorm:"foreignKey:OneUserID"`
	AnotherUserID      int           `json:"another_user_id" gorm:"index;index:idx_chat_one_another,priority:2"`
	AnotherUser        *User         `json:"-" gorm:"foreignKey:AnotherUserID"`
	LastMessageID      int           `json:"last_message_id"`
	LastMessageContent string        `json:"last_message_content"`
	Messages           []ChatMessage `json:"messages"`

	// 统计数据
	MessageCount int `json:"message_count"`
}

type ChatMessage struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	Content   string    `json:"content"`

	// 关联数据\
	ChatID     int   `json:"chat_id" gorm:"index"`
	Chat       *Chat `json:"-" gorm:"foreignKey:ChatID"`
	FromUserID int   `json:"from_user_id"`
	FromUser   *User `json:"-" gorm:"foreignKey:FromUserID"`
	ToUserID   int   `json:"to_user_id"`
	ToUser     *User `json:"-" gorm:"foreignKey:ToUserID"`
}
