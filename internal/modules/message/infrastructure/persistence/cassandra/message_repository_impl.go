package cassandra

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/message/domain/entity"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/message/domain/repository"
)

type messageRepositoryImpl struct {
	session *gocql.Session
}

func NewMessageRepository(session *gocql.Session) repository.MessageRepository {
	return &messageRepositoryImpl{
		session: session,
	}
}

func (r *messageRepositoryImpl) Create(ctx context.Context, message *entity.Message) error {
	query := `INSERT INTO messages (
		conversation_id, message_id, sender_id, parent_message_id,
		message_type, content, metadata, status, is_edited, is_deleted,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err := r.session.Query(query,
		message.ConversationID,
		message.MessageID,
		message.SenderID,
		message.ParentMessageID,
		string(message.Type),
		message.Content,
		message.Metadata,
		string(message.Status),
		message.IsEdited,
		message.IsDeleted,
		message.CreatedAt,
		message.UpdatedAt,
	).WithContext(ctx).Exec()

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Also insert into messages_by_user for search
	go r.insertMessageByUser(context.Background(), message)

	return nil
}

func (r *messageRepositoryImpl) insertMessageByUser(ctx context.Context, message *entity.Message) error {
	query := `INSERT INTO messages_by_user (
		user_id, message_id, conversation_id, sender_id, content, created_at
	) VALUES (?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		message.SenderID,
		message.MessageID,
		message.ConversationID,
		message.SenderID,
		message.Content,
		message.CreatedAt,
	).WithContext(ctx).Exec()
}

