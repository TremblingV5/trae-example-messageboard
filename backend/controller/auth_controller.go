package controller

import (
	"messageboard/dto/request"
	"messageboard/dto/response"
	"messageboard/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	user, err := c.authService.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse(response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}))
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, err.Error()))
		return
	}

	loginResp, err := c.authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(loginResp))
}