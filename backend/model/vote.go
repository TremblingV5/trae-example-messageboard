package model

import (
	"time"
)

type Vote struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_comment"`
	CommentID uint      `json:"comment_id" gorm:"not null;uniqueIndex:idx_user_comment"`
	Value     int       `json:"value" gorm:"not null"` // 1 for upvote, -1 for downvote
	CreatedAt time.Time `json:"created_at"`
}

func (Vote) TableName() string {
	return "votes"
}