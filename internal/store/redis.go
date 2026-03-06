package store

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisStore 实现了 Store 接口，使用 Redis 作为存储后端
// Redis 具有高性能、持久化（AOF/RDB）、支持过期时间等特性，非常适合短链接服务。
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore 创建一个新的 RedisStore 实例
// addr: Redis 地址，例如 "localhost:6379"
// password: Redis 密码，无密码传空字符串
// db: Redis 数据库索引，默认 0
func NewRedisStore(addr string, password string, db int) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisStore{
		client: client,
	}
}

// Save 实现 Store 接口的 Save 方法
func (s *RedisStore) Save(shortCode string, originalURL string) error {
	ctx := context.Background()
	// 使用 SET 命令保存键值对
	// 0 表示不过期（永久保存）。
	// 如果需要设置过期时间（例如 7 天），可以把 0 改为 7 * 24 * time.Hour
	err := s.client.Set(ctx, shortCode, originalURL, 0).Err()
	return err
}

// Load 实现 Store 接口的 Load 方法
func (s *RedisStore) Load(shortCode string) (string, error) {
	ctx := context.Background()
	// 使用 GET 命令获取值
	url, err := s.client.Get(ctx, shortCode).Result()
	if err == redis.Nil {
		// redis.Nil 表示 key 不存在
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}
	return url, nil
}
