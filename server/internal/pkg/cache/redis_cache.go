package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache 是 Cache 接口的 Redis 实现。
type redisCache struct {
	client *redis.Client
}

// NewRedisCache 创建 Redis Cache。
func NewRedisCache(
	client *redis.Client,
) Cache {

	return &redisCache{
		client: client,
	}
}

// Get 从 Redis 获取缓存。
func (c *redisCache) Get(
	ctx context.Context,
	key string,
) (string, error) {

	value, err := c.client.Get(
		ctx,
		key,
	).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", err
		}

		return "", err
	}

	return value, nil
}

// Set 写入 Redis。
func (c *redisCache) Set(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {

	return c.client.Set(
		ctx,
		key,
		value,
		ttl,
	).Err()
}

// Delete 删除 Redis 缓存。
func (c *redisCache) Delete(
	ctx context.Context,
	key string,
) error {

	return c.client.Del(
		ctx,
		key,
	).Err()
}
