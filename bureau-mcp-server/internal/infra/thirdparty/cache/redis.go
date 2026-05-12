package cache

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type redisAdapter struct {
	rdb *redis.Client
}

func NewRedisAdapter() *redisAdapter {
	adapter := &redisAdapter{}
	adapter.rdb = adapter.Connect()
	return adapter
}

func (m *redisAdapter) Connect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		log.Fatalf("failed to connect to redis: %v", status.Err())
	}

	return rdb
}
