package redis

import (
	"context"
	"fmt"
	"time"
)

const (
	ChatMembersKeyPrefix = "chat_members:"
	UserChatsKeyPrefix   = "user_chats:"

	ChatMembersTTL = 24 * time.Hour
	UserChatsTTL   = 24 * time.Hour
)

func (r *RedisClient) AddUser(chatID, userID int) error {
	ctx := context.Background()

	chatMembersKey := fmt.Sprintf("%s%d", ChatMembersKeyPrefix, chatID)
	if err := r.Client.SAdd(ctx, chatMembersKey, userID).Err(); err != nil {
		return fmt.Errorf("failed to add user to chat members: %w", err)
	}

	userChatsKey := fmt.Sprintf("%s%d", UserChatsKeyPrefix, userID)
	if err := r.Client.SAdd(ctx, userChatsKey, chatID).Err(); err != nil {
		r.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
		}).Warn("Failed to add chat to user chats")
	}

	r.Client.Expire(ctx, chatMembersKey, ChatMembersTTL)
	r.Client.Expire(ctx, userChatsKey, UserChatsTTL)

	if err := r.updateChatMemberCount(chatID); err != nil {
		r.logger.WithError(err).WithField("chat_id", chatID).Warn("Failed to update chat member count")
	}

	return nil
}

func (r *RedisClient) RemoveUser(chatID, userID int) error {
	ctx := context.Background()

	chatMembersKey := fmt.Sprintf("%s%d", ChatMembersKeyPrefix, chatID)
	if err := r.Client.SRem(ctx, chatMembersKey, userID).Err(); err != nil {
		return fmt.Errorf("failed to remove user from chat members: %w", err)
	}

	userChatsKey := fmt.Sprintf("%s%d", UserChatsKeyPrefix, userID)
	if err := r.Client.SRem(ctx, userChatsKey, chatID).Err(); err != nil {
		r.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id": userID,
			"chat_id": chatID,
		}).Warn("Failed to remove chat from user chats")
	}

	if err := r.updateChatMemberCount(chatID); err != nil {
		r.logger.WithError(err).WithField("chat_id", chatID).Warn("Failed to update chat member count")
	}

	return nil
}

func (r *RedisClient) GetChatMembers(chatID int) ([]int, error) {
	ctx := context.Background()

	chatMembersKey := fmt.Sprintf("%s%d", ChatMembersKeyPrefix, chatID)
	members, err := r.Client.SMembers(ctx, chatMembersKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get chat members: %w", err)
	}

	var memberIDs []int
	for _, memberStr := range members {
		var memberID int
		if _, err := fmt.Sscanf(memberStr, "%d", &memberID); err == nil {
			memberIDs = append(memberIDs, memberID)
		}
	}

	return memberIDs, nil
}

func (r *RedisClient) GetUserChats(userID int) ([]int, error) {
	ctx := context.Background()

	userChatsKey := fmt.Sprintf("%s%d", UserChatsKeyPrefix, userID)
	chats, err := r.Client.SMembers(ctx, userChatsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}

	var chatIDs []int
	for _, chatStr := range chats {
		var chatID int
		if _, err := fmt.Sscanf(chatStr, "%d", &chatID); err == nil {
			chatIDs = append(chatIDs, chatID)
		}
	}

	return chatIDs, nil
}

func (r *RedisClient) IsUserInChat(chatID, userID int) (bool, error) {
	ctx := context.Background()

	chatMembersKey := fmt.Sprintf("%s%d", ChatMembersKeyPrefix, chatID)
	exists, err := r.Client.SIsMember(ctx, chatMembersKey, userID).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if user is in chat: %w", err)
	}

	return exists, nil
}

func (r *RedisClient) GetChatMemberCount(chatID int) (int, error) {
	ctx := context.Background()

	chatMembersKey := fmt.Sprintf("%s%d", ChatMembersKeyPrefix, chatID)
	count, err := r.Client.SCard(ctx, chatMembersKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get chat member count: %w", err)
	}

	return int(count), nil
}

func (r *RedisClient) ClearChatMembers(chatID int) error {
	ctx := context.Background()

	chatMembersKey := fmt.Sprintf("%s%d", ChatMembersKeyPrefix, chatID)

	members, err := r.GetChatMembers(chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat members for clearing: %w", err)
	}

	for _, userID := range members {
		userChatsKey := fmt.Sprintf("%s%d", UserChatsKeyPrefix, userID)
		r.Client.SRem(ctx, userChatsKey, chatID)
	}

	if err := r.Client.Del(ctx, chatMembersKey).Err(); err != nil {
		return fmt.Errorf("failed to delete chat members: %w", err)
	}

	return nil
}

func (r *RedisClient) ClearUserChats(userID int) error {
	ctx := context.Background()

	userChatsKey := fmt.Sprintf("%s%d", UserChatsKeyPrefix, userID)

	chats, err := r.GetUserChats(userID)
	if err != nil {
		return fmt.Errorf("failed to get user chats for clearing: %w", err)
	}

	for _, chatID := range chats {
		chatMembersKey := fmt.Sprintf("%s%d", ChatMembersKeyPrefix, chatID)
		r.Client.SRem(ctx, chatMembersKey, userID)
	}

	if err := r.Client.Del(ctx, userChatsKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user chats: %w", err)
	}

	return nil
}

func (r *RedisClient) updateChatMemberCount(chatID int) error {
	count, err := r.GetChatMemberCount(chatID)
	if err != nil {
		return err
	}

	return r.UpdateChatMemberCount(chatID, count)
}
