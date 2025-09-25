package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"time"
)

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg RedisConfig) *RedisClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		logrus.Fatal(pong, err)
	}

	return &RedisClient{
		Client: client,
	}
}
