package utils

import (
	"messageboard/dto/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success 成功响应
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, response.SuccessResponse(data))
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, response.SuccessWithMessage(message, data))
}

// Created 创建成功响应
func Created(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, response.SuccessResponse(data))
}

// BadRequest 400 错误响应
func BadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusBadRequest, response.ErrorResponse(400, message))
}

// Unauthorized 401 错误响应
func Unauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, message))
}

// Forbidden 403 错误响应
func Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, response.ErrorResponse(403, message))
}

// NotFound 404 错误响应
func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, response.ErrorResponse(404, message))
}

// InternalError 500 错误响应
func InternalError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(500, message))
}