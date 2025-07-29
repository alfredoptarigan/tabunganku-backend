package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisOnce sync.Once
	redisDb   *redis.Client
)

func InitRedis() *redis.Client {

	redisOnce.Do(func() {
		redisDb = redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         fmt.Sprintf("%s:%s", RedisHost, RedisPort),
			Password:     RedisPass,
			DB:           0, // use default DB
			MaxRetries:   3,
			MaxIdleConns: 10,
			MinIdleConns: 5,
		})

		if err := redisDb.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}
	})

	return redisDb
}
