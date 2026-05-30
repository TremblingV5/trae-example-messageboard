package controller

import (
	"messageboard/dto/request"
	"messageboard/dto/response"
	"messageboard/service"
	"messageboard/utils"
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
		utils.BadRequest(ctx, "invalid request parameters")
		return
	}

	user, err := c.authService.Register(&req)
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "registration failed", err)
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse(response.NewUserResponse(user)))
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "invalid request parameters")
		return
	}

	loginResp, err := c.authService.Login(&req)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "login failed", err)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(loginResp))
}
