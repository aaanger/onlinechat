package ws

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ChatService interface {
	CreateChat(req ChatRequest, userID int) (*ChatResponse, error)
	GetChatByID(chatID int) (*ChatResponse, error)
	GetUserChats(userID int, limit, offset int) (*ChatListResponse, error)
	SearchPublicChats(userID int, searchTerm string, limit, offset int) (*ChatListResponse, error)
	JoinChat(userID, chatID int) error
	LeaveChat(userID, chatID int) error
	SaveMessage(message *Message) error
	GetMessages(chatID int, limit, offset int) (*MessageListResponse, error)
	UpdateChat(chatID int, userID int, req ChatRequest) (*ChatResponse, error)
	DeleteChat(chatID int, userID int) error
	GetChatMembers(chatID int) ([]int, error)
}

type chatService struct {
	repo   ChatRepository
	logger *logrus.Logger
}

func NewChatService(repo ChatRepository, logger *logrus.Logger) ChatService {
	return &chatService{
		repo:   repo,
		logger: logger,
	}
}

func (s *chatService) CreateChat(req ChatRequest, userID int) (*ChatResponse, error) {
	if req.MaxMembers == 0 {
		req.MaxMembers = 100
	}

	chat := &Chat{
		Name:        req.Name,
		Description: &req.Description,
		CreatedBy:   userID,
		IsPrivate:   req.IsPrivate,
		MaxMembers:  req.MaxMembers,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
	}

	createdChat, err := s.repo.CreateChat(chat)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"name":    req.Name,
		}).Error("Failed to create chat")
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	if err := s.repo.AddUserToChat(userID, createdChat.ID, "owner"); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": createdChat.ID,
		}).Error("Failed to add creator to chat")
		return nil, fmt.Errorf("failed to add creator to chat: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": createdChat.ID,
		"name":    req.Name,
	}).Info("Chat created successfully")

	response := createdChat.ToResponse()
	return &response, nil
}

func (s *chatService) GetChatByID(chatID int) (*ChatResponse, error) {
	chat, err := s.repo.GetChatByID(chatID)
	if err != nil {
		s.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to get chat")
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	response := chat.ToResponse()
	return &response, nil
}

func (s *chatService) GetUserChats(userID int, limit, offset int) (*ChatListResponse, error) {
	chats, total, err := s.repo.GetUserChats(userID, limit, offset)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user chats")
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}

	var chatResponses []ChatResponse
	for _, chat := range chats {
		chatResponses = append(chatResponses, chat.ToResponse())
	}

	return &ChatListResponse{
		Chats: chatResponses,
		Total: total,
	}, nil
}

func (s *chatService) SearchPublicChats(userID int, searchTerm string, limit, offset int) (*ChatListResponse, error) {
	chats, total, err := s.repo.SearchPublicChats(userID, searchTerm, limit, offset)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to search public chats")
		return nil, fmt.Errorf("failed to search public chats: %w", err)
	}

	var chatResponses []ChatResponse
	for _, chat := range chats {
		chatResponses = append(chatResponses, chat.ToResponse())
	}

	return &ChatListResponse{
		Chats: chatResponses,
		Total: total,
	}, nil
}

func (s *chatService) JoinChat(userID, chatID int) error {
	chat, err := s.repo.GetChatByID(chatID)
	if err != nil {
		return fmt.Errorf("chat not found: %w", err)
	}

	if !chat.IsActive {
		return fmt.Errorf("chat is not active")
	}

	members, err := s.repo.GetChatMembers(chatID)
	if err != nil {
		return fmt.Errorf("failed to get chat members: %w", err)
	}

	for _, memberID := range members {
		if memberID == userID {
			return fmt.Errorf("user already in chat")
		}
	}

	if len(members) >= chat.MaxMembers {
		return fmt.Errorf("chat is full")
	}

	if err := s.repo.AddUserToChat(userID, chatID, "member"); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to add user to chat")
		return fmt.Errorf("failed to join chat: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("User joined chat")

	return nil
}

func (s *chatService) LeaveChat(userID, chatID int) error {
	if err := s.repo.RemoveUserFromChat(userID, chatID); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to remove user from chat")
		return fmt.Errorf("failed to leave chat: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("User left chat")

	return nil
}

func (s *chatService) SaveMessage(message *Message) error {
	if err := s.repo.SaveMessage(message); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": message.UserID,
			"chat_id": message.ChatID,
		}).Error("Failed to save message")
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

func (s *chatService) GetMessages(chatID int, limit, offset int) (*MessageListResponse, error) {
	messages, total, err := s.repo.GetMessages(chatID, limit, offset)
	if err != nil {
		s.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to get messages")
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	hasMore := (offset + len(messages)) < total

	return &MessageListResponse{
		Messages: messages,
		Total:    total,
		HasMore:  hasMore,
	}, nil
}

func (s *chatService) UpdateChat(chatID int, userID int, req ChatRequest) (*ChatResponse, error) {
	role, err := s.repo.GetUserRoleInChat(userID, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	if role != "owner" && role != "admin" {
		return nil, fmt.Errorf("insufficient permissions")
	}

	chat, err := s.repo.UpdateChat(chatID, req)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to update chat")
		return nil, fmt.Errorf("failed to update chat: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Chat updated successfully")

	response := chat.ToResponse()
	return &response, nil
}

func (s *chatService) DeleteChat(chatID int, userID int) error {
	role, err := s.repo.GetUserRoleInChat(userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}

	if role != "owner" {
		return fmt.Errorf("only chat owner can delete the chat")
	}

	if err := s.repo.DeleteChat(chatID); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to delete chat")
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"chat_id": chatID,
	}).Info("Chat deleted successfully")

	return nil
}

func (s *chatService) GetChatMembers(chatID int) ([]int, error) {
	members, err := s.repo.GetChatMembers(chatID)
	if err != nil {
		s.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to get chat members")
		return nil, fmt.Errorf("failed to get chat members: %w", err)
	}

	return members, nil
}
