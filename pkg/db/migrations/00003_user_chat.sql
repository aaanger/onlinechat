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

-- Create indexes
CREATE INDEX idx_user_chat_user_id ON user_chat(user_id);
CREATE INDEX idx_user_chat_chat_id ON user_chat(chat_id);
CREATE INDEX idx_user_chat_joined_at ON user_chat(joined_at);
CREATE INDEX idx_user_chat_role ON user_chat(role);
CREATE INDEX idx_user_chat_is_banned ON user_chat(is_banned);

-- Add trigger to update current_members count
CREATE OR REPLACE FUNCTION update_chat_member_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE chats SET current_members = current_members + 1 WHERE id = NEW.chat_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE chats SET current_members = current_members - 1 WHERE id = OLD.chat_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_chat_member_count
    AFTER INSERT OR DELETE ON user_chat
    FOR EACH ROW EXECUTE FUNCTION update_chat_member_count();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_chat_member_count ON user_chat;
DROP FUNCTION IF EXISTS update_chat_member_count();
DROP INDEX IF EXISTS idx_user_chat_is_banned;
DROP INDEX IF EXISTS idx_user_chat_role;
DROP INDEX IF EXISTS idx_user_chat_joined_at;
DROP INDEX IF EXISTS idx_user_chat_chat_id;
DROP INDEX IF EXISTS idx_user_chat_user_id;
DROP TABLE IF EXISTS user_chat;
-- +goose StatementEnd