func (r *messageRepositoryImpl) GetByID(ctx context.Context, conversationID int64, messageID gocql.UUID) (*entity.Message, error) {
	query := `SELECT 
		conversation_id, message_id, sender_id, parent_message_id,
		message_type, content, metadata, status, is_edited, is_deleted,
		created_at, updated_at
	FROM messages 
	WHERE conversation_id = ? AND message_id = ?`

	message := &entity.Message{}
	var msgType, status string

	err := r.session.Query(query, conversationID, messageID).
		WithContext(ctx).
		Scan(
			&message.ConversationID,
			&message.MessageID,
			&message.SenderID,
			&message.ParentMessageID,
			&msgType,
			&message.Content,
			&message.Metadata,
			&status,
			&message.IsEdited,
			&message.IsDeleted,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, entity.ErrMessageNotFound
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	message.Type = entity.MessageType(msgType)
	message.Status = entity.MessageStatus(status)

	return message, nil
}

func (r *messageRepositoryImpl) GetByConversation(ctx context.Context, conversationID int64, limit int, pagingState []byte) ([]*entity.Message, []byte, error) {
	query := `SELECT 
		conversation_id, message_id, sender_id, parent_message_id,
		message_type, content, metadata, status, is_edited, is_deleted,
		created_at, updated_at
	FROM messages 
	WHERE conversation_id = ?
	LIMIT ?`

	iter := r.session.Query(query, conversationID, limit).
		WithContext(ctx).
		PageSize(limit).
		PageState(pagingState).
		Iter()

	var messages []*entity.Message

	for {
		message := &entity.Message{}
		var msgType, status string

		if !iter.Scan(
			&message.ConversationID,
			&message.MessageID,
			&message.SenderID,
			&message.ParentMessageID,
			&msgType,
			&message.Content,
			&message.Metadata,
			&status,
			&message.IsEdited,
			&message.IsDeleted,
			&message.CreatedAt,
			&message.UpdatedAt,
		) {
			break
		}

		message.Type = entity.MessageType(msgType)
		message.Status = entity.MessageStatus(status)

		if !message.IsDeleted {
			messages = append(messages, message)
		}
	}

	nextPageState := iter.PageState()

	if err := iter.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nextPageState, nil
}

func (r *messageRepositoryImpl) Update(ctx context.Context, message *entity.Message) error {
	query := `UPDATE messages SET 
		content = ?, 
		is_edited = ?, 
		is_deleted = ?,
		updated_at = ?
	WHERE conversation_id = ? AND message_id = ?`

	err := r.session.Query(query,
		message.Content,
		message.IsEdited,
		message.IsDeleted,
		time.Now(),
		message.ConversationID,
		message.MessageID,
	).WithContext(ctx).Exec()

	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}

func (r *messageRepositoryImpl) Delete(ctx context.Context, conversationID int64, messageID gocql.UUID) error {
	// Soft delete by updating is_deleted flag
	query := `UPDATE messages SET 
		is_deleted = ?,
		updated_at = ?
	WHERE conversation_id = ? AND message_id = ?`

	err := r.session.Query(query, true, time.Now(), conversationID, messageID).
		WithContext(ctx).Exec()

	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (r *messageRepositoryImpl) GetMessagesByUser(ctx context.Context, userID int64, limit int) ([]*entity.Message, error) {
	query := `SELECT 
		user_id, message_id, conversation_id, sender_id, content, created_at
	FROM messages_by_user 
	WHERE user_id = ?
	LIMIT ?`

	iter := r.session.Query(query, userID, limit).
		WithContext(ctx).
		Iter()

	var messages []*entity.Message

	for {
		message := &entity.Message{}
		var userID int64

		if !iter.Scan(
			&userID,
			&message.MessageID,
			&message.ConversationID,
			&message.SenderID,
			&message.Content,
			&message.CreatedAt,
		) {
			break
		}

		messages = append(messages, message)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get messages by user: %w", err)
	}

	return messages, nil
}

func (r *messageRepositoryImpl) UpdateDeliveryStatus(ctx context.Context, messageID gocql.UUID, userID int64, status entity.MessageStatus) error {
	query := `INSERT INTO message_delivery (
		message_id, user_id, status, delivered_at, read_at
	) VALUES (?, ?, ?, ?, ?)`

	var deliveredAt, readAt *time.Time
	now := time.Now()

	if status == entity.MessageStatusDelivered {
		deliveredAt = &now
	} else if status == entity.MessageStatusRead {
		deliveredAt = &now
		readAt = &now
	}

	return r.session.Query(query, messageID, userID, string(status), deliveredAt, readAt).
		WithContext(ctx).Exec()
}

func (r *messageRepositoryImpl) GetDeliveryStatus(ctx context.Context, messageID gocql.UUID) (map[int64]entity.MessageStatus, error) {
	query := `SELECT user_id, status FROM message_delivery WHERE message_id = ?`

	iter := r.session.Query(query, messageID).
		WithContext(ctx).
		Iter()

	statuses := make(map[int64]entity.MessageStatus)

	for {
		var userID int64
		var status string

		if !iter.Scan(&userID, &status) {
			break
		}

		statuses[userID] = entity.MessageStatus(status)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get delivery status: %w", err)
	}

	return statuses, nil
}

func (r *messageRepositoryImpl) AddReaction(ctx context.Context, messageID gocql.UUID, userID int64, emoji string) error {
	// Insert reaction
	query1 := `INSERT INTO message_reactions (message_id, user_id, emoji, created_at) 
			   VALUES (?, ?, ?, ?)`

	if err := r.session.Query(query1, messageID, userID, emoji, time.Now()).
		WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	// Increment counter
	query2 := `UPDATE reaction_counts SET count = count + 1 
			   WHERE message_id = ? AND emoji = ?`

	if err := r.session.Query(query2, messageID, emoji).
		WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to update reaction count: %w", err)
	}

	return nil
}

func (r *messageRepositoryImpl) RemoveReaction(ctx context.Context, messageID gocql.UUID, userID int64, emoji string) error {
	// Delete reaction
	query1 := `DELETE FROM message_reactions 
			   WHERE message_id = ? AND user_id = ? AND emoji = ?`

	if err := r.session.Query(query1, messageID, userID, emoji).
		WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	// Decrement counter
	query2 := `UPDATE reaction_counts SET count = count - 1 
			   WHERE message_id = ? AND emoji = ?`

	if err := r.session.Query(query2, messageID, emoji).
		WithContext(ctx).Exec(); err != nil {
		return fmt.Errorf("failed to update reaction count: %w", err)
	}

	return nil
}

func (r *messageRepositoryImpl) GetReactions(ctx context.Context, messageID gocql.UUID) (map[string][]int64, error) {
	query := `SELECT emoji, user_id FROM message_reactions WHERE message_id = ?`

	iter := r.session.Query(query, messageID).
		WithContext(ctx).
		Iter()

	reactions := make(map[string][]int64)

	for {
		var emoji string
		var userID int64

		if !iter.Scan(&emoji, &userID) {
			break
		}

		reactions[emoji] = append(reactions[emoji], userID)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get reactions: %w", err)
	}

	return reactions, nil
}

func (r *messageRepositoryImpl) GetReactionCounts(ctx context.Context, messageID gocql.UUID) (map[string]int64, error) {
	query := `SELECT emoji, count FROM reaction_counts WHERE message_id = ?`

	iter := r.session.Query(query, messageID).
		WithContext(ctx).
		Iter()

	counts := make(map[string]int64)

	for {
		var emoji string
		var count int64

		if !iter.Scan(&emoji, &count) {
			break
		}

		counts[emoji] = count
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get reaction counts: %w", err)
	}

	return counts, nil
}

func (r *messageRepositoryImpl) SetTypingIndicator(ctx context.Context, conversationID int64, userID int64, isTyping bool) error {
	query := `INSERT INTO typing_indicators (conversation_id, user_id, is_typing, updated_at) 
			  VALUES (?, ?, ?, ?) USING TTL 10`

	return r.session.Query(query, conversationID, userID, isTyping, time.Now()).
		WithContext(ctx).Exec()
}

func (r *messageRepositoryImpl) GetTypingUsers(ctx context.Context, conversationID int64) ([]int64, error) {
	query := `SELECT user_id FROM typing_indicators 
			  WHERE conversation_id = ? AND is_typing = true`

	iter := r.session.Query(query, conversationID).
		WithContext(ctx).
		Iter()

	var userIDs []int64

	for {
		var userID int64
		if !iter.Scan(&userID) {
			break
		}
		userIDs = append(userIDs, userID)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get typing users: %w", err)
	}

	return userIDs, nil
}

