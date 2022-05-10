package redis

import (
	"time"

	"github.com/DeniesKresna/jhapi2/repository/cache"
	"github.com/go-redis/redis"
)

type CacheRepositoryImpl struct {
	redisClient *redis.Client
}

func NewRepository(redis *redis.Client) cache.IRepository {
	return &CacheRepositoryImpl{redisClient: redis}
}

func (redis *CacheRepositoryImpl) Set(key string, value interface{}, duration time.Duration) error {
	_, err := redis.redisClient.Set(key, value, duration).Result()
	return err
}

func (redis *CacheRepositoryImpl) Get(key string) (value []byte, err error) {
	result, err := redis.redisClient.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

func (redis *CacheRepositoryImpl) Del(key string) error {
	result := redis.redisClient.Del(key)
	return result.Err()
}

func (redis *CacheRepositoryImpl) Expire(key string, duration time.Duration) error {
	return redis.redisClient.Expire(key, duration).Err()
}
