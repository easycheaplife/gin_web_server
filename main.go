package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"gin_web_server/api"
	"gin_web_server/config"
	"gin_web_server/database"
	"gin_web_server/middleware"
	"gin_web_server/models"
)

var cfg *config.Config

func init() {
	// 加载配置
	var err error
	cfg, err = loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	if err := database.InitMySQL(cfg.MySQL); err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}

	// 初始化Redis连接（可选）
	if err := database.InitRedis(cfg.Redis); err != nil {
		log.Printf("Warning: Redis initialization failed: %v", err)
		log.Println("Continuing without Redis cache...")
	}

	// 自动迁移数据库表
	if err := database.GetDB().AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func main() {
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loadConfig 加载配置文件
func loadConfig() (*config.Config, error) {
	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}

	// 读取配置文件
	data, err := os.ReadFile(pwd + "/config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 替换环境变量
	if mysqlPassword := os.Getenv("MYSQL_PASSWORD"); mysqlPassword != "" {
		data = []byte(strings.Replace(string(data), "${MYSQL_PASSWORD}", mysqlPassword, 1))
	} else {
		return nil, fmt.Errorf("MYSQL_PASSWORD environment variable is not set")
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// 打印配置信息
	log.Printf("Config loaded: MySQL(%s:%d), Redis(%s:%d)",
		cfg.MySQL.Host, cfg.MySQL.Port,
		cfg.Redis.Host, cfg.Redis.Port)

	return &cfg, nil
}
