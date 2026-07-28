package sprint

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, sprint *Sprint) error
	FindByID(ctx context.Context, id uuid.UUID) (*Sprint, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]Sprint, error)
	FindActiveByProjectID(ctx context.Context, projectID uuid.UUID) (*Sprint, error)
	Update(ctx context.Context, sprint *Sprint) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
