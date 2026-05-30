package middleware

import (
	"messageboard/dto/response"
	"messageboard/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// parseToken 解析 JWT token，返回 claims 和错误
func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(model.GetJWTSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// extractTokenFromHeader 从 Authorization 头中提取 Bearer token
func extractTokenFromHeader(authHeader string) (string, bool) {
	if authHeader == "" {
		return "", false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", false
	}

	return parts[1], true
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

		tokenString, valid := extractTokenFromHeader(authHeader)
		if !valid {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "invalid authorization header format"))
			ctx.Abort()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "invalid token"))
			ctx.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "user_id not found in token"))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", uint(userIDFloat))
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

		tokenString, valid := extractTokenFromHeader(authHeader)
		if !valid {
			ctx.Next()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			ctx.Next()
			return
		}

		ctx.Set("user_id", uint(userIDFloat))
		ctx.Next()
	}
}
