package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type RedisClient struct {
	Client *redis.Client
	logger *logrus.Logger
}

type UserSession struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	ChatID    int       `json:"chat_id"`
	Connected bool      `json:"connected"`
	LastSeen  time.Time `json:"last_seen"`
}

type ChatMetadata struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	MemberCount int       `json:"member_count"`
	LastMessage time.Time `json:"last_message"`
	IsActive    bool      `json:"is_active"`
}

type MessageCache struct {
	ID          int       `json:"id"`
	ChatID      int       `json:"chat_id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewRedisClient(cfg RedisConfig, logger *logrus.Logger) *RedisClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to Redis")
	}

	logger.WithField("response", pong).Info("Redis connection established")

	return &RedisClient{
		Client: client,
		logger: logger,
	}
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

func (r *RedisClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := r.Client.Ping(ctx).Result()
	return err
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.Client
}
