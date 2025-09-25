package ws

import (
	"time"
)

type Chat struct {
	ID             int             `json:"id" db:"id"`
	Name           string          `json:"name" db:"name"`
	Description    *string         `json:"description,omitempty" db:"description"`
	CreatedBy      int             `json:"created_by" db:"created_by"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
	IsPrivate      bool            `json:"is_private" db:"is_private"`
	IsActive       bool            `json:"is_active" db:"is_active"`
	MaxMembers     int             `json:"max_members" db:"max_members"`
	CurrentMembers int             `json:"current_members" db:"current_members"`
	Clients        map[int]*Client `json:"-" db:"-"`
}

type Message struct {
	ID          int        `json:"id" db:"id"`
	ChatID      int        `json:"chat_id" db:"chat_id"`
	UserID      int        `json:"user_id" db:"user_id"`
	Username    string     `json:"username" db:"username"`
	Content     string     `json:"content" db:"content"`
	MessageType string     `json:"message_type" db:"message_type"`
	ReplyToID   *int       `json:"reply_to_id,omitempty" db:"reply_to_id"`
	EditedAt    *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	IsDeleted   bool       `json:"is_deleted" db:"is_deleted"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type ChatRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description,omitempty" binding:"omitempty,max=500"`
	IsPrivate   bool   `json:"is_private"`
	MaxMembers  int    `json:"max_members,omitempty" binding:"omitempty,min=2,max=1000"`
}

type ChatResponse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    *string   `json:"description,omitempty"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	IsPrivate      bool      `json:"is_private"`
	MaxMembers     int       `json:"max_members"`
	CurrentMembers int       `json:"current_members"`
}

type MessageRequest struct {
	Content     string `json:"content" binding:"required,min=1,max=4000"`
	MessageType string `json:"message_type,omitempty" binding:"omitempty,oneof=text image file system"`
	ReplyToID   *int   `json:"reply_to_id,omitempty"`
}

type JoinChatRequest struct {
	ChatID int `json:"chat_id" binding:"required"`
}

type ChatListResponse struct {
	Chats []ChatResponse `json:"chats"`
	Total int            `json:"total"`
}

type MessageListResponse struct {
	Messages []Message `json:"messages"`
	Total    int       `json:"total"`
	HasMore  bool      `json:"has_more"`
}

type UserChatRole struct {
	UserID      int        `json:"user_id" db:"user_id"`
	ChatID      int        `json:"chat_id" db:"chat_id"`
	Role        string     `json:"role" db:"role"`
	JoinedAt    time.Time  `json:"joined_at" db:"joined_at"`
	IsMuted     bool       `json:"is_muted" db:"is_muted"`
	IsBanned    bool       `json:"is_banned" db:"is_banned"`
	BannedUntil *time.Time `json:"banned_until,omitempty" db:"banned_until"`
	LastReadAt  *time.Time `json:"last_read_at,omitempty" db:"last_read_at"`
}

func (c *Chat) ToResponse() ChatResponse {
	return ChatResponse{
		ID:             c.ID,
		Name:           c.Name,
		Description:    c.Description,
		CreatedBy:      c.CreatedBy,
		CreatedAt:      c.CreatedAt,
		IsPrivate:      c.IsPrivate,
		MaxMembers:     c.MaxMembers,
		CurrentMembers: c.CurrentMembers,
	}
}
