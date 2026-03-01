package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/greeschenko/projectMimir/auth-service/internal/domain/user"
)

// ------------------------
// Fake Repository
// ------------------------

type fakeRepo struct {
	findByEmail func(email string) (*domain.User, error)
	save        func(u *domain.User) error
}

func (f *fakeRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.findByEmail(email)
}

func (f *fakeRepo) Save(ctx context.Context, u *domain.User) error {
	return f.save(u)
}

// REQUIRED because Repository interface likely has it
func (f *fakeRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

// ------------------------
// Fake Hasher
// ------------------------

type fakeHasher struct{}

func (f *fakeHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (f *fakeHasher) Compare(hash, password string) error {
	if hash != "hashed-"+password {
		return errors.New("password mismatch")
	}
	return nil
}

// ------------------------
// Fake Token Service
// ------------------------

type fakeToken struct{}

func (f *fakeToken) GenerateAccess(userID string) (string, error) {
	return "access-" + userID, nil
}

func (f *fakeToken) GenerateRefresh(userID string) (string, error) {
	return "refresh-" + userID, nil
}

// REQUIRED if TokenService has it
func (f *fakeToken) ValidateToken(token string) (string, error) {
	return "user-id", nil
}

// ------------------------
// Tests
// ------------------------

func TestRegisterUseCase_Success(t *testing.T) {
	repo := &fakeRepo{
		findByEmail: func(email string) (*domain.User, error) {
			return nil, nil
		},
		save: func(u *domain.User) error {
			return nil
		},
	}

	hasher := &fakeHasher{}
	tokens := &fakeToken{}

	uc := NewRegisterUseCase(repo, hasher, tokens)

	cmd := RegisterCommand{
		Email:    "test@example.com",
		Password: "password123",
	}

	res, err := uc.Execute(context.Background(), cmd)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", res.Email)
	}

	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Errorf("expected tokens to be set")
	}
}

func TestRegisterUseCase_UserAlreadyExists(t *testing.T) {
	// IMPORTANT: create via constructor
	existingUser, _ := domain.New("exist@example.com", "hashed-password")

	repo := &fakeRepo{
		findByEmail: func(email string) (*domain.User, error) {
			return existingUser, nil
		},
		save: func(u *domain.User) error {
			return nil
		},
	}

	hasher := &fakeHasher{}
	tokens := &fakeToken{}

	uc := NewRegisterUseCase(repo, hasher, tokens)

	cmd := RegisterCommand{
		Email:    "exist@example.com",
		Password: "password123",
	}

	_, err := uc.Execute(context.Background(), cmd)
	if err != ErrUserAlreadyExists {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestRegisterUseCase_SaveError(t *testing.T) {
	repo := &fakeRepo{
		findByEmail: func(email string) (*domain.User, error) {
			return nil, nil
		},
		save: func(u *domain.User) error {
			return errors.New("save failed")
		},
	}

	hasher := &fakeHasher{}
	tokens := &fakeToken{}

	uc := NewRegisterUseCase(repo, hasher, tokens)

	cmd := RegisterCommand{
		Email:    "new@example.com",
		Password: "password123",
	}

	_, err := uc.Execute(context.Background(), cmd)
	if err == nil || err.Error() != "save failed" {
		t.Errorf("expected save failed error, got %v", err)
	}
}
