package redis

import (
	"context"
	"fmt"
)

func (r *RedisClient) AddUser(chatID, userID int) error {
	ctx := context.Background()
	key := fmt.Sprintf("chat:%d:users", chatID)
	return r.Client.SAdd(ctx, key, userID).Err()
}

func (r *RedisClient) RemoveUser(chatID, userID int) error {
	ctx := context.Background()
	key := fmt.Sprintf("chat:%d:users", chatID)
	return r.Client.SRem(ctx, key, userID).Err()
}
