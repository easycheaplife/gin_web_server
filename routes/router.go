package routes

import (
	"gin_web_server/api"
	"gin_web_server/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置所有路由
func SetupRouter() *gin.Engine {
	// 初始化Gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})

	// API 路由组 v1
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", api.GetUsers)          // GET /api/v1/users
			users.GET("/:id", api.GetUser)       // GET /api/v1/users/:id
			users.POST("", api.CreateUser)       // POST /api/v1/users
			users.PUT("/:id", api.UpdateUser)    // PUT /api/v1/users/:id
			users.DELETE("/:id", api.DeleteUser) // DELETE /api/v1/users/:id
		}
	}

	return r
}
