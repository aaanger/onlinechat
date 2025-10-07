-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_chat (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(20) DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'moderator', 'member')),
    is_muted BOOLEAN DEFAULT false,
    is_banned BOOLEAN DEFAULT false,
    banned_until TIMESTAMP WITH TIME ZONE,
    last_read_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, chat_id)
);

CREATE INDEX idx_user_chat_user_id ON user_chat(user_id);
CREATE INDEX idx_user_chat_chat_id ON user_chat(chat_id);
CREATE INDEX idx_user_chat_joined_at ON user_chat(joined_at);
CREATE INDEX idx_user_chat_role ON user_chat(role);
CREATE INDEX idx_user_chat_is_banned ON user_chat(is_banned);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_chat_is_banned;
DROP INDEX IF EXISTS idx_user_chat_role;
DROP INDEX IF EXISTS idx_user_chat_joined_at;
DROP INDEX IF EXISTS idx_user_chat_chat_id;
DROP INDEX IF EXISTS idx_user_chat_user_id;
DROP TABLE IF EXISTS user_chat;
-- +goose StatementEnd
