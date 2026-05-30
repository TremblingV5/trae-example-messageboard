package service

import (
	"errors"
	"messageboard/dto/response"
	"messageboard/model"
	"messageboard/repository"

	"gorm.io/gorm"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	voteRepo    *repository.VoteRepository
	postRepo    *repository.PostRepository
}

func NewCommentService(commentRepo *repository.CommentRepository, voteRepo *repository.VoteRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		voteRepo:    voteRepo,
		postRepo:    postRepo,
	}
}

func (s *CommentService) CreateComment(postID, authorID uint, content string, parentID *uint) (*model.Comment, error) {
	// Check if post exists
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Validate parent comment depth if parent_id is provided
	if parentID != nil {
		// Check if parent comment exists
		parentComment, err := s.commentRepo.FindByID(*parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("parent comment not found")
			}
			return nil, err
		}

		// Check if parent comment belongs to the same post
		if parentComment.PostID != postID {
			return nil, errors.New("parent comment does not belong to this post")
		}

		// Check depth limit
		depth, err := s.commentRepo.GetCommentDepth(*parentID)
		if err != nil {
			return nil, err
		}
		if depth >= model.MaxCommentDepth-1 {
			return nil, errors.New("maximum comment depth reached")
		}
	}

	comment := &model.Comment{
		Content:  content,
		PostID:   postID,
		AuthorID: authorID,
		ParentID: parentID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return s.commentRepo.FindByID(comment.ID)
}

func (s *CommentService) GetCommentTree(postID uint, currentUserID uint) ([]response.CommentResponse, error) {
	// Get all comments for this post in one query
	allComments, err := s.commentRepo.FindByPostID(postID)
	if err != nil {
		return nil, err
	}

	if len(allComments) == 0 {
		return []response.CommentResponse{}, nil
	}

	// Collect all comment IDs for batch vote queries
	commentIDs := make([]uint, len(allComments))
	for i, c := range allComments {
		commentIDs[i] = c.ID
	}

	// Batch query vote counts and user votes
	voteCountMap, _ := s.voteRepo.BatchGetVoteCount(commentIDs)
	userVoteMap, _ := s.voteRepo.BatchGetUserVotes(currentUserID, commentIDs)

	// Build comment lookup map
	commentMap := make(map[uint]*model.Comment, len(allComments))
	for i := range allComments {
		commentMap[allComments[i].ID] = &allComments[i]
	}

	// Build tree in memory: only keep root comments (parent_id == nil)
	commentTree := make([]response.CommentResponse, 0)
	for i := range allComments {
		if allComments[i].ParentID == nil {
			commentTree = append(commentTree, s.buildCommentResponseBatch(&allComments[i], currentUserID, commentMap, voteCountMap, userVoteMap))
		}
	}

	return commentTree, nil
}

func (s *CommentService) buildCommentResponseBatch(
	comment *model.Comment,
	currentUserID uint,
	commentMap map[uint]*model.Comment,
	voteCountMap map[uint]int,
	userVoteMap map[uint]int,
) response.CommentResponse {
	var author *response.UserResponse
	if comment.Author != nil {
		author = response.NewUserResponsePtr(comment.Author)
	}

	// Get children from the flat map
	childResponses := make([]response.CommentResponse, 0)
	for _, c := range commentMap {
		if c.ParentID != nil && *c.ParentID == comment.ID {
			childResponses = append(childResponses, s.buildCommentResponseBatch(c, currentUserID, commentMap, voteCountMap, userVoteMap))
		}
	}

	return response.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		PostID:    comment.PostID,
		Author:    author,
		ParentID:  comment.ParentID,
		Children:  childResponses,
		VoteCount: voteCountMap[comment.ID],
		VotedValue: userVoteMap[comment.ID],
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func (s *CommentService) DeleteComment(id uint, authorID uint) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return err
	}

	if comment.AuthorID != authorID {
		return errors.New("unauthorized")
	}

	return s.commentRepo.Delete(id)
}
