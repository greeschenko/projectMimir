package user

import (
	"context"
	"errors"

	"github.com/greeschenko/projectMimir/auth-service/internal/domain/user"
	"github.com/greeschenko/projectMimir/auth-service/internal/ports"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type LoginUseCase struct {
	repo   user.Repository
	hasher ports.PasswordHasher
	tokens ports.TokenService
}

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	UserID       string
	Email        string
	AccessToken  string
	RefreshToken string
}

func NewLoginUseCase(repo user.Repository, hasher ports.PasswordHasher, tokens ports.TokenService) *LoginUseCase {
	return &LoginUseCase{repo: repo, hasher: hasher, tokens: tokens}
}

func (uc *LoginUseCase) Execute(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	u, err := uc.repo.FindByEmail(ctx, cmd.Email)
	if err != nil || u == nil {
		return nil, ErrInvalidCredentials
	}

	if err := uc.hasher.Compare(u.Password(), cmd.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	access, err := uc.tokens.GenerateAccess(u.ID().String())
	if err != nil {
		return nil, err
	}

	refresh, err := uc.tokens.GenerateRefresh(u.ID().String())
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		UserID:       u.ID().String(),
		Email:        u.Email(),
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
