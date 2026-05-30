package middleware

import (
	"messageboard/dto/response"
	"messageboard/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// parseToken 解析 JWT token，返回 user_id 和是否解析成功
func parseToken(authHeader string) (uint, bool) {
	// Check Bearer prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, false
	}

	tokenString := parts[1]

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(model.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		return 0, false
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	// Get user_id from claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, false
	}

	return uint(userIDFloat), true
}

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "authorization header is required"))
			ctx.Abort()
			return
		}

		userID, ok := parseToken(authHeader)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "invalid token"))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

// OptionalAuth 可选的认证中间件（不强制要求认证）
func OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}

		userID, ok := parseToken(authHeader)
		if ok {
			ctx.Set("user_id", userID)
		}

		ctx.Next()
	}
}
