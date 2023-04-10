package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username" gorm:"index"`
	Email        string    `json:"email" gorm:"index"` // 邮箱验证码
	Password     string    `json:"password"`
	RegisterIP   string    `json:"-"`
	LoginIP      []string  `json:"-" gorm:"serializer:json"`
	LoginTime    time.Time `json:"-" gorm:"autoUpdateTime"`
	RegisterTime time.Time `json:"-" gorm:"autoCreateTime"`
	Banned       bool      `json:"banned"`
	IsAdmin      bool      `json:"is_admin"`
	DeletedAt    time.Time `json:"-"`
}
