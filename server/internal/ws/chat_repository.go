package ws

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ChatRepository interface {
	CreateChat(chat *Chat) (*Chat, error)
	GetChatByID(chatID int) (*Chat, error)
	GetUserChats(userID int, limit, offset int) ([]Chat, int, error)
	SearchPublicChats(userID int, searchTerm string, limit, offset int) ([]Chat, int, error)
	UpdateChat(chatID int, req ChatRequest) (*Chat, error)
	DeleteChat(chatID int) error
	AddUserToChat(userID, chatID int, role string) error
	RemoveUserFromChat(userID, chatID int) error
	GetChatMembers(chatID int) ([]int, error)
	GetUserRoleInChat(userID, chatID int) (string, error)
	SaveMessage(message *Message) error
	GetMessages(chatID int, limit, offset int) ([]Message, int, error)
}

type chatRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewChatRepository(db *sql.DB, logger *logrus.Logger) ChatRepository {
	return &chatRepository{
		db:     db,
		logger: logger,
	}
}

func (r *chatRepository) CreateChat(chat *Chat) (*Chat, error) {
	query := `
		INSERT INTO chats (name, description, created_by, created_at, updated_at, 
		                  is_private, is_active, max_members, current_members)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	row := r.db.QueryRow(query,
		chat.Name, chat.Description, chat.CreatedBy, now, now,
		chat.IsPrivate, chat.IsActive, chat.MaxMembers, chat.CurrentMembers,
	)

	err := row.Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		r.logger.WithError(err).Error("Failed to create chat")
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return chat, nil
}

func (r *chatRepository) GetChatByID(chatID int) (*Chat, error) {
	query := `
		SELECT id, name, description, created_by, created_at, updated_at,
		       is_private, is_active, max_members, current_members
		FROM chats
		WHERE id = $1 AND is_active = true
	`

	chat := &Chat{}
	row := r.db.QueryRow(query, chatID)

	err := row.Scan(
		&chat.ID, &chat.Name, &chat.Description, &chat.CreatedBy,
		&chat.CreatedAt, &chat.UpdatedAt, &chat.IsPrivate, &chat.IsActive,
		&chat.MaxMembers, &chat.CurrentMembers,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chat not found")
		}
		r.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to get chat")
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	return chat, nil
}

func (r *chatRepository) GetUserChats(userID int, limit, offset int) ([]Chat, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM chats c
		INNER JOIN user_chat uc ON c.id = uc.chat_id
		WHERE uc.user_id = $1 AND c.is_active = true
	`

	var total int
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		r.logger.WithError(err).WithField("user_id", userID).Error("Failed to count user chats")
		return nil, 0, fmt.Errorf("failed to count user chats: %w", err)
	}

	query := `
		SELECT c.id, c.name, c.description, c.created_by, c.created_at, c.updated_at,
		       c.is_private, c.is_active, c.max_members, c.current_members
		FROM chats c
		INNER JOIN user_chat uc ON c.id = uc.chat_id
		WHERE uc.user_id = $1 AND c.is_active = true
		ORDER BY c.updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		r.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user chats")
		return nil, 0, fmt.Errorf("failed to get user chats: %w", err)
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		err := rows.Scan(
			&chat.ID, &chat.Name, &chat.Description, &chat.CreatedBy,
			&chat.CreatedAt, &chat.UpdatedAt, &chat.IsPrivate, &chat.IsActive,
			&chat.MaxMembers, &chat.CurrentMembers,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan chat")
			continue
		}
		chats = append(chats, chat)
	}

	return chats, total, nil
}

func (r *chatRepository) SearchPublicChats(userID int, searchTerm string, limit, offset int) ([]Chat, int, error) {
	searchCondition := ""
	args := []interface{}{userID}
	argIndex := 2

	if searchTerm != "" {
		searchCondition = " AND (c.name ILIKE $2 OR c.description ILIKE $2)"
		searchPattern := "%" + searchTerm + "%"
		args = append(args, searchPattern)
		argIndex = 3
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM chats c
		WHERE c.is_active = true 
		AND c.is_private = false
		AND c.id NOT IN (
			SELECT chat_id 
			FROM user_chat 
			WHERE user_id = $1
		)
		%s
	`, searchCondition)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.WithError(err).Error("Failed to count public chats")
		return nil, 0, fmt.Errorf("failed to count public chats: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT c.id, c.name, c.description, c.created_by, c.created_at, c.updated_at,
		       c.is_private, c.is_active, c.max_members, c.current_members
		FROM chats c
		WHERE c.is_active = true 
		AND c.is_private = false
		AND c.id NOT IN (
			SELECT chat_id 
			FROM user_chat 
			WHERE user_id = $1
		)
		%s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, searchCondition, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to search public chats")
		return nil, 0, fmt.Errorf("failed to search public chats: %w", err)
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		err := rows.Scan(
			&chat.ID, &chat.Name, &chat.Description, &chat.CreatedBy,
			&chat.CreatedAt, &chat.UpdatedAt, &chat.IsPrivate, &chat.IsActive,
			&chat.MaxMembers, &chat.CurrentMembers,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan chat")
			continue
		}
		chats = append(chats, chat)
	}

	return chats, total, nil
}

