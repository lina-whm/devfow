package team

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, team *Team) error
	FindByID(ctx context.Context, id uuid.UUID) (*Team, error)
	FindByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]Team, error)
	Update(ctx context.Context, team *Team) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type MemberRepository interface {
	AddMember(ctx context.Context, member *TeamMember) error
	FindByTeamID(ctx context.Context, teamID uuid.UUID) ([]TeamMember, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]TeamMember, error)
	FindMember(ctx context.Context, teamID, userID uuid.UUID) (*TeamMember, error)
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
}
