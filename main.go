package main

import (
	"gin_web_server/api"
	"gin_web_server/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 路由组 v1
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		users := v1.Group("/users")
		{
			users.GET("", api.GetUsers)    // GET /api/v1/users
			users.GET("/:id", api.GetUser) // GET /api/v1/users/:id
			users.POST("", api.CreateUser) // POST /api/v1/users
		}
	}

	r.Run() // 默认监听 :8080
}
