package cache

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

// Cache 定义缓存系统需要提供的基础能力。
//
// Service 层只依赖这个接口，
// 不需要知道底层究竟是 Redis、内存缓存还是其他缓存系统。
type Cache interface {

	// Get 获取缓存数据。
	//
	// key：
	//   缓存唯一标识。
	//
	// 返回：
	//   value：缓存内容。
	//   err：读取缓存过程中出现的错误。
	Get(
		ctx context.Context,
		key string,
	) (string, error)

	// Set 写入缓存。
	//
	// ttl：
	//   缓存有效时间。
	Set(
		ctx context.Context,
		key string,
		value string,
		ttl time.Duration,
	) error

	// Delete 删除指定缓存。
	Delete(
		ctx context.Context,
		key string,
	) error
}
