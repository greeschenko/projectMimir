package user

import (
	"context"
	"errors"

	domain "github.com/greeschenko/projectMimir/auth-service/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type RegisterUseCase struct {
	repo domain.Repository
}

func NewRegisterUseCase(repo domain.Repository) *RegisterUseCase {
	return &RegisterUseCase{repo: repo}
}

type RegisterCommand struct {
	Email    string
	Password string
}

func (uc *RegisterUseCase) Execute(ctx context.Context, cmd RegisterCommand) (*domain.User, error) {
	existing, err := uc.repo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := domain.New(cmd.Email, string(hashed))
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil // 🔥 повертаємо user
}
