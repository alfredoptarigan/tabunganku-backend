package services

import "alfredo/tabunganku/pkg/repositories"

type RedisService interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) error
}

type redisServiceImpl struct {
	repository repositories.RedisRepository
}

func (r *redisServiceImpl) Set(key string, value interface{}) error {
	if err := r.repository.Set(key, value); err != nil {
		return err
	}

	return nil
}

func (r *redisServiceImpl) Get(key string) (string, error) {
	res, err := r.repository.Get(key)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (r *redisServiceImpl) Delete(key string) error {
	if err := r.repository.Delete(key); err != nil {
		return err
	}

	return nil
}

func NewRedisService(repository repositories.RedisRepository) RedisService {
	return &redisServiceImpl{
		repository: repository,
	}
}
