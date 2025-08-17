# Gin Web Server

本项目基于 Gin 框架，提供一个结构化的 Go Web Server 示例。

## 项目结构

```
gin_web_server/
├── api/            # API 处理函数
├── middleware/     # 中间件
├── main.go        # 主程序入口
├── go.mod         # Go 模块文件
└── README.md      # 项目文档
```

## 功能特性

- 结构化的项目组织
- 自定义日志中间件
- RESTful API 示例
- 路由分组管理

## 快速开始

1. 安装依赖：
   ```bash
   go mod tidy
   ```
2. 启动服务：
   ```bash
   go run main.go
   ```

## API 接口

### 健康检查
- `GET /ping` - 服务健康检查
  ```bash
  curl http://localhost:8080/ping
  ```

### 用户管理
- `GET /api/v1/users` - 获取用户列表
  ```bash
  curl http://localhost:8080/api/v1/users
  ```
- `GET /api/v1/users/:id` - 获取指定用户
  ```bash
  curl http://localhost:8080/api/v1/users/1
  ```
- `POST /api/v1/users` - 创建新用户
  ```bash
  curl -X POST http://localhost:8080/api/v1/users \
    -H "Content-Type: application/json" \
    -d '{"username":"test","email":"test@example.com"}'
  ```

## 依赖
- [Gin](https://github.com/gin-gonic/gin) - Web 框架
