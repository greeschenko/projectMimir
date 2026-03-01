package user

import (
    "context"
    "errors"

    "github.com/greeschenko/projectMimir/auth-service/internal/domain/user"
    "github.com/greeschenko/projectMimir/auth-service/internal/ports"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type RegisterUseCase struct {
    repo   user.Repository
    hasher ports.PasswordHasher
    tokens ports.TokenService
}

type RegisterCommand struct {
    Email    string
    Password string
}

type RegisterResult struct {
    UserID       string
    Email        string
    AccessToken  string
    RefreshToken string
}

func NewRegisterUseCase(repo user.Repository, hasher ports.PasswordHasher, tokens ports.TokenService) *RegisterUseCase {
    return &RegisterUseCase{repo: repo, hasher: hasher, tokens: tokens}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, cmd RegisterCommand) (*RegisterResult, error) {
    existing, err := uc.repo.FindByEmail(ctx, cmd.Email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, ErrUserAlreadyExists
    }

    hashed, err := uc.hasher.Hash(cmd.Password)
    if err != nil {
        return nil, err
    }

    userEntity, err := user.New(cmd.Email, hashed)
    if err != nil {
        return nil, err
    }

    if err := uc.repo.Save(ctx, userEntity); err != nil {
        return nil, err
    }

    access, err := uc.tokens.GenerateAccess(userEntity.ID().String())
    if err != nil {
        return nil, err
    }
    refresh, err := uc.tokens.GenerateRefresh(userEntity.ID().String())
    if err != nil {
        return nil, err
    }

    return &RegisterResult{
        UserID:       userEntity.ID().String(),
        Email:        userEntity.Email(),
        AccessToken:  access,
        RefreshToken: refresh,
    }, nil
}
