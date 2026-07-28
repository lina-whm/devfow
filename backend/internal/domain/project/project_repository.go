package project

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, project *Project) error
	FindByID(ctx context.Context, id uuid.UUID) (*Project, error)
	FindByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]Project, error)
	FindByKey(ctx context.Context, orgID uuid.UUID, key string) (*Project, error)
	Update(ctx context.Context, project *Project) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
