package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return param.TimeStamp.Format("2006/01/02 - 15:04:05") +
			" | " + param.ClientIP +
			" | " + param.Method +
			" | " + param.Path +
			" | " + string(rune(param.StatusCode)) +
			" | " + param.Latency.String() +
			" | " + param.ErrorMessage + "\n"
	})
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}