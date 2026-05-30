package middleware

import (
	"messageboard/config"

	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	allowedOrigins := config.AppConfig.AllowedOrigins
	isProduction := config.AppConfig.Server.Mode == "release"

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")

		if isProduction {
			// Production: only allow configured origins, with credentials
			allowed := false
			for _, o := range allowedOrigins {
				if o == origin {
					allowed = true
					break
				}
			}
			if allowed {
				ctx.Header("Access-Control-Allow-Origin", origin)
				ctx.Header("Access-Control-Allow-Credentials", "true")
			}
		} else {
			// Development: echo back the request Origin instead of using wildcard
			if origin != "" {
				ctx.Header("Access-Control-Allow-Origin", origin)
				ctx.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		// Handle preflight request
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
