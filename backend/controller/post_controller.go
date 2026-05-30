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
		utils.BadRequest(ctx, "invalid request parameters")
		return
	}

	userID := ctx.GetUint("user_id")

	post, err := c.postService.CreatePost(userID, &req)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to create post", err)
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
		utils.BadRequest(ctx, "invalid request parameters")
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
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to get post list", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(list))
}

// GetPost 获取帖子详情
func (c *PostController) GetPost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid post id")
		return
	}

	post, err := c.postService.GetPostByID(uint(id))
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "post not found", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Author:    response.NewUserResponsePtr(post.Author),
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}))
}

// SearchPosts 搜索帖子
func (c *PostController) SearchPosts(ctx *gin.Context) {
	var req request.SearchPostRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "invalid request parameters")
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
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to search posts", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(list))
}

// UploadPostImage 上传帖子图片
func (c *PostController) UploadPostImage(ctx *gin.Context) {
	// Get file from request
	file, err := ctx.FormFile("image")
	if err != nil {
		utils.BadRequest(ctx, "no file uploaded")
		return
	}

	// Validate file size (5MB max for post image)
	if file.Size > model.MaxImageSize {
		utils.BadRequest(ctx, "file size exceeds 5MB limit")
		return
	}

	// Save file
	imageURL, err := utils.SaveFile(file, "posts")
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to upload image", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(gin.H{
		"image_url": imageURL,
	}))
}

// UpdatePost 更新帖子
func (c *PostController) UpdatePost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid post id")
		return
	}

	var req request.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "invalid request parameters")
		return
	}

	userID := ctx.GetUint("user_id")

	post, err := c.postService.UpdatePost(uint(id), userID, &req)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "failed to update post", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}))
}

// DeletePost 删除帖子
func (c *PostController) DeletePost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid post id")
		return
	}

	userID := ctx.GetUint("user_id")

	if err := c.postService.DeletePost(uint(id), userID); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "failed to delete post", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(nil))
}
