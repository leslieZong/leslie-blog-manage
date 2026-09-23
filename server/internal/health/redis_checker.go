package health

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisChecker 用于检查 Redis 是否正常。
type RedisChecker struct {
	client *redis.Client
}

// NewRedisChecker 创建 Redis Checker。
func NewRedisChecker(
	client *redis.Client,
) *RedisChecker {

	return &RedisChecker{
		client: client,
	}
}

// Check 通过 Redis PING 检查连接。
func (c *RedisChecker) Check(
	ctx context.Context,
) error {

	return c.client.Ping(ctx).Err()
}