func (r *chatRepository) UpdateChat(chatID int, req ChatRequest) (*Chat, error) {
	query := `
		UPDATE chats 
		SET name = $1, description = $2, max_members = $3, updated_at = $4
		WHERE id = $5 AND is_active = true
		RETURNING id, name, description, created_by, created_at, updated_at,
		          is_private, is_active, max_members, current_members
	`

	chat := &Chat{}
	row := r.db.QueryRow(query, req.Name, req.Description, req.MaxMembers, time.Now(), chatID)

	err := row.Scan(
		&chat.ID, &chat.Name, &chat.Description, &chat.CreatedBy,
		&chat.CreatedAt, &chat.UpdatedAt, &chat.IsPrivate, &chat.IsActive,
		&chat.MaxMembers, &chat.CurrentMembers,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chat not found")
		}
		r.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to update chat")
		return nil, fmt.Errorf("failed to update chat: %w", err)
	}

	return chat, nil
}

func (r *chatRepository) DeleteChat(chatID int) error {
	query := `UPDATE chats SET is_active = false, updated_at = $1 WHERE id = $2`

	result, err := r.db.Exec(query, time.Now(), chatID)
	if err != nil {
		r.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to delete chat")
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat not found")
	}

	return nil
}

func (r *chatRepository) AddUserToChat(userID, chatID int, role string) error {
	query := `
		INSERT INTO user_chat (user_id, chat_id, role, joined_at, last_read_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, chat_id) DO UPDATE SET
		role = EXCLUDED.role,
		is_banned = false,
		banned_until = NULL
	`

	now := time.Now()
	_, err := r.db.Exec(query, userID, chatID, role, now, now)
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to add user to chat")
		return fmt.Errorf("failed to add user to chat: %w", err)
	}

	return nil
}

func (r *chatRepository) RemoveUserFromChat(userID, chatID int) error {
	query := `DELETE FROM user_chat WHERE user_id = $1 AND chat_id = $2`

	result, err := r.db.Exec(query, userID, chatID)
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to remove user from chat")
		return fmt.Errorf("failed to remove user from chat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found in chat")
	}

	return nil
}

func (r *chatRepository) GetChatMembers(chatID int) ([]int, error) {
	query := `
		SELECT user_id
		FROM user_chat
		WHERE chat_id = $1 AND is_banned = false
	`

	rows, err := r.db.Query(query, chatID)
	if err != nil {
		r.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to get chat members")
		return nil, fmt.Errorf("failed to get chat members: %w", err)
	}
	defer rows.Close()

	var members []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			r.logger.WithError(err).Error("Failed to scan user ID")
			continue
		}
		members = append(members, userID)
	}

	return members, nil
}

func (r *chatRepository) GetUserRoleInChat(userID, chatID int) (string, error) {
	query := `
		SELECT role
		FROM user_chat
		WHERE user_id = $1 AND chat_id = $2 AND is_banned = false
	`

	var role string
	err := r.db.QueryRow(query, userID, chatID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found in chat")
		}
		r.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": userID,
			"chat_id": chatID,
		}).Error("Failed to get user role in chat")
		return "", fmt.Errorf("failed to get user role: %w", err)
	}

	return role, nil
}

func (r *chatRepository) SaveMessage(message *Message) error {
	query := `
		INSERT INTO messages (chat_id, user_id, content, message_type, reply_to_id, 
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	row := r.db.QueryRow(query,
		message.ChatID, message.UserID, message.Content, message.MessageType,
		message.ReplyToID, now, now,
	)

	err := row.Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"user_id": message.UserID,
			"chat_id": message.ChatID,
		}).Error("Failed to save message")
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

func (r *chatRepository) GetMessages(chatID int, limit, offset int) ([]Message, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM messages
		WHERE chat_id = $1 AND is_deleted = false
	`

	var total int
	err := r.db.QueryRow(countQuery, chatID).Scan(&total)
	if err != nil {
		r.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to count messages")
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	query := `
		SELECT m.id, m.chat_id, m.user_id, u.username, m.content, m.message_type,
		       m.reply_to_id, m.edited_at, m.is_deleted, m.deleted_at,
		       m.created_at, m.updated_at
		FROM messages m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.chat_id = $1 AND m.is_deleted = false
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, chatID, limit, offset)
	if err != nil {
		r.logger.WithError(err).WithField("chat_id", chatID).Error("Failed to get messages")
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.ID, &message.ChatID, &message.UserID, &message.Username,
			&message.Content, &message.MessageType, &message.ReplyToID,
			&message.EditedAt, &message.IsDeleted, &message.DeletedAt,
			&message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan message")
			continue
		}
		messages = append(messages, message)
	}

	return messages, total, nil
}
