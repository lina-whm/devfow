package user

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	DisplayName     string     `json:"display_name"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewUser(email, passwordHash, displayName string) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		DisplayName:  displayName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.PasswordHash == "" {
		return errors.New("password hash is required")
	}
	if u.DisplayName == "" {
		return errors.New("display name is required")
	}
	nameLen := utf8.RuneCountInString(u.DisplayName)
	if nameLen < 3 || nameLen > 50 {
		return errors.New("display name must be between 3 and 50 characters")
	}
	return nil
}

func (u *User) IsVerified() bool {
	return u.EmailVerifiedAt != nil
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}
