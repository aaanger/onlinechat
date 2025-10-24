-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file', 'system')),
    reply_to_id INT REFERENCES messages(id) ON DELETE SET NULL,
    edited_at TIMESTAMP WITH TIME ZONE,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_chat_created ON messages(chat_id, created_at);
CREATE INDEX idx_messages_reply_to ON messages(reply_to_id);
CREATE INDEX idx_messages_type ON messages(message_type);
CREATE INDEX idx_messages_is_deleted ON messages(is_deleted);

ALTER TABLE messages ADD CONSTRAINT check_content_length CHECK (length(content) >= 1 AND length(content) <= 4000);
ALTER TABLE messages ADD CONSTRAINT check_content_not_empty CHECK (trim(content) != '');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_is_deleted;
DROP INDEX IF EXISTS idx_messages_type;
DROP INDEX IF EXISTS idx_messages_reply_to;
DROP INDEX IF EXISTS idx_messages_chat_created;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_user_id;
DROP INDEX IF EXISTS idx_messages_chat_id;
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
