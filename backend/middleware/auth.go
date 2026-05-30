package middleware

import (
	"messageboard/dto/response"
	"messageboard/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "authorization header is required"))
			ctx.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "invalid authorization header format"))
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(model.GetJWTSecret()), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "invalid token"))
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "token is invalid"))
			ctx.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "invalid token claims"))
			ctx.Abort()
			return
		}

		// Get user_id from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(401, "user_id not found in token"))
			ctx.Abort()
			return
		}

		userID := uint(userIDFloat)
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

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Next()
			return
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
			ctx.Next()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.Next()
			return
		}

		// Get user_id from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			ctx.Next()
			return
		}

		userID := uint(userIDFloat)
		ctx.Set("user_id", userID)

		ctx.Next()
	}
}