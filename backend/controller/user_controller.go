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
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	user, err := c.authService.GetUserByID(uint(id))
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "user not found", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(response.NewUserResponse(user)))
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	// Get current user from context
	userID := ctx.GetUint("user_id")
	if userID != uint(id) {
		utils.Forbidden(ctx, "cannot update other user's profile")
		return
	}

	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "invalid request parameters")
		return
	}

	user, err := c.authService.UpdateUser(uint(id), &req)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(response.NewUserResponse(user)))
}

// UploadAvatar 上传头像
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	// Get current user from context
	userID := ctx.GetUint("user_id")
	if userID != uint(id) {
		utils.Forbidden(ctx, "cannot update other user's avatar")
		return
	}

	// Get file from request
	file, err := ctx.FormFile("avatar")
	if err != nil {
		utils.BadRequest(ctx, "no file uploaded")
		return
	}

	// Validate file size (2MB max for avatar)
	if file.Size > model.MaxAvatarSize {
		utils.BadRequest(ctx, "file size exceeds 2MB limit")
		return
	}

	// Save file
	avatarURL, err := utils.SaveFile(file, "avatars")
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to upload avatar", err)
		return
	}

	user, err := c.authService.UpdateAvatar(uint(id), avatarURL)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "failed to update avatar", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(gin.H{
		"avatar_url": avatarURL,
		"user":       response.NewUserResponse(user),
	}))
}
