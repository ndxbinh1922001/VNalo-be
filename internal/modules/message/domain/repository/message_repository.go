package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/message/domain/entity"
)

// MessageRepository defines the contract for message data access
type MessageRepository interface {
	// Message CRUD
	Create(ctx context.Context, message *entity.Message) error
	GetByID(ctx context.Context, conversationID int64, messageID gocql.UUID) (*entity.Message, error)
	GetByConversation(ctx context.Context, conversationID int64, limit int, pagingState []byte) ([]*entity.Message, []byte, error)
	Update(ctx context.Context, message *entity.Message) error
	Delete(ctx context.Context, conversationID int64, messageID gocql.UUID) error

	// Search
	GetMessagesByUser(ctx context.Context, userID int64, limit int) ([]*entity.Message, error)

	// Delivery status
	UpdateDeliveryStatus(ctx context.Context, messageID gocql.UUID, userID int64, status entity.MessageStatus) error
	GetDeliveryStatus(ctx context.Context, messageID gocql.UUID) (map[int64]entity.MessageStatus, error)

	// Reactions
	AddReaction(ctx context.Context, messageID gocql.UUID, userID int64, emoji string) error
	RemoveReaction(ctx context.Context, messageID gocql.UUID, userID int64, emoji string) error
	GetReactions(ctx context.Context, messageID gocql.UUID) (map[string][]int64, error)
	GetReactionCounts(ctx context.Context, messageID gocql.UUID) (map[string]int64, error)

	// Typing indicators
	SetTypingIndicator(ctx context.Context, conversationID int64, userID int64, isTyping bool) error
	GetTypingUsers(ctx context.Context, conversationID int64) ([]int64, error)
}

