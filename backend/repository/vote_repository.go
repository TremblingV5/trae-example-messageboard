package repository

import (
	"errors"
	"messageboard/model"

	"gorm.io/gorm"
)

type VoteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) Create(vote *model.Vote) error {
	return r.db.Create(vote).Error
}

func (r *VoteRepository) FindByUserAndComment(userID, commentID uint) (*model.Vote, error) {
	var vote model.Vote
	err := r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&vote).Error
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) Update(vote *model.Vote) error {
	return r.db.Save(vote).Error
}

func (r *VoteRepository) Delete(id uint) error {
	return r.db.Delete(&model.Vote{}, id).Error
}

func (r *VoteRepository) GetVoteCount(commentID uint) (int, error) {
	var result struct {
		Total int
	}
	err := r.db.Model(&model.Vote{}).
		Select("COALESCE(SUM(value), 0) as total").
		Where("comment_id = ?", commentID).
		Scan(&result).Error
	return result.Total, err
}

// BatchGetVoteCounts 批量获取多个评论的投票数
func (r *VoteRepository) BatchGetVoteCounts(commentIDs []uint) (map[uint]int, error) {
	result := make(map[uint]int)
	if len(commentIDs) == 0 {
		return result, nil
	}

	type VoteCount struct {
		CommentID uint `gorm:"column:comment_id"`
		Total     int  `gorm:"column:total"`
	}

	var counts []VoteCount
	err := r.db.Model(&model.Vote{}).
		Select("comment_id, COALESCE(SUM(value), 0) as total").
		Where("comment_id IN ?", commentIDs).
		Group("comment_id").
		Find(&counts).Error
	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		result[c.CommentID] = c.Total
	}
	return result, nil
}

// BatchGetUserVotes 批量获取用户对多个评论的投票状态
func (r *VoteRepository) BatchGetUserVotes(userID uint, commentIDs []uint) (map[uint]int, error) {
	result := make(map[uint]int)
	if len(commentIDs) == 0 {
		return result, nil
	}

	var votes []model.Vote
	err := r.db.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&votes).Error
	if err != nil {
		return nil, err
	}

	for _, v := range votes {
		result[v.CommentID] = v.Value
	}
	return result, nil
}

// GetUserVote 获取用户对评论的投票状态
func (r *VoteRepository) GetUserVote(userID, commentID uint) (int, error) {
	vote, err := r.FindByUserAndComment(userID, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // 用户未投票
		}
		return 0, err
	}
	return vote.Value, nil
}
