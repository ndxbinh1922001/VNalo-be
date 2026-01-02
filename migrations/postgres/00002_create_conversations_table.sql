-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS conversations (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('DIRECT', 'GROUP', 'CHANNEL')),
    name VARCHAR(255),
    avatar_url VARCHAR(500),
    description TEXT,
    creator_id BIGINT NOT NULL REFERENCES users(id),
    is_public BOOLEAN DEFAULT FALSE,
    member_count INT DEFAULT 0,
    last_message_id VARCHAR(100),
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_conversations_creator ON conversations(creator_id);
CREATE INDEX idx_conversations_type ON conversations(type);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC NULLS LAST);
CREATE INDEX idx_conversations_is_deleted ON conversations(is_deleted);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS conversations;
-- +goose StatementEnd

