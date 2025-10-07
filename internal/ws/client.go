package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
	maxMessageSize  = 4096
)

type Client struct {
	ID         int             `json:"id"`
	Username   string          `json:"username"`
	ChatID     int             `json:"chat_id"`
	Connection *websocket.Conn `json:"-"`
	Message    chan *Message   `json:"-"`
	Send       chan []byte     `json:"-"`
	Hub        *Hub            `json:"-"`
	LastPing   time.Time       `json:"-"`
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Connection.Close()
	}()

	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageData, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Hub.logger.WithError(err).WithFields(logrus.Fields{
					"client_id": c.ID,
					"username":  c.Username,
					"chat_id":   c.ChatID,
				}).Error("WebSocket connection closed unexpectedly")
			}
			break
		}

		var msg MessageRequest
		if err := json.Unmarshal(messageData, &msg); err != nil {
			c.Hub.logger.WithError(err).WithFields(logrus.Fields{
				"client_id": c.ID,
				"username":  c.Username,
				"chat_id":   c.ChatID,
			}).Error("Failed to parse message")
			continue
		}

		if err := c.validateMessage(msg); err != nil {
			c.sendError(fmt.Sprintf("Invalid message: %v", err))
			continue
		}

		message := &Message{
			ChatID:      c.ChatID,
			UserID:      c.ID,
			Username:    c.Username,
			Content:     msg.Content,
			MessageType: msg.MessageType,
			ReplyToID:   msg.ReplyToID,
			CreatedAt:   time.Now(),
		}

		if err := c.Hub.service.SaveMessage(message); err != nil {
			c.Hub.logger.WithError(err).WithFields(logrus.Fields{
				"client_id": c.ID,
				"username":  c.Username,
				"chat_id":   c.ChatID,
			}).Error("Failed to save message")
			c.sendError("Failed to save message")
			continue
		}

		c.Hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
				c.Hub.logger.WithError(err).WithFields(logrus.Fields{
					"client_id": c.ID,
					"username":  c.Username,
					"chat_id":   c.ChatID,
				}).Error("Failed to write message to WebSocket")
				return
			}

		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Hub.logger.WithError(err).WithFields(logrus.Fields{
					"client_id": c.ID,
					"username":  c.Username,
					"chat_id":   c.ChatID,
				}).Error("Failed to send ping")
				return
			}
		}
	}
}

func (c *Client) validateMessage(msg MessageRequest) error {
	if len(msg.Content) == 0 {
		return fmt.Errorf("message content cannot be empty")
	}

	if len(msg.Content) > 4000 {
		return fmt.Errorf("message content too long")
	}

	if msg.MessageType == "" {
		msg.MessageType = "text"
	}

	validTypes := map[string]bool{
		"text":   true,
		"image":  true,
		"file":   true,
		"system": true,
	}

	if !validTypes[msg.MessageType] {
		return fmt.Errorf("invalid message type")
	}

	return nil
}

func (c *Client) sendError(message string) {
	errorMsg := map[string]interface{}{
		"type":    "error",
		"message": message,
		"time":    time.Now(),
	}

	if data, err := json.Marshal(errorMsg); err == nil {
		select {
		case c.Send <- data:
		default:
			// If we can't send the error message, close the connection
			c.Connection.Close()
		}
	}
}

func (c *Client) sendMessage(message *Message) {
	if data, err := json.Marshal(message); err == nil {
		select {
		case c.Send <- data:
		default:
			// If we can't send the message, close the connection
			c.Connection.Close()
		}
	}
}

func (c *Client) IsActive() bool {
	return c.Connection != nil
}

func (c *Client) Close() {
	if c.Connection != nil {
		c.Connection.Close()
	}
	close(c.Send)
}
