package ws

import (
	"encoding/json"
	"fmt"
	"onlineChat/pkg/config"
	"sync"
	"time"

	"onlineChat/pkg/redis"

	"github.com/sirupsen/logrus"
)

type Hub struct {
	chats      map[int]map[int]*Client
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	redis      *redis.RedisClient
	service    ChatService
	logger     *logrus.Logger
	mu         sync.RWMutex
}

func NewHub(redisCfg config.RedisConfig, service ChatService, logger *logrus.Logger) *Hub {
	return &Hub{
		chats:      make(map[int]map[int]*Client),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis: redis.NewRedisClient(redis.RedisConfig{
			Address:  redisCfg.Address,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		}, logger),
		service: service,
		logger:  logger,
	}
}

func (h *Hub) Run() {
	h.logger.Info("Starting WebSocket hub")

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.chats[client.ChatID] == nil {
		h.chats[client.ChatID] = make(map[int]*Client)
	}

	if existingClient, exists := h.chats[client.ChatID][client.ID]; exists {
		h.logger.WithFields(logrus.Fields{
			"user_id":  client.ID,
			"username": client.Username,
			"chat_id":  client.ChatID,
		}).Warn("User already connected to chat, closing existing connection")

		existingClient.Close()
	}

	h.chats[client.ChatID][client.ID] = client
	if err := h.redis.AddUser(client.ChatID, client.ID); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":  client.ID,
			"username": client.Username,
			"chat_id":  client.ChatID,
		}).Error("Failed to add user to Redis")
	}

	h.sendSystemMessage(client.ChatID, fmt.Sprintf("User %s joined the chat", client.Username))

	h.logger.WithFields(logrus.Fields{
		"user_id":  client.ID,
		"username": client.Username,
		"chat_id":  client.ChatID,
		"clients":  len(h.chats[client.ChatID]),
	}).Info("Client registered")
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if chat, exists := h.chats[client.ChatID]; exists {
		if _, clientExists := chat[client.ID]; clientExists {
			if err := h.redis.RemoveUser(client.ChatID, client.ID); err != nil {
				h.logger.WithError(err).WithFields(logrus.Fields{
					"user_id":  client.ID,
					"username": client.Username,
					"chat_id":  client.ChatID,
				}).Error("Failed to remove user from Redis")
			}

			client.Close()

			delete(chat, client.ID)

			if len(chat) == 0 {
				delete(h.chats, client.ChatID)
			}

			h.sendSystemMessage(client.ChatID, fmt.Sprintf("User %s left the chat", client.Username))

			h.logger.WithFields(logrus.Fields{
				"user_id":  client.ID,
				"username": client.Username,
				"chat_id":  client.ChatID,
				"clients":  len(h.chats[client.ChatID]),
			}).Info("Client unregistered")
		}
	}
}

func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	messageCache := &redis.MessageCache{
		ID:          message.ID,
		ChatID:      message.ChatID,
		UserID:      message.UserID,
		Username:    message.Username,
		Content:     message.Content,
		MessageType: message.MessageType,
		CreatedAt:   message.CreatedAt,
	}

	if err := h.redis.CacheMessage(messageCache); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"chat_id":    message.ChatID,
			"user_id":    message.UserID,
			"message_id": message.ID,
		}).Error("Failed to cache message")
	}

	if err := h.redis.UpdateChatLastMessage(message.ChatID, message.CreatedAt); err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"chat_id": message.ChatID,
		}).Warn("Failed to update chat last message timestamp")
	}

	if chat, exists := h.chats[message.ChatID]; exists {
		messageData, err := json.Marshal(message)
		if err != nil {
			h.logger.WithError(err).WithFields(logrus.Fields{
				"chat_id": message.ChatID,
				"user_id": message.UserID,
			}).Error("Failed to marshal message")
			return
		}

		for clientID, client := range chat {
			select {
			case client.Send <- messageData:
			default:
				h.logger.WithFields(logrus.Fields{
					"client_id": clientID,
					"chat_id":   message.ChatID,
				}).Warn("Failed to send message to client, closing connection")

				client.Close()
				delete(chat, clientID)
			}
		}
	}
}

func (h *Hub) sendSystemMessage(chatID int, content string) {
	systemMessage := &Message{
		ChatID:      chatID,
		UserID:      0,
		Username:    "System",
		Content:     content,
		MessageType: "system",
		CreatedAt:   time.Now(),
	}

	h.broadcast <- systemMessage
}

func (h *Hub) GetChatClients(chatID int) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []*Client
	if chat, exists := h.chats[chatID]; exists {
		for _, client := range chat {
			clients = append(clients, client)
		}
	}

	return clients
}

func (h *Hub) GetChatClientCount(chatID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if chat, exists := h.chats[chatID]; exists {
		return len(chat)
	}

	return 0
}

func (h *Hub) IsUserInChat(chatID, userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if chat, exists := h.chats[chatID]; exists {
		_, userExists := chat[userID]
		return userExists
	}

	return false
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, chat := range h.chats {
		for _, client := range chat {
			client.Close()
		}
	}

	h.chats = make(map[int]map[int]*Client)
	h.logger.Info("Hub closed")
}
