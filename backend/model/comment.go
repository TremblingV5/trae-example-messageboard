package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	PostID    uint           `json:"post_id" gorm:"not null;index"`
	Post      *Post          `json:"post,omitempty" gorm:"foreignKey:PostID"`
	AuthorID  uint           `json:"author_id" gorm:"not null;index"`
	Author    *User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	ParentID  *uint          `json:"parent_id" gorm:"index"` // nil for root comments
	Parent    *Comment       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children  []Comment      `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Comment) TableName() string {
	return "comments"
}

// MaxCommentDepth is the maximum nesting depth for comments
const MaxCommentDepth = 5