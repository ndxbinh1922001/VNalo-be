-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'User ID',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT 'User email address',
    password_hash VARCHAR(255) NOT NULL COMMENT 'Hashed password',
    username VARCHAR(255) NOT NULL COMMENT 'Username',
    status TINYINT NOT NULL DEFAULT 1 COMMENT 'Status: 1=active, 2=disabled',
    language VARCHAR(50) NOT NULL DEFAULT 'en' COMMENT 'Language preference',
    is_vip BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'VIP status flag',
    last_login_time BIGINT DEFAULT 0 COMMENT 'Last login timestamp',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last update time',
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Soft delete flag',
    INDEX idx_email (email),
    INDEX idx_is_deleted (is_deleted),
    INDEX idx_created_at (created_at)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='User accounts table';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

