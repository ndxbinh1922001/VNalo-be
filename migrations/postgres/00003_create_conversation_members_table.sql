-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS conversation_members (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'MEMBER' CHECK (role IN ('OWNER', 'ADMIN', 'MEMBER')),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    is_muted BOOLEAN DEFAULT FALSE,
    notification_enabled BOOLEAN DEFAULT TRUE,
    last_read_message_id VARCHAR(100),
    last_read_at TIMESTAMP,
    unread_count INT DEFAULT 0,
    UNIQUE(conversation_id, user_id)
);

CREATE INDEX idx_members_conversation ON conversation_members(conversation_id);
CREATE INDEX idx_members_user ON conversation_members(user_id);
CREATE INDEX idx_members_unread ON conversation_members(user_id, unread_count) WHERE unread_count > 0;
CREATE INDEX idx_members_active ON conversation_members(user_id) WHERE left_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS conversation_members;
-- +goose StatementEnd

