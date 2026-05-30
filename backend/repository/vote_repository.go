package repository

import (
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

func (r *VoteRepository) DeleteByUserAndComment(userID, commentID uint) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&model.Vote{}).Error
}

// GetUserVote 获取用户对评论的投票状态
func (r *VoteRepository) GetUserVote(userID, commentID uint) (int, error) {
	vote, err := r.FindByUserAndComment(userID, commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // 用户未投票
		}
		return 0, err
	}
	return vote.Value, nil
}