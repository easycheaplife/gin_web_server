package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 是一个简单的日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		t := time.Now()
		path := c.Request.URL.Path

		// 处理请求
		c.Next()

		// 请求后
		latency := time.Since(t)
		statusCode := c.Writer.Status()

		fmt.Printf("[%s] %s %d %v\n",
			path,
			c.Request.Method,
			statusCode,
			latency,
		)
	}
}
