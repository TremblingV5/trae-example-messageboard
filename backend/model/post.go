package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	AuthorID  uint           `json:"author_id" gorm:"not null;index"`
	Author    *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	ImageURL  string         `json:"image_url" gorm:"size:255"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Post) TableName() string {
	return "posts"
}