package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID             int            `json:"id"`
	Username       string         `json:"username" gorm:"index,size:30"`
	Email          string         `json:"email" gorm:"index"`
	HashedPassword string         `json:"-" gorm:"size:256"`
	LoginTime      time.Time      `json:"-" gorm:"autoUpdateTime"`
	RegisterTime   time.Time      `json:"-" gorm:"autoCreateTime"`
	Banned         bool           `json:"banned"`
	IsAdmin        bool           `json:"is_admin"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}
