package controller

import (
	"messageboard/dto/request"
	"messageboard/dto/response"
	"messageboard/model"
	"messageboard/service"
	"messageboard/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService *service.PostService
}

func NewPostController(postService *service.PostService) *PostController {
	return &PostController{postService: postService}
}

// CreatePost 创建帖子
func (c *PostController) CreatePost(ctx *gin.Context) {
	var req request.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	userID := ctx.GetUint("user_id")

	post, err := c.postService.CreatePost(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse(response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}))
}

// GetPostList 获取帖子列表
func (c *PostController) GetPostList(ctx *gin.Context) {
	var req request.PostListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	list, err := c.postService.GetPostList(req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(list))
}

// GetPost 获取帖子详情
func (c *PostController) GetPost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid post id"))
		return
	}

	post, err := c.postService.GetPostByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse(404, "post not found"))
		return
	}

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

	ctx.JSON(http.StatusOK, response.SuccessResponse(response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Author:    author,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}))
}

// SearchPosts 搜索帖子
func (c *PostController) SearchPosts(ctx *gin.Context) {
	var req request.SearchPostRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	list, err := c.postService.SearchPosts(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(list))
}

// UploadPostImage 上传帖子图片
func (c *PostController) UploadPostImage(ctx *gin.Context) {
	// Get file from request
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "no file uploaded"))
		return
	}

	// Validate file size (5MB max for post image)
	if file.Size > model.MaxImageSize {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "file size exceeds 5MB limit"))
		return
	}

	// Save file
	imageURL, err := utils.SaveFile(file, "posts")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(gin.H{
		"image_url": imageURL,
	}))
}