package service

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/message/application/dto"
)

// MessageService defines application use cases for messages
type MessageService interface {
	// Send and manage messages
	SendMessage(ctx context.Context, senderID int64, req dto.SendMessageRequest) (*dto.MessageResponse, error)
	GetMessage(ctx context.Context, conversationID int64, messageID gocql.UUID) (*dto.MessageResponse, error)
	GetConversationMessages(ctx context.Context, conversationID int64, userID int64, limit int, pagingState string) (*dto.MessageListResponse, error)
	EditMessage(ctx context.Context, conversationID int64, messageID gocql.UUID, userID int64, req dto.EditMessageRequest) (*dto.MessageResponse, error)
	DeleteMessage(ctx context.Context, conversationID int64, messageID gocql.UUID, userID int64) error

	// Reactions
	AddReaction(ctx context.Context, conversationID int64, messageID gocql.UUID, userID int64, emoji string) error
	RemoveReaction(ctx context.Context, conversationID int64, messageID gocql.UUID, userID int64, emoji string) error

	// Typing indicators
	SetTyping(ctx context.Context, conversationID int64, userID int64, isTyping bool) error
	GetTypingUsers(ctx context.Context, conversationID int64) ([]int64, error)

	// Mark as read
	MarkAsRead(ctx context.Context, conversationID int64, messageID gocql.UUID, userID int64) error
}

