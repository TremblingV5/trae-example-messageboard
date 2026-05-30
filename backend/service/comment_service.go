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
	// 一次性加载帖子所有评论
	allComments, err := s.commentRepo.FindByPostID(postID)
	if err != nil {
		return nil, err
	}

	// 批量查询所有评论的投票信息
	commentIDs := make([]uint, len(allComments))
	for i, c := range allComments {
		commentIDs[i] = c.ID
	}

	voteCounts, err := s.voteRepo.BatchGetVoteCounts(commentIDs)
	if err != nil {
		voteCounts = make(map[uint]int)
	}

	userVotes := make(map[uint]int)
	if currentUserID > 0 {
		userVotes, err = s.voteRepo.BatchGetUserVotes(currentUserID, commentIDs)
		if err != nil {
			userVotes = make(map[uint]int)
		}
	}

	// 在内存中构建树
	commentMap := make(map[uint]*model.Comment)
	for i := range allComments {
		commentMap[allComments[i].ID] = &allComments[i]
	}

	// 构建评论响应映射
	responseMap := make(map[uint]*response.CommentResponse)
	for i := range allComments {
		c := &allComments[i]
		var author *response.UserResponse
		if c.Author != nil {
			author = response.NewUserResponsePtr(c.Author)
		}

		resp := &response.CommentResponse{
			ID:        c.ID,
			Content:   c.Content,
			PostID:    c.PostID,
			Author:    author,
			ParentID:  c.ParentID,
			VoteCount: voteCounts[c.ID],
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
		responseMap[c.ID] = resp
	}

	// 组装树结构
	var roots []response.CommentResponse
	for i := range allComments {
		c := &allComments[i]
		resp := responseMap[c.ID]

		if c.ParentID == nil {
			roots = append(roots, *resp)
		} else if parentResp, ok := responseMap[*c.ParentID]; ok {
			parentResp.Children = append(parentResp.Children, *resp)
		}
	}

	return roots, nil
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
