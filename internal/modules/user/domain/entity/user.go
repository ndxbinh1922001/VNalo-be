package entity

import (
	"time"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/value_object"
)

// User entity represents the core business object
type User struct {
	ID            int64
	Email         value_object.Email
	Password      value_object.Password
	Username      string
	Status        UserStatus
	Language      string
	IsVIP         bool
	LastLoginTime int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsDeleted     bool
}

type UserStatus int

const (
	UserStatusActive UserStatus = iota + 1
	UserStatusDisabled
)

// NewUser factory method to create a new user
func NewUser(email value_object.Email, password value_object.Password, username string) (*User, error) {
	// Domain validation logic
	if username == "" {
		return nil, ErrInvalidUsername
	}

	return &User{
		Email:         email,
		Password:      password,
		Username:      username,
		Status:        UserStatusActive,
		Language:      "en",
		IsVIP:         false,
		LastLoginTime: 0,
		IsDeleted:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// Domain methods

// Activate activates the user account
func (u *User) Activate() error {
	if u.Status == UserStatusActive {
		return ErrUserAlreadyActive
	}
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate() error {
	if u.Status == UserStatusDisabled {
		return ErrUserAlreadyDisabled
	}
	u.Status = UserStatusDisabled
	u.UpdatedAt = time.Now()
	return nil
}

// PromoteToVIP promotes user to VIP status
func (u *User) PromoteToVIP() {
	u.IsVIP = true
	u.UpdatedAt = time.Now()
}

// DemoteFromVIP removes VIP status
func (u *User) DemoteFromVIP() {
	u.IsVIP = false
	u.UpdatedAt = time.Now()
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && !u.IsDeleted
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	u.LastLoginTime = time.Now().Unix()
	u.UpdatedAt = time.Now()
}

// ChangePassword changes user password
func (u *User) ChangePassword(newPassword value_object.Password) {
	u.Password = newPassword
	u.UpdatedAt = time.Now()
}

// ChangeLanguage changes user language preference
func (u *User) ChangeLanguage(language string) {
	u.Language = language
	u.UpdatedAt = time.Now()
}

// SoftDelete marks user as deleted
func (u *User) SoftDelete() {
	u.IsDeleted = true
	u.UpdatedAt = time.Now()
}

