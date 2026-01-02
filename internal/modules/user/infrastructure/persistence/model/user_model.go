package model

import (
	"time"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/entity"
	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/domain/value_object"
)

// UserModel is the GORM model for database persistence
type UserModel struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Email         string    `gorm:"column:email;unique;not null;index"`
	PasswordHash  string    `gorm:"column:password_hash;not null"`
	Username      string    `gorm:"column:username;not null"`
	Status        int       `gorm:"column:status;not null;default:1"`
	Language      string    `gorm:"column:language;not null;default:en"`
	IsVIP         bool      `gorm:"column:is_vip;not null;default:false"`
	LastLoginTime int64     `gorm:"column:last_login_time;default:0"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
	IsDeleted     bool      `gorm:"column:is_deleted;not null;default:false;index"`
}

func (UserModel) TableName() string {
	return "users"
}

// ToEntity converts GORM model to domain entity
func (m *UserModel) ToEntity() (*entity.User, error) {
	email, err := value_object.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}

	password := value_object.NewPasswordFromHash(m.PasswordHash)

	return &entity.User{
		ID:            m.ID,
		Email:         email,
		Password:      password,
		Username:      m.Username,
		Status:        entity.UserStatus(m.Status),
		Language:      m.Language,
		IsVIP:         m.IsVIP,
		LastLoginTime: m.LastLoginTime,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		IsDeleted:     m.IsDeleted,
	}, nil
}

// FromEntity converts domain entity to GORM model
func FromEntity(user *entity.User) *UserModel {
	return &UserModel{
		ID:            user.ID,
		Email:         user.Email.Value(),
		PasswordHash:  user.Password.Hash(),
		Username:      user.Username,
		Status:        int(user.Status),
		Language:      user.Language,
		IsVIP:         user.IsVIP,
		LastLoginTime: user.LastLoginTime,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		IsDeleted:     user.IsDeleted,
	}
}

