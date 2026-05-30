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
	// Get all root comments
	rootComments, err := s.commentRepo.FindRootCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	// Build comment tree
	commentTree := make([]response.CommentResponse, 0, len(rootComments))
	for _, comment := range rootComments {
		commentResp := s.buildCommentResponse(&comment, currentUserID)
		commentTree = append(commentTree, commentResp)
	}

	return commentTree, nil
}

func (s *CommentService) buildCommentResponse(comment *model.Comment, currentUserID uint) response.CommentResponse {
	// Get vote count
	voteCount, _ := s.voteRepo.GetVoteCount(comment.ID)

	// Check if current user has voted
	votedValue, _ := s.voteRepo.GetUserVote(currentUserID, comment.ID)

	var author *response.UserResponse
	if comment.Author != nil {
		author = &response.UserResponse{
			ID:        comment.Author.ID,
			Username:  comment.Author.Username,
			Nickname:  comment.Author.Nickname,
			Avatar:    comment.Author.Avatar,
			CreatedAt: comment.Author.CreatedAt,
			UpdatedAt: comment.Author.UpdatedAt,
		}
	}

	// Get children
	children, _ := s.commentRepo.FindChildrenByParentID(comment.ID)
	childResponses := make([]response.CommentResponse, 0, len(children))
	for _, child := range children {
		childResponses = append(childResponses, s.buildCommentResponse(&child, currentUserID))
	}

	return response.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		PostID:    comment.PostID,
		Author:    author,
		ParentID:  comment.ParentID,
		Children:  childResponses,
		VoteCount: voteCount,
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