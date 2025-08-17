package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GetUsers 返回用户列表
func GetUsers(c *gin.Context) {
	users := []User{
		{ID: 1, Username: "user1", Email: "user1@example.com"},
		{ID: 2, Username: "user2", Email: "user2@example.com"},
	}
	c.JSON(http.StatusOK, users)
}

// GetUser 返回单个用户信息
func GetUser(c *gin.Context) {
	// 获取用户ID参数
	_ = c.Param("id") // 在实际应用中，这里应该使用userID查询数据库

	// 这里演示返回模拟数据
	user := User{ID: 1, Username: "user1", Email: "user1@example.com"}
	c.JSON(http.StatusOK, user)
}

// CreateUser 创建新用户
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 这里应该添加数据库操作
	c.JSON(http.StatusCreated, user)
}
