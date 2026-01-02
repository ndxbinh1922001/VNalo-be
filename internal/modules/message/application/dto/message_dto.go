package dto

import (
	"time"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/message/domain/entity"
)

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	ConversationID  int64             `json:"conversation_id" validate:"required"`
	Type            string            `json:"type" validate:"required,oneof=TEXT IMAGE VIDEO AUDIO FILE LOCATION"`
	Content         string            `json:"content" validate:"required"`
	ParentMessageID *string           `json:"parent_message_id,omitempty"` // UUID string for reply
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// MessageResponse represents a message response
type MessageResponse struct {
	MessageID       string            `json:"message_id"`
	ConversationID  int64             `json:"conversation_id"`
	SenderID        int64             `json:"sender_id"`
	ParentMessageID *string           `json:"parent_message_id,omitempty"`
	Type            string            `json:"type"`
	Content         string            `json:"content"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	Status          string            `json:"status"`
	IsEdited        bool              `json:"is_edited"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	ReactionCounts  map[string]int64  `json:"reaction_counts,omitempty"`
}

// NewMessageResponse converts domain entity to DTO
func NewMessageResponse(message *entity.Message) *MessageResponse {
	resp := &MessageResponse{
		MessageID:      message.MessageID.String(),
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		Type:           string(message.Type),
		Content:        message.Content,
		Metadata:       message.Metadata,
		Status:         string(message.Status),
		IsEdited:       message.IsEdited,
		CreatedAt:      message.CreatedAt,
		UpdatedAt:      message.UpdatedAt,
	}

	if message.ParentMessageID != nil {
		parentID := message.ParentMessageID.String()
		resp.ParentMessageID = &parentID
	}

	return resp
}

// MessageListResponse represents a paginated list of messages
type MessageListResponse struct {
	Messages       []*MessageResponse `json:"messages"`
	NextPageState  string             `json:"next_page_state,omitempty"`
	HasMore        bool               `json:"has_more"`
	ConversationID int64              `json:"conversation_id"`
}

// EditMessageRequest represents the request to edit a message
type EditMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

// ReactToMessageRequest represents the request to react to a message
type ReactToMessageRequest struct {
	Emoji string `json:"emoji" validate:"required,min=1,max=10"`
}

// TypingIndicatorRequest represents typing status update
type TypingIndicatorRequest struct {
	IsTyping bool `json:"is_typing"`
}

