package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/devflow/devflow-backend/internal/domain/user"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *mockUserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := NewService(repo)

	repo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
		return u.Email == "test@example.com" && u.DisplayName == "Test"
	})).Return(nil)

	result, err := svc.Register(context.Background(), "test@example.com", "password123", "Test")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	repo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := new(mockUserRepo)
	svc := NewService(repo)

	existing := user.NewUser("dup@example.com", "hash", "Dup")
	repo.On("FindByEmail", mock.Anything, "dup@example.com").Return(existing, nil)

	result, err := svc.Register(context.Background(), "dup@example.com", "password123", "Dup")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrEmailAlreadyExists, err)
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := NewService(repo)

	u := user.NewUser("test@example.com", "$2a$10$dummyhash", "Test")
	repo.On("FindByEmail", mock.Anything, "test@example.com").Return(u, nil)

	result, err := svc.Login(context.Background(), "test@example.com", "wrong")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mockUserRepo)
	svc := NewService(repo)

	repo.On("FindByEmail", mock.Anything, "notfound@example.com").Return(nil, assert.AnError)

	result, err := svc.Login(context.Background(), "notfound@example.com", "password")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrInvalidCredentials, err)
}
