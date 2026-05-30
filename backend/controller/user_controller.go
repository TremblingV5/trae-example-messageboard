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

type UserController struct {
	authService *service.AuthService
}

func NewUserController(authService *service.AuthService) *UserController {
	return &UserController{authService: authService}
}

// GetUser 获取用户信息
func (c *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid user id"))
		return
	}

	user, err := c.authService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse(404, "user not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}))
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid user id"))
		return
	}

	// Get current user from context
	userID := ctx.GetUint("user_id")
	if userID != uint(id) {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse(403, "cannot update other user's profile"))
		return
	}

	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	user, err := c.authService.UpdateUser(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}))
}

// UploadAvatar 上传头像
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "invalid user id"))
		return
	}

	// Get current user from context
	userID := ctx.GetUint("user_id")
	if userID != uint(id) {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse(403, "cannot update other user's avatar"))
		return
	}

	// Get file from request
	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "no file uploaded"))
		return
	}

	// Validate file size (2MB max for avatar)
	if file.Size > model.MaxAvatarSize {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, "file size exceeds 2MB limit"))
		return
	}

	// Save file
	avatarURL, err := utils.SaveFile(file, "avatars")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	user, err := c.authService.UpdateAvatar(uint(id), avatarURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(gin.H{
		"avatar_url": avatarURL,
		"user": response.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}))
}