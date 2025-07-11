package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/riskibarqy/bq-account-service/config"
)

var ctx = context.Background()

// RedisClient holds the Redis client instance
var RedisClient *redis.Client

// Init initializes the Redis client
func Init() {
	opt, err := redis.ParseURL(config.AppConfig.RedisURL)
	if err != nil {
		// logger.Log(ctx, logger.LevelError, err.Error(), err)
		// err = &types.Error{}
		panic(err)
	}
	RedisClient = redis.NewClient(opt)

	// Test the connection
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		// logger.Log(ctx, logger.LevelError, err.Error(), err)
		panic(err)
	}
}
