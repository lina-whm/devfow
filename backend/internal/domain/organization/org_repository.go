package organization

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, org *Organization) error
	FindByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	FindBySlug(ctx context.Context, slug string) (*Organization, error)
	Update(ctx context.Context, org *Organization) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type MemberRepository interface {
	AddMember(ctx context.Context, member *OrganizationMember) error
	FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]OrganizationMember, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]OrganizationMember, error)
	FindMember(ctx context.Context, orgID, userID uuid.UUID) (*OrganizationMember, error)
	UpdateRole(ctx context.Context, orgID, userID uuid.UUID, role Role) error
	RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error
}

type InvitationRepository interface {
	Create(ctx context.Context, invitation *Invitation) error
	FindByID(ctx context.Context, id uuid.UUID) (*Invitation, error)
	FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]Invitation, error)
	FindByEmail(ctx context.Context, email string) ([]Invitation, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*Invitation, error)
	Update(ctx context.Context, invitation *Invitation) error
	Delete(ctx context.Context, id uuid.UUID) error
}
