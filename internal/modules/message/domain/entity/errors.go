package entity

import "errors"

var (
	ErrMessageNotFound        = errors.New("message not found")
	ErrInvalidMessageType     = errors.New("invalid message type")
	ErrEmptyMessageContent    = errors.New("message content cannot be empty")
	ErrConversationNotFound   = errors.New("conversation not found")
	ErrUserNotInConversation  = errors.New("user is not a member of this conversation")
	ErrCannotEditMessage      = errors.New("cannot edit this message")
	ErrCannotDeleteMessage    = errors.New("cannot delete this message")
	ErrInvalidReaction        = errors.New("invalid reaction emoji")
)

