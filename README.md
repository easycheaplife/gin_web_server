# Gin Web Server

本项目基于 Gin 框架，提供一个结构化的 Go Web Server 示例，集成了 MySQL 数据库和 Redis 缓存。

## 特性

- 基于 Gin 的 RESTful API
- MySQL 数据持久化存储（使用 GORM）
- Redis 缓存支持（可选）
- 环境变量配置支持
- 自定义日志中间件
- 结构化项目组织

## 项目结构

```
gin_web_server/
├── api/            # API 处理函数
├── config/         # 配置文件
├── database/       # 数据库和缓存连接
├── middleware/     # 中间件
├── models/         # 数据模型
├── main.go        # 主程序入口
├── go.mod         # Go 模块文件
└── README.md      # 项目文档
```

## 环境要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+（可选）

## 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/easycheaplife/gin_web_server.git
cd gin_web_server
```

### 2. 配置环境
1. 准备数据库：
```bash
mysql -u root -p
CREATE DATABASE gin_web_server;
```

2. 设置环境变量：
```bash
export MYSQL_PASSWORD=your_password
```

3. 检查配置文件 `config/config.yaml`：
```yaml
mysql:
  host: localhost
  port: 3306
  user: root
  password: ${MYSQL_PASSWORD}  # 从环境变量读取
  database: gin_web_server

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

### 3. 启动服务
```bash
go mod tidy  # 安装依赖
go run main.go
```

## API 接口

### 用户管理接口

| 方法   | 路径              | 描述         | 参数                    | 示例命令 |
|--------|-------------------|--------------|------------------------|----------|
| GET    | /api/v1/users    | 获取用户列表 | page, page_size       | `curl "http://localhost:8080/api/v1/users?page=1&page_size=10"` |
| GET    | /api/v1/users/:id| 获取用户详情 | -                     | `curl http://localhost:8080/api/v1/users/1` |
| POST   | /api/v1/users    | 创建用户     | username, email       | `curl -X POST http://localhost:8080/api/v1/users -H "Content-Type: application/json" -d '{"username":"test","email":"test@example.com"}'` |
| PUT    | /api/v1/users/:id| 更新用户     | username, email       | `curl -X PUT http://localhost:8080/api/v1/users/1 -H "Content-Type: application/json" -d '{"username":"updated","email":"new@example.com"}'` |
| DELETE | /api/v1/users/:id| 删除用户     | -                     | `curl -X DELETE http://localhost:8080/api/v1/users/1` |

### 响应格式

所有接口返回 JSON 格式数据，基本结构如下：

#### 成功响应
```json
{
    "id": 1,
    "username": "test",
    "email": "test@example.com",
    "created_at": "2025-08-17T19:53:29Z",
    "updated_at": "2025-08-17T19:53:29Z"
}
```

#### 错误响应
```json
{
    "error": "错误信息描述"
}
```
```

## 功能特点

### 1. 环境变量配置
- 支持通过环境变量配置敏感信息
- 便于在不同环境下部署
- 遵循 12-Factor App 原则

### 2. 数据库集成
- 使用 GORM 进行数据库操作
- 支持数据模型的自动迁移
- 提供基础的 CRUD 操作封装

### 3. Redis 缓存（可选）
- 支持数据缓存
- 可配置的过期时间
- 自动缓存清理

### 4. 项目最佳实践
- RESTful API 设计
- 中间件实现的访问日志
- 结构化的错误处理
- 统一的响应格式

## 常见问题

### 1. 数据库连接失败
检查以下几点：
- 环境变量 `MYSQL_PASSWORD` 是否正确设置
- MySQL 服务是否运行
- 数据库用户权限是否配置正确

### 2. Redis 集成问题
Redis 是可选的，如果不需要可以在配置文件中禁用：
```yaml
redis:
  enabled: false
```

## 开发计划

- [ ] 添加用户认证
- [ ] 实现 API 限流
- [ ] 添加单元测试
- [ ] 支持 Docker 部署
   ```

3. 启动服务：
   ```bash
   go run main.go
   ```

## 依赖
- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [GORM](https://gorm.io) - ORM 框架
- [go-redis](https://github.com/redis/go-redis) - Redis 客户端
