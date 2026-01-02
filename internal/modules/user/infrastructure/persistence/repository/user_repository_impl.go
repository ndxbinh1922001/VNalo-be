package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/entity"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/repository"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/value_object"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/infrastructure/persistence/model"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository implementation
// This is an ADAPTER in hexagonal architecture
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	userModel := model.FromEntity(user)

	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set the generated ID back to the entity
	user.ID = userModel.ID
	return nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var userModel model.UserModel

	err := r.db.WithContext(ctx).
		Where("id = ? AND is_deleted = ?", id, false).
		First(&userModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return userModel.ToEntity()
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email value_object.Email) (*entity.User, error) {
	var userModel model.UserModel

	err := r.db.WithContext(ctx).
		Where("email = ? AND is_deleted = ?", email.Value(), false).
		First(&userModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return userModel.ToEntity()
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	userModel := model.FromEntity(user)

	result := r.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", user.ID).
		Updates(userModel)

	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&model.UserModel{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	var userModels []model.UserModel

	err := r.db.WithContext(ctx).
		Where("is_deleted = ?", false).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&userModels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*entity.User, 0, len(userModels))
	for _, userModel := range userModels {
		user, err := userModel.ToEntity()
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("is_deleted = ?", false).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

func (r *userRepositoryImpl) Exists(ctx context.Context, email value_object.Email) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("email = ? AND is_deleted = ?", email.Value(), false).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return count > 0, nil
}

