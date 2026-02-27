package user

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type User struct {
	id        uuid.UUID
	email     string
	password  string
	createdAt time.Time
	updatedAt time.Time
}

func New(email, hashedPassword string) (*User, error) {
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	if len(hashedPassword) == 0 {
		return nil, ErrInvalidPassword
	}

	now := time.Now()

	return &User{
		id:        uuid.New(),
		email:     email,
		password:  hashedPassword,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Restore(
	id uuid.UUID,
	email string,
	password string,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		email:     email,
		password:  password,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() uuid.UUID        { return u.id }
func (u *User) Email() string        { return u.email }
func (u *User) Password() string     { return u.password }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
