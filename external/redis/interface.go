package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	FlushDB(ctx context.Context) *redis.StatusCmd
	Ping(ctx context.Context) *redis.StatusCmd
}
