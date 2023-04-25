package model

import (
	"gorm.io/gorm"
	"time"
)

type Box struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	OwnerID   string         `json:"owner_id"`
	Title     string         `json:"title"`
}
