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

	// 使用 Set + NX 选项原子性地保存
	// NX: 仅在 key 不存在时设置，避免竞态条件
	err := s.client.SetNX(ctx, shortCode, originalURL, 0).Err()
	if err == redis.Nil {
		return ErrExists
	}
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

// IncrementVisits 增加访问计数
func (s *RedisStore) IncrementVisits(shortCode string) error {
	ctx := context.Background()
	// 使用 INCR 命令原子递增计数器
	// key 为 "visits:" + shortCode，避免与 URL 映射冲突
	key := "visits:" + shortCode
	return s.client.Incr(ctx, key).Err()
}

// GetVisits 获取访问计数
func (s *RedisStore) GetVisits(shortCode string) (int64, error) {
	ctx := context.Background()
	key := "visits:" + shortCode
	val, err := s.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return val, nil
}

// Delete 删除短链接
func (s *RedisStore) Delete(shortCode string) error {
	ctx := context.Background()

	// Del 返回被删除的 key 数量
	deleted, err := s.client.Del(ctx, shortCode).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return ErrNotFound
	}

	// 删除访问计数（忽略错误可能不存在）
	s.client.Del(ctx, "visits:"+shortCode)
	return nil
}
