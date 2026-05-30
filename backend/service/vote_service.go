package service

import (
	"errors"
	"messageboard/dto/response"
	"messageboard/model"
	"messageboard/repository"

	"gorm.io/gorm"
)

type VoteService struct {
	voteRepo    *repository.VoteRepository
	commentRepo *repository.CommentRepository
}

func NewVoteService(voteRepo *repository.VoteRepository, commentRepo *repository.CommentRepository) *VoteService {
	return &VoteService{
		voteRepo:    voteRepo,
		commentRepo: commentRepo,
	}
}

// Vote handles voting on a comment
// action: "upvote", "downvote", "cancel"
func (s *VoteService) Vote(userID, commentID uint, action string) (*response.VoteResponse, error) {
	// Check if comment exists
	_, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	// Get existing vote
	existingVote, err := s.voteRepo.FindByUserAndComment(userID, commentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	switch action {
	case "upvote":
		if existingVote != nil {
			if existingVote.Value == 1 {
				// Already upvoted, cancel
				if err := s.voteRepo.Delete(existingVote.ID); err != nil {
					return nil, err
				}
			} else {
				// Change from downvote to upvote
				existingVote.Value = 1
				if err := s.voteRepo.Update(existingVote); err != nil {
					return nil, err
				}
			}
		} else {
			// Create new upvote
			vote := &model.Vote{
				UserID:    userID,
				CommentID: commentID,
				Value:     1,
			}
			if err := s.voteRepo.Create(vote); err != nil {
				return nil, err
			}
		}

	case "downvote":
		if existingVote != nil {
			if existingVote.Value == -1 {
				// Already downvoted, cancel
				if err := s.voteRepo.Delete(existingVote.ID); err != nil {
					return nil, err
				}
			} else {
				// Change from upvote to downvote
				existingVote.Value = -1
				if err := s.voteRepo.Update(existingVote); err != nil {
					return nil, err
				}
			}
		} else {
			// Create new downvote
			vote := &model.Vote{
				UserID:    userID,
				CommentID: commentID,
				Value:     -1,
			}
			if err := s.voteRepo.Create(vote); err != nil {
				return nil, err
			}
		}

	case "cancel":
		if existingVote != nil {
			if err := s.voteRepo.Delete(existingVote.ID); err != nil {
				return nil, err
			}
		}

	default:
		return nil, errors.New("invalid action")
	}

	// Get updated vote count
	voteCount, err := s.voteRepo.GetVoteCount(commentID)
	if err != nil {
		return nil, err
	}

	// Check current user's vote status
	currentValue, err := s.voteRepo.GetUserVote(userID, commentID)
	if err != nil {
		return nil, err
	}

	voted := currentValue != 0

	return &response.VoteResponse{
		VoteCount: voteCount,
		Voted:     voted,
		Value:     currentValue,
	}, nil
}

// GetVoteInfo gets vote information for a comment
func (s *VoteService) GetVoteInfo(commentID uint, userID uint) (*response.VoteResponse, error) {
	voteCount, err := s.voteRepo.GetVoteCount(commentID)
	if err != nil {
		return nil, err
	}

	value, err := s.voteRepo.GetUserVote(userID, commentID)
	if err != nil {
		return nil, err
	}

	voted := value != 0

	return &response.VoteResponse{
		VoteCount: voteCount,
		Voted:     voted,
		Value:     value,
	}, nil
}