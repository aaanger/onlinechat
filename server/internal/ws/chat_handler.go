package ws

import (
	"net/http"
	"strconv"
	"time"

	"onlineChat/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	hub     *Hub
	service ChatService
	logger  *logrus.Logger
}

func NewChatHandler(hub *Hub, service ChatService, logger *logrus.Logger) *Handler {
	return &Handler{
		hub:     hub,
		service: service,
		logger:  logger,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should implement proper origin checking
		return true
	},
}

func (h *Handler) ServeWS(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, err := utils.GetUsername(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get username from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	chatIDStr := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		h.logger.WithError(err).WithField("chat_id", chatIDStr).Error("Invalid chat ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	members, err := h.service.GetChatMembers(chatID)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to get chat members")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify chat membership"})
		return
	}

	isMember := false
	for _, memberID := range members {
		if memberID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		h.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Warn("User attempted to connect to chat they're not a member of")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to upgrade connection to WebSocket")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish WebSocket connection"})
		return
	}

	client := &Client{
		ID:         userID,
		Username:   username,
		ChatID:     chatID,
		Connection: conn,
		Send:       make(chan []byte, 256),
		Hub:        h.hub,
		LastPing:   time.Now(),
	}

	h.hub.register <- client

	go client.writePump()
	go client.readPump()

	h.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"username": username,
		"chat_id":  chatID,
	}).Info("WebSocket connection established")
}

func (h *Handler) CreateChat(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid chat creation request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	chat, err := h.service.CreateChat(req, userID)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"name":    req.Name,
		}).Error("Failed to create chat")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat": chat})
}

func (h *Handler) GetAllChats(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	chats, err := h.service.GetUserChats(userID, limit, offset)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user chats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *Handler) JoinChat(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	chatIDStr := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		h.logger.WithError(err).WithField("chat_id", chatIDStr).Error("Invalid chat ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	err = h.service.JoinChat(userID, chatID)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to join chat")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined chat"})
}

func (h *Handler) LeaveChat(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	chatIDStr := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		h.logger.WithError(err).WithField("chat_id", chatIDStr).Error("Invalid chat ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	err = h.service.LeaveChat(userID, chatID)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to leave chat")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left chat"})
}

func (h *Handler) SearchPublicChats(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	searchTerm := c.Query("search")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	chats, err := h.service.SearchPublicChats(userID, searchTerm, limit, offset)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Error("Failed to search public chats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *Handler) GetChatMessages(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	chatIDStr := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		h.logger.WithError(err).WithField("chat_id", chatIDStr).Error("Invalid chat ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	members, err := h.service.GetChatMembers(chatID)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to get chat members")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify chat membership"})
		return
	}

	isMember := false
	for _, memberID := range members {
		if memberID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	messages, err := h.service.GetMessages(chatID, limit, offset)
	if err != nil {
		h.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to get messages")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) GetClientsByChatID(c *gin.Context) {
	userID, err := utils.GetUserID(c)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	chatIDStr := c.Param("chatID")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		h.logger.WithError(err).WithField("chat_id", chatIDStr).Error("Invalid chat ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	members, err := h.service.GetChatMembers(chatID)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to get chat members")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify chat membership"})
		return
	}

	isMember := false
	for _, memberID := range members {
		if memberID == userID {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	clients := h.hub.GetChatClients(chatID)

	var clientResponses []gin.H
	for _, client := range clients {
		clientResponses = append(clientResponses, gin.H{
			"id":       client.ID,
			"username": client.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"clients": clientResponses,
		"count":   len(clientResponses),
	})
}
