package dto

import (
	"time"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/entity"
)

// CreateUserRequest represents the input for creating a user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Username string `json:"username" validate:"required,min=3"`
}

// UpdateUserRequest represents the input for updating a user
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3"`
	Language *string `json:"language,omitempty" validate:"omitempty,min=2,max=5"`
	Status   *int    `json:"status,omitempty" validate:"omitempty,oneof=1 2"`
}

// ChangePasswordRequest represents the input for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UserResponse represents the output for user data
type UserResponse struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Status        int       `json:"status"`
	Language      string    `json:"language"`
	IsVIP         bool      `json:"is_vip"`
	LastLoginTime int64     `json:"last_login_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewUserResponse converts domain entity to DTO
func NewUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:            user.ID,
		Email:         user.Email.Value(),
		Username:      user.Username,
		Status:        int(user.Status),
		Language:      user.Language,
		IsVIP:         user.IsVIP,
		LastLoginTime: user.LastLoginTime,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
}

