package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ChatMessagesPrefix = "chat_messages:"
	MessageKeyPrefix   = "message:"
	ChatRecentPrefix   = "chat_recent:"

	MessageTTL        = 7 * 24 * time.Hour // 7 days
	RecentMessagesTTL = 24 * time.Hour     // 1 day
	ChatMetadataTTL   = 1 * time.Hour      // 1 hour

	MaxRecentMessages = 100
	MaxCachedMessages = 1000
)

func (r *RedisClient) CacheMessage(message *MessageCache) error {
	ctx := context.Background()

	messageKey := fmt.Sprintf("%s%d", MessageKeyPrefix, message.ID)
	messageData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := r.Client.Set(ctx, messageKey, messageData, MessageTTL).Err(); err != nil {
		return fmt.Errorf("failed to cache message: %w", err)
	}

	chatMessagesKey := fmt.Sprintf("%s%d", ChatMessagesPrefix, message.ChatID)
	timestamp := float64(message.CreatedAt.Unix())

	if err := r.Client.ZAdd(ctx, chatMessagesKey, &redis.Z{
		Score:  timestamp,
		Member: message.ID,
	}).Err(); err != nil {
		return fmt.Errorf("failed to add message to chat list: %w", err)
	}

	if err := r.Client.ZRemRangeByRank(ctx, chatMessagesKey, 0, -MaxCachedMessages-1).Err(); err != nil {
		r.logger.WithError(err).WithField("chat_id", message.ChatID).Warn("Failed to trim chat messages")
	}

	r.Client.Expire(ctx, chatMessagesKey, MessageTTL)

	if err := r.UpdateRecentMessages(message.ChatID, message); err != nil {
		r.logger.WithError(err).WithField("chat_id", message.ChatID).Warn("Failed to update recent messages")
	}

	return nil
}

func (r *RedisClient) GetCachedMessage(messageID int) (*MessageCache, error) {
	ctx := context.Background()

	messageKey := fmt.Sprintf("%s%d", MessageKeyPrefix, messageID)
	messageData, err := r.Client.Get(ctx, messageKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cached message: %w", err)
	}

	var message MessageCache
	if err := json.Unmarshal([]byte(messageData), &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached message: %w", err)
	}

	return &message, nil
}

func (r *RedisClient) GetChatMessages(chatID int, limit, offset int) ([]*MessageCache, error) {
	ctx := context.Background()

	chatMessagesKey := fmt.Sprintf("%s%d", ChatMessagesPrefix, chatID)

	messageIDs, err := r.Client.ZRevRange(ctx, chatMessagesKey, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get message IDs: %w", err)
	}

	var messages []*MessageCache
	for _, messageIDStr := range messageIDs {
		messageID, err := strconv.Atoi(messageIDStr)
		if err != nil {
			r.logger.WithError(err).WithField("message_id", messageIDStr).Warn("Invalid message ID in cache")
			continue
		}

		message, err := r.GetCachedMessage(messageID)
		if err != nil {
			r.logger.WithError(err).WithField("message_id", messageID).Warn("Failed to get cached message")
			continue
		}

		if message != nil {
			messages = append(messages, message)
		}
	}

	return messages, nil
}

func (r *RedisClient) UpdateRecentMessages(chatID int, message *MessageCache) error {
	ctx := context.Background()

	recentKey := fmt.Sprintf("%s%d", ChatRecentPrefix, chatID)

	messageData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal recent message: %w", err)
	}

	if err := r.Client.LPush(ctx, recentKey, messageData).Err(); err != nil {
		return fmt.Errorf("failed to add to recent messages: %w", err)
	}

	if err := r.Client.LTrim(ctx, recentKey, 0, MaxRecentMessages-1).Err(); err != nil {
		return fmt.Errorf("failed to trim recent messages: %w", err)
	}

	r.Client.Expire(ctx, recentKey, RecentMessagesTTL)

	return nil
}

