package service

import (
	"context"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/application/dto"
)

// UserService defines application use cases
type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponse, error)
	PromoteUserToVIP(ctx context.Context, id int64) error
	DemoteUserFromVIP(ctx context.Context, id int64) error
	ChangePassword(ctx context.Context, id int64, req dto.ChangePasswordRequest) error
	ActivateUser(ctx context.Context, id int64) error
	DeactivateUser(ctx context.Context, id int64) error
}

