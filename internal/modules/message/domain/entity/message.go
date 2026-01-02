package entity

import (
	"time"

	"github.com/gocql/gocql"
)

type MessageType string

const (
	MessageTypeText     MessageType = "TEXT"
	MessageTypeImage    MessageType = "IMAGE"
	MessageTypeVideo    MessageType = "VIDEO"
	MessageTypeAudio    MessageType = "AUDIO"
	MessageTypeFile     MessageType = "FILE"
	MessageTypeLocation MessageType = "LOCATION"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "SENT"
	MessageStatusDelivered MessageStatus = "DELIVERED"
	MessageStatusRead      MessageStatus = "READ"
)

type Message struct {
	MessageID       gocql.UUID  // TIMEUUID from Cassandra
	ConversationID  int64
	SenderID        int64
	ParentMessageID *gocql.UUID // For replies
	Type            MessageType
	Content         string
	Metadata        map[string]string // Store media URLs, file info, etc.
	Status          MessageStatus
	IsEdited        bool
	IsDeleted       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewMessage creates a new message
func NewMessage(conversationID, senderID int64, msgType MessageType, content string) *Message {
	return &Message{
		MessageID:      gocql.TimeUUID(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Type:           msgType,
		Content:        content,
		Metadata:       make(map[string]string),
		Status:         MessageStatusSent,
		IsEdited:       false,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// NewReplyMessage creates a reply to another message
func NewReplyMessage(conversationID, senderID int64, parentMessageID gocql.UUID, content string) *Message {
	msg := NewMessage(conversationID, senderID, MessageTypeText, content)
	msg.ParentMessageID = &parentMessageID
	return msg
}

// Domain methods

func (m *Message) Edit(newContent string) {
	m.Content = newContent
	m.IsEdited = true
	m.UpdatedAt = time.Now()
}

func (m *Message) Delete() {
	m.IsDeleted = true
	m.UpdatedAt = time.Now()
}

func (m *Message) MarkAsDelivered() {
	m.Status = MessageStatusDelivered
	m.UpdatedAt = time.Now()
}

func (m *Message) MarkAsRead() {
	m.Status = MessageStatusRead
	m.UpdatedAt = time.Now()
}

func (m *Message) IsTextMessage() bool {
	return m.Type == MessageTypeText
}

func (m *Message) IsMediaMessage() bool {
	return m.Type == MessageTypeImage ||
		m.Type == MessageTypeVideo ||
		m.Type == MessageTypeAudio ||
		m.Type == MessageTypeFile
}

func (m *Message) IsReply() bool {
	return m.ParentMessageID != nil
}

func (m *Message) AddMetadata(key, value string) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
	m.Metadata[key] = value
	m.UpdatedAt = time.Now()
}

