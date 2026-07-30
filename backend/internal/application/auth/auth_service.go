package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/devflow/devflow-backend/internal/domain/user"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
)

type Service struct {
	userRepo user.Repository
}

func NewService(userRepo user.Repository) *Service {
	return &Service{userRepo: userRepo}
}

type AuthUser struct {
	ID            string
	Email         string
	PasswordHash  string
	DisplayName   string
	AvatarURL     string
	EmailVerified bool
	OrgID         string
	Role          string
}

func (s *Service) Register(ctx context.Context, email, password, displayName string) (*AuthUser, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := user.NewUser(email, string(hash), displayName)
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &AuthUser{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*AuthUser, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &AuthUser{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*AuthUser, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse id: %w", err)
	}
	u, err := s.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, ErrUserNotFound
	}
	avatarURL := ""
	if u.AvatarURL != nil {
		avatarURL = *u.AvatarURL
	}
	return &AuthUser{
		ID:            u.ID.String(),
		Email:         u.Email,
		DisplayName:   u.DisplayName,
		AvatarURL:     avatarURL,
		EmailVerified: u.EmailVerifiedAt != nil,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*AuthUser, error) {
	return nil, ErrInvalidToken
}

func (s *Service) VerifyEmail(ctx context.Context, email string) error {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}
	now := time.Now()
	u.EmailVerifiedAt = &now
	return s.userRepo.Update(ctx, u)
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, displayName string, avatarURL string) (*AuthUser, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("parse id: %w", err)
	}
	u, err := s.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if displayName != "" {
		u.DisplayName = displayName
	}
	if avatarURL != "" {
		u.AvatarURL = &avatarURL
	}
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	avatar := ""
	if u.AvatarURL != nil {
		avatar = *u.AvatarURL
	}
	return &AuthUser{
		ID:          u.ID.String(),
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarURL:   avatar,
	}, nil
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, email, newPassword string) error {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	u.PasswordHash = string(hash)
	return s.userRepo.Update(ctx, u)
}
