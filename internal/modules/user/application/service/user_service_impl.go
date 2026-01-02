package service

import (
	"context"
	"fmt"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/application/dto"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/entity"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/repository"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/value_object"
	"github.com/ndxbinh1922001/VNalo-be/pkg/validator"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user application service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 1. Validate request
	if err := validator.Validate(&req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Check if user already exists
	email, err := value_object.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	exists, err := s.userRepo.Exists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, entity.ErrEmailAlreadyExists
	}

	// 3. Create password value object
	password, err := value_object.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 4. Create domain entity
	user, err := entity.NewUser(email, password, req.Username)
	if err != nil {
		return nil, err
	}

	// 5. Persist to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. Return DTO
	return dto.NewUserResponse(user), nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return dto.NewUserResponse(user), nil
}

func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	emailVO, err := value_object.NewEmail(email)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, emailVO)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return dto.NewUserResponse(user), nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Validate request
	if err := validator.Validate(&req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Language != nil {
		user.ChangeLanguage(*req.Language)
	}
	if req.Status != nil {
		user.Status = entity.UserStatus(*req.Status)
	}

	// Save
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return dto.NewUserResponse(user), nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	// Get user first to ensure it exists
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Soft delete
	user.SoftDelete()

	// Save
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) ListUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, err := s.userRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	totalCount, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, dto.NewUserResponse(user))
	}

	return &dto.UserListResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *userServiceImpl) PromoteUserToVIP(ctx context.Context, id int64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.PromoteToVIP()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) DemoteUserFromVIP(ctx context.Context, id int64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.DemoteFromVIP()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) ChangePassword(ctx context.Context, id int64, req dto.ChangePasswordRequest) error {
	// Validate request
	if err := validator.Validate(&req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify old password
	if err := user.Password.Compare(req.OldPassword); err != nil {
		return value_object.ErrPasswordMismatch
	}

	// Create new password
	newPassword, err := value_object.NewPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Change password
	user.ChangePassword(newPassword)

	// Save
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) ActivateUser(ctx context.Context, id int64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.Activate(); err != nil {
		return err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) DeactivateUser(ctx context.Context, id int64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.Deactivate(); err != nil {
		return err
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

