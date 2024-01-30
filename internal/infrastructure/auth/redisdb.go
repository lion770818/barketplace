package auth

import (
	//"github.com/go-redis/redis/v7"
	redis "github.com/redis/go-redis/v9"
)

type RedisService struct {
	Auth   AuthInterface
	Client *redis.Client
}

func NewRedisDB(client *redis.Client) (*RedisService, error) {

	return &RedisService{
		Auth:   NewAuth(client),
		Client: client,
	}, nil
}
