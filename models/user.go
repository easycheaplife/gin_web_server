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

// getCacheKey 获取用户缓存key
func getUserCacheKey(id uint) string {
	return fmt.Sprintf("user:%d", id)
}

// getUsersListCacheKey 获取用户列表缓存key
func getUsersListCacheKey(page, pageSize int) string {
	return fmt.Sprintf("users:list:%d:%d", page, pageSize)
}

// Create 创建用户
func (u *User) Create() error {
	if err := database.GetDB().Create(u).Error; err != nil {
		return err
	}

	// 创建用户后清除用户列表缓存
	if redisClient := database.GetRedis(); redisClient != nil {
		ctx := context.Background()
		pattern := "users:list:*"
		iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			redisClient.Del(ctx, iter.Val())
		}
	}

	return nil
}

// GetUserByID 根据ID获取用户信息
func GetUserByID(id uint) (*User, error) {
	var user User
	redisClient := database.GetRedis()
	cacheKey := getUserCacheKey(id)
	ctx := context.Background()

	// 尝试从Redis缓存获取
	if redisClient != nil {
		data, err := redisClient.Get(ctx, cacheKey).Bytes()
		if err == nil {
			if err := json.Unmarshal(data, &user); err == nil {
				return &user, nil
			}
		}
	}

	// 从数据库读取
	if err := database.GetDB().First(&user, id).Error; err != nil {
		return nil, err
	}

	// 写入Redis缓存
	if redisClient != nil {
		if data, err := json.Marshal(user); err == nil {
			redisClient.Set(ctx, cacheKey, data, time.Hour)
		}
	}

	return &user, nil
}

// GetUsers 获取用户列表
func GetUsers(page, pageSize int) ([]User, error) {
	var users []User
	redisClient := database.GetRedis()
	cacheKey := getUsersListCacheKey(page, pageSize)
	ctx := context.Background()

	// 尝试从Redis缓存获取
	if redisClient != nil {
		data, err := redisClient.Get(ctx, cacheKey).Bytes()
		if err == nil {
			if err := json.Unmarshal(data, &users); err == nil {
				return users, nil
			}
		}
	}

	// 从数据库读取
	offset := (page - 1) * pageSize
	if err := database.GetDB().Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	// 写入Redis缓存
	if redisClient != nil {
		if data, err := json.Marshal(users); err == nil {
			redisClient.Set(ctx, cacheKey, data, 5*time.Minute) // 列表缓存时间较短
		}
	}

	return users, nil
}

// Update 更新用户信息
func (u *User) Update() error {
	// 更新数据库
	if err := database.GetDB().Save(u).Error; err != nil {
		return err
	}

	redisClient := database.GetRedis()
	if redisClient != nil {
		ctx := context.Background()

		// 删除用户详情缓存
		cacheKey := getUserCacheKey(u.ID)
		redisClient.Del(ctx, cacheKey)

		// 删除所有用户列表缓存
		pattern := "users:list:*"
		iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			redisClient.Del(ctx, iter.Val())
		}
	}

	return nil
}

// Delete 删除用户
func (u *User) Delete() error {
	// 删除数据库记录
	if err := database.GetDB().Delete(u).Error; err != nil {
		return err
	}

	// 删除缓存
	redisClient := database.GetRedis()
	if redisClient != nil {
		ctx := context.Background()

		// 删除用户详情缓存
		cacheKey := getUserCacheKey(u.ID)
		redisClient.Del(ctx, cacheKey)

		// 删除所有用户列表缓存
		pattern := "users:list:*"
		iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			redisClient.Del(ctx, iter.Val())
		}
	}

	return nil
}
