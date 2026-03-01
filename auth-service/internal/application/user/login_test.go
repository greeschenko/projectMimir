package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/greeschenko/projectMimir/auth-service/internal/domain/user"
)

// ------------------------
// Fake repository
// ------------------------

type fakeLoginRepo struct {
	findByEmail func(email string) (*domain.User, error)
}

func (f *fakeLoginRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.findByEmail(email)
}

// Required by interface
func (f *fakeLoginRepo) Save(ctx context.Context, u *domain.User) error {
	return nil
}

// Required by interface
func (f *fakeLoginRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

// ------------------------
// Fake hasher
// ------------------------

type fakeLoginHasher struct{}

func (f *fakeLoginHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (f *fakeLoginHasher) Compare(hash, password string) error {
	if hash != "hashed-"+password {
		return errors.New("password mismatch")
	}
	return nil
}

// ------------------------
// Fake token service
// ------------------------

type fakeLoginToken struct{}

func (f *fakeLoginToken) GenerateAccess(userID string) (string, error) {
	return "access-" + userID, nil
}

func (f *fakeLoginToken) GenerateRefresh(userID string) (string, error) {
	return "refresh-" + userID, nil
}

// Required by interface
func (f *fakeLoginToken) ValidateToken(token string) (string, error) {
	return "user-id", nil
}

// ------------------------
// Tests
// ------------------------

func TestLoginUseCase_Success(t *testing.T) {
	userDomain, _ := domain.New("login@example.com", "hashed-password123")

	repo := &fakeLoginRepo{
		findByEmail: func(email string) (*domain.User, error) { return userDomain, nil },
	}

	hasher := &fakeLoginHasher{}
	tokens := &fakeLoginToken{}

	uc := NewLoginUseCase(repo, hasher, tokens)
	cmd := LoginCommand{
		Email:    "login@example.com",
		Password: "password123",
	}

	res, err := uc.Execute(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Email != "login@example.com" {
		t.Errorf("expected email login@example.com, got %s", res.Email)
	}

	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Errorf("expected tokens to be set")
	}
}

func TestLoginUseCase_InvalidCredentials_EmailNotFound(t *testing.T) {
	repo := &fakeLoginRepo{
		findByEmail: func(email string) (*domain.User, error) { return nil, nil }, // email not found
	}

	hasher := &fakeLoginHasher{}
	tokens := &fakeLoginToken{}

	uc := NewLoginUseCase(repo, hasher, tokens)
	cmd := LoginCommand{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	_, err := uc.Execute(context.Background(), cmd)
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginUseCase_InvalidCredentials_WrongPassword(t *testing.T) {
	userDomain, _ := domain.New("login@example.com", "hashed-password123")

	repo := &fakeLoginRepo{
		findByEmail: func(email string) (*domain.User, error) { return userDomain, nil },
	}

	hasher := &fakeLoginHasher{}
	tokens := &fakeLoginToken{}

	uc := NewLoginUseCase(repo, hasher, tokens)
	cmd := LoginCommand{
		Email:    "login@example.com",
		Password: "wrongpassword",
	}

	_, err := uc.Execute(context.Background(), cmd)
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials for wrong password, got %v", err)
	}
}
