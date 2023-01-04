package database

import (
	errors "errors"
	"github.com/go-redis/redis"
)

func NewInitRedis() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:63879",
		Password: "",
		DB:       0,
	})
	if redisClient == nil {
		return nil, errors.New("can not connect redis")
	}

	return redisClient, nil
}