func (r *RedisClient) GetRecentMessages(chatID int, limit int) ([]*MessageCache, error) {
	ctx := context.Background()

	recentKey := fmt.Sprintf("%s%d", ChatRecentPrefix, chatID)

	messageDataList, err := r.Client.LRange(ctx, recentKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent messages: %w", err)
	}

	var messages []*MessageCache
	for _, messageData := range messageDataList {
		var message MessageCache
		if err := json.Unmarshal([]byte(messageData), &message); err != nil {
			r.logger.WithError(err).Warn("Failed to unmarshal recent message")
			continue
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (r *RedisClient) DeleteMessageFromCache(messageID int) error {
	ctx := context.Background()

	messageKey := fmt.Sprintf("%s%d", MessageKeyPrefix, messageID)

	message, err := r.GetCachedMessage(messageID)
	if err != nil {
		return fmt.Errorf("failed to get message for deletion: %w", err)
	}

	if message == nil {
		return nil
	}

	if err := r.Client.Del(ctx, messageKey).Err(); err != nil {
		return fmt.Errorf("failed to delete message from cache: %w", err)
	}

	chatMessagesKey := fmt.Sprintf("%s%d", ChatMessagesPrefix, message.ChatID)
	if err := r.Client.ZRem(ctx, chatMessagesKey, messageID).Err(); err != nil {
		r.logger.WithError(err).WithField("chat_id", message.ChatID).Warn("Failed to remove message from chat list")
	}

	return nil
}

func (r *RedisClient) ClearChatMessages(chatID int) error {
	ctx := context.Background()

	chatMessagesKey := fmt.Sprintf("%s%d", ChatMessagesPrefix, chatID)
	recentKey := fmt.Sprintf("%s%d", ChatRecentPrefix, chatID)

	messageIDs, err := r.Client.ZRange(ctx, chatMessagesKey, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get chat message IDs: %w", err)
	}

	for _, messageIDStr := range messageIDs {
		messageKey := fmt.Sprintf("%s%s", MessageKeyPrefix, messageIDStr)
		r.Client.Del(ctx, messageKey)
	}

	if err := r.Client.Del(ctx, chatMessagesKey).Err(); err != nil {
		return fmt.Errorf("failed to delete chat messages list: %w", err)
	}

	if err := r.Client.Del(ctx, recentKey).Err(); err != nil {
		return fmt.Errorf("failed to delete recent messages: %w", err)
	}

	return nil
}

func (r *RedisClient) CacheChatMetadata(chat *ChatMetadata) error {
	ctx := context.Background()

	chatKey := fmt.Sprintf("chat_meta:%d", chat.ID)
	chatData, err := json.Marshal(chat)
	if err != nil {
		return fmt.Errorf("failed to marshal chat metadata: %w", err)
	}

	if err := r.Client.Set(ctx, chatKey, chatData, ChatMetadataTTL).Err(); err != nil {
		return fmt.Errorf("failed to cache chat metadata: %w", err)
	}

	return nil
}

func (r *RedisClient) GetCachedChatMetadata(chatID int) (*ChatMetadata, error) {
	ctx := context.Background()

	chatKey := fmt.Sprintf("chat_meta:%d", chatID)
	chatData, err := r.Client.Get(ctx, chatKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil // Chat metadata not found in cache
		}
		return nil, fmt.Errorf("failed to get cached chat metadata: %w", err)
	}

	var chat ChatMetadata
	if err := json.Unmarshal([]byte(chatData), &chat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached chat metadata: %w", err)
	}

	return &chat, nil
}

func (r *RedisClient) UpdateChatMemberCount(chatID int, count int) error {
	chat, err := r.GetCachedChatMetadata(chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat metadata for update: %w", err)
	}

	if chat == nil {
		chat = &ChatMetadata{
			ID:          chatID,
			MemberCount: count,
			IsActive:    true,
		}
	} else {
		chat.MemberCount = count
	}

	return r.CacheChatMetadata(chat)
}

func (r *RedisClient) UpdateChatLastMessage(chatID int, lastMessage time.Time) error {
	chat, err := r.GetCachedChatMetadata(chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat metadata for update: %w", err)
	}

	if chat == nil {
		chat = &ChatMetadata{
			ID:          chatID,
			LastMessage: lastMessage,
			IsActive:    true,
		}
	} else {
		chat.LastMessage = lastMessage
	}

	return r.CacheChatMetadata(chat)
}
