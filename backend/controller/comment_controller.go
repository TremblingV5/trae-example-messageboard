package controller

import (
	"messageboard/dto/request"
	"messageboard/dto/response"
	"messageboard/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentService *service.CommentService
	voteService    *service.VoteService
}

func NewCommentController(commentService *service.CommentService, voteService *service.VoteService) *CommentController {
	return &CommentController{
		commentService: commentService,
		voteService:    voteService,
	}
}

// CreateComment 创建评论
func (c *CommentController) CreateComment(ctx *gin.Context) {
	postIDStr := ctx.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid post id"))
		return
	}

	var req request.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	userID := ctx.GetUint("user_id")

	comment, err := c.commentService.CreateComment(uint(postID), userID, req.Content, req.ParentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse(gin.H{
		"id":         comment.ID,
		"content":    comment.Content,
		"post_id":    comment.PostID,
		"author_id":  comment.AuthorID,
		"parent_id":  comment.ParentID,
		"created_at": comment.CreatedAt,
	}))
}

// GetComments 获取评论树
func (c *CommentController) GetComments(ctx *gin.Context) {
	postIDStr := ctx.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid post id"))
		return
	}

	// Get current user ID (may be 0 if not authenticated)
	userID, _ := ctx.Get("user_id")
	var uid uint
	if userID != nil {
		uid = userID.(uint)
	}

	comments, err := c.commentService.GetCommentTree(uint(postID), uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(comments))
}

// Vote 点赞/取消点赞
func (c *CommentController) Vote(ctx *gin.Context) {
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid comment id"))
		return
	}

	var req request.VoteActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	userID := ctx.GetUint("user_id")

	voteResp, err := c.voteService.Vote(userID, uint(commentID), req.Action)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(voteResp))
}

// GetVotes 获取点赞数
func (c *CommentController) GetVotes(ctx *gin.Context) {
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid comment id"))
		return
	}

	// Get current user ID (may be 0 if not authenticated)
	userID, _ := ctx.Get("user_id")
	var uid uint
	if userID != nil {
		uid = userID.(uint)
	}

	voteResp, err := c.voteService.GetVoteInfo(uint(commentID), uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(voteResp))
}