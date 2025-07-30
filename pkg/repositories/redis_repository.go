package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisRepositoryImpl struct {
	client *redis.Client
}

type RedisRepository interface {
	GetClient() *redis.Client
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) error
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepositoryImpl{
		client: client,
	}
}

func (r *redisRepositoryImpl) GetClient() *redis.Client {
	return r.client
}

func (r *redisRepositoryImpl) Set(key string, value interface{}) error {
	ctx := context.Background()
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return errors.New(fmt.Sprint("Please contact our customer service."))
	}
	return nil
}

func (r *redisRepositoryImpl) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *redisRepositoryImpl) Delete(key string) error {
	ctx := context.Background()
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errors.New(fmt.Sprint("Please contact our customer service."))
	}
	return nil
}
