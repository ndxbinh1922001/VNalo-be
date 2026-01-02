package entity

import "errors"

var (
	ErrInvalidUsername     = errors.New("invalid username")
	ErrUserAlreadyActive   = errors.New("user is already active")
	ErrUserAlreadyDisabled = errors.New("user is already disabled")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
)

