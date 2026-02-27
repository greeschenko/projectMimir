package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/greeschenko/projectMimir/auth-service/internal/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

type PostgresUserRepository struct {
	db *gorm.DB
}

var _ user.Repository = (*PostgresUserRepository)(nil)

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, u *user.User) error {
	model := ToModel(u)
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).
		First(&model, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return ToDomain(&model), nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).
		First(&model, "email = ?", email).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return ToDomain(&model), nil
}

func ToModel(u *user.User) UserModel {
	return UserModel{
		ID:        u.ID(),
		Email:     u.Email(),
		Password:  u.Password(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}

func ToDomain(m *UserModel) *user.User {
	return user.Restore(
		m.ID,
		m.Email,
		m.Password,
		m.CreatedAt,
		m.UpdatedAt,
	)
}
