package models

import (
"context"
"encoding/json"
"fmt"
"gin_web_server/database"
"time"

"gorm.io/gorm"
)

type User struct {
ID        uint           `json:"id" gorm:"primarykey"`
Username  string         `json:"username" gorm:"uniqueIndex;not null"`
Email     string         `json:"email" gorm:"uniqueIndex;not null"`
CreatedAt time.Time      `json:"created_at"`
UpdatedAt time.Time      `json:"updated_at"`
DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
return "users"
}

// Create 创建用户
func (u *User) Create() error {
return database.GetDB().Create(u).Error
}

// GetUserByID 根据ID获取用户信息
func GetUserByID(id uint) (*User, error) {
var user User

// 尝试从Redis缓存获取（如果Redis可用）
if redisClient := database.GetRedis(); redisClient != nil {
cacheKey := fmt.Sprintf("user:%d", id)
ctx := context.Background()

// 从Redis读取
data, err := redisClient.Get(ctx, cacheKey).Bytes()
if err == nil {
// 缓存命中，解析数据
if err := json.Unmarshal(data, &user); err == nil {
return &user, nil
}
}
}

// 从数据库读取
if err := database.GetDB().First(&user, id).Error; err != nil {
return nil, err
}

// 写入Redis缓存（如果Redis可用）
if redisClient := database.GetRedis(); redisClient != nil {
if data, err := json.Marshal(user); err == nil {
cacheKey := fmt.Sprintf("user:%d", id)
ctx := context.Background()
redisClient.Set(ctx, cacheKey, data, time.Hour)
}
}

return &user, nil
}

// GetUsers 获取用户列表
func GetUsers(page, pageSize int) ([]User, error) {
var users []User
offset := (page - 1) * pageSize

err := database.GetDB().Offset(offset).Limit(pageSize).Find(&users).Error
return users, err
}

// Update 更新用户信息
func (u *User) Update() error {
// 更新数据库
if err := database.GetDB().Save(u).Error; err != nil {
return err
}

// 删除缓存，让其重新生成（如果Redis可用）
if redisClient := database.GetRedis(); redisClient != nil {
cacheKey := fmt.Sprintf("user:%d", u.ID)
ctx := context.Background()
redisClient.Del(ctx, cacheKey)
}

return nil
}

// Delete 删除用户
func (u *User) Delete() error {
// 删除数据库记录
if err := database.GetDB().Delete(u).Error; err != nil {
return err
}

// 删除缓存（如果Redis可用）
if redisClient := database.GetRedis(); redisClient != nil {
cacheKey := fmt.Sprintf("user:%d", u.ID)
ctx := context.Background()
redisClient.Del(ctx, cacheKey)
}

return nil
}
