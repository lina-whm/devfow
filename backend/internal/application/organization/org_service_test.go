package organization

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/devflow/devflow-backend/internal/domain/organization"
	"github.com/devflow/devflow-backend/internal/domain/user"
)

type mockOrgRepo struct{ mock.Mock }
func (m *mockOrgRepo) Create(ctx context.Context, org *organization.Organization) error {
	return m.Called(ctx, org).Error(0)
}
func (m *mockOrgRepo) FindByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*organization.Organization), args.Error(1)
}
func (m *mockOrgRepo) FindBySlug(ctx context.Context, slug string) (*organization.Organization, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*organization.Organization), args.Error(1)
}
func (m *mockOrgRepo) Update(ctx context.Context, org *organization.Organization) error {
	return m.Called(ctx, org).Error(0)
}
func (m *mockOrgRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockMemRepo struct{ mock.Mock }
func (m *mockMemRepo) AddMember(ctx context.Context, member *organization.OrganizationMember) error {
	return m.Called(ctx, member).Error(0)
}
func (m *mockMemRepo) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]organization.OrganizationMember, error) {
	args := m.Called(ctx, orgID); return args.Get(0).([]organization.OrganizationMember), args.Error(1)
}
func (m *mockMemRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]organization.OrganizationMember, error) {
	args := m.Called(ctx, userID); return args.Get(0).([]organization.OrganizationMember), args.Error(1)
}
func (m *mockMemRepo) FindMember(ctx context.Context, orgID, userID uuid.UUID) (*organization.OrganizationMember, error) {
	args := m.Called(ctx, orgID, userID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*organization.OrganizationMember), args.Error(1)
}
func (m *mockMemRepo) UpdateRole(ctx context.Context, orgID, userID uuid.UUID, role organization.Role) error {
	return m.Called(ctx, orgID, userID, role).Error(0)
}
func (m *mockMemRepo) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	return m.Called(ctx, orgID, userID).Error(0)
}

type mockInvRepo struct{ mock.Mock }
func (m *mockInvRepo) Create(ctx context.Context, inv *organization.Invitation) error { return m.Called(ctx, inv).Error(0) }
func (m *mockInvRepo) FindByID(ctx context.Context, id uuid.UUID) (*organization.Invitation, error) {
	args := m.Called(ctx, id); return args.Get(0).(*organization.Invitation), args.Error(1)
}
func (m *mockInvRepo) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]organization.Invitation, error) { return nil, nil }
func (m *mockInvRepo) FindByEmail(ctx context.Context, email string) ([]organization.Invitation, error) { return nil, nil }
func (m *mockInvRepo) FindByTokenHash(ctx context.Context, tokenHash string) (*organization.Invitation, error) { return nil, nil }
func (m *mockInvRepo) Update(ctx context.Context, inv *organization.Invitation) error { return nil }
func (m *mockInvRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

type mockUserRepo struct{ mock.Mock }
func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error { return nil }
func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error { return nil }
func (m *mockUserRepo) SoftDelete(ctx context.Context, id uuid.UUID) error { return nil }

func TestCreateOrg_Success(t *testing.T) {
	orgRepo := new(mockOrgRepo)
	memRepo := new(mockMemRepo)
	invRepo := new(mockInvRepo)
	userRepo := new(mockUserRepo)
	svc := NewService(orgRepo, memRepo, invRepo, userRepo)

	ownerID := uuid.New()

	orgRepo.On("FindBySlug", mock.Anything, "my-org").Return(nil, assert.AnError)
	orgRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	memRepo.On("AddMember", mock.Anything, mock.Anything).Return(nil)

	org, err := svc.Create(context.Background(), "My Org", "my-org", ownerID)
	assert.NoError(t, err)
	assert.NotNil(t, org)
	assert.Equal(t, "My Org", org.Name)
	assert.Equal(t, "my-org", org.Slug)
	orgRepo.AssertExpectations(t)
}

func TestCreateOrg_DuplicateSlug(t *testing.T) {
	orgRepo := new(mockOrgRepo)
	memRepo := new(mockMemRepo)
	invRepo := new(mockInvRepo)
	userRepo := new(mockUserRepo)
	svc := NewService(orgRepo, memRepo, invRepo, userRepo)

	existing := organization.NewOrganization("Existing", "my-org", "", uuid.New())
	orgRepo.On("FindBySlug", mock.Anything, "my-org").Return(existing, nil)

	_, err := svc.Create(context.Background(), "My Org", "my-org", uuid.New())
	assert.Error(t, err)
	assert.Equal(t, ErrOrgSlugTaken, err)
}
