package repository

import (
	"context"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/entity"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/value_object"
)

// UserRepository defines the contract for user data access
// This is a PORT in hexagonal architecture
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email value_object.Email) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context, email value_object.Email) (bool, error)
}

