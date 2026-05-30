package service

import (
	"messageboard/dto/request"
	"messageboard/dto/response"
	"messageboard/model"
	"messageboard/repository"

	"gorm.io/gorm"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(authorID uint, req *request.CreatePostRequest) (*model.Post, error) {
	post := &model.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
		ImageURL: req.ImageURL,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	// Load author
	return s.postRepo.FindByID(post.ID)
}

func (s *PostService) GetPostByID(id uint) (*model.Post, error) {
	return s.postRepo.FindByID(id)
}

func (s *PostService) GetPostList(page, pageSize int) (*response.PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := s.postRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	postResponses := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = s.toPostResponse(&post)
	}

	return &response.PostListResponse{
		Posts:    postResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *PostService) SearchPosts(keyword string, page, pageSize int) (*response.PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := s.postRepo.Search(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	postResponses := make([]response.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = s.toPostResponse(&post)
	}

	return &response.PostListResponse{
		Posts:    postResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *PostService) UpdatePost(id uint, authorID uint, req *request.CreatePostRequest) (*model.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != authorID {
		return nil, gorm.ErrRecordNotFound
	}

	post.Title = req.Title
	post.Content = req.Content
	post.ImageURL = req.ImageURL

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	return s.postRepo.FindByID(id)
}

func (s *PostService) DeletePost(id uint, authorID uint) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if post.AuthorID != authorID {
		return gorm.ErrRecordNotFound
	}

	return s.postRepo.Delete(id)
}

func (s *PostService) toPostResponse(post *model.Post) response.PostResponse {
	var author *response.UserResponse
	if post.Author != nil {
		author = &response.UserResponse{
			ID:        post.Author.ID,
			Username:  post.Author.Username,
			Nickname:  post.Author.Nickname,
			Avatar:    post.Author.Avatar,
			CreatedAt: post.Author.CreatedAt,
			UpdatedAt: post.Author.UpdatedAt,
		}
	}

	return response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Author:    author,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}