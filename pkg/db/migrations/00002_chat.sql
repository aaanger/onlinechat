-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_private BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    max_members INT DEFAULT 100,
    current_members INT DEFAULT 0
);

CREATE INDEX idx_chats_created_by ON chats(created_by);
CREATE INDEX idx_chats_created_at ON chats(created_at);
CREATE INDEX idx_chats_is_active ON chats(is_active);
CREATE INDEX idx_chats_is_private ON chats(is_private);

ALTER TABLE chats ADD CONSTRAINT check_name_length CHECK (length(name) >= 1 AND length(name) <= 100);
ALTER TABLE chats ADD CONSTRAINT check_max_members CHECK (max_members > 0 AND max_members <= 1000);
ALTER TABLE chats ADD CONSTRAINT check_current_members CHECK (current_members >= 0 AND current_members <= max_members);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chats_is_private;
DROP INDEX IF EXISTS idx_chats_is_active;
DROP INDEX IF EXISTS idx_chats_created_at;
DROP INDEX IF EXISTS idx_chats_created_by;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd
