package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/riskibarqy/bq-account-service/config"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
)

var RedisClient CacheClient

// Init initializes the Redis client
func Init(ctx context.Context) {
	opt, err := redis.ParseURL(config.AppConfig.RedisURL)
	if err != nil {
		// logger.Log(ctx, logger.LevelError, err.Error(), err)
		// err = &types.Error{}
		panic(err)
	}
	RedisClient = redistrace.NewClient(opt)

	// Test the connection
	ping, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		// logger.Log(ctx, logger.LevelError, err.Error(), err)
		panic(err)
	}
	// config.AppConfig.RedisClient = &RedisClient

	log.Printf("[Redis] %s", ping)
}
