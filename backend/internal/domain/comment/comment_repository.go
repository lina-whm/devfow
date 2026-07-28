package comment

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	FindByID(ctx context.Context, id uuid.UUID) (*Comment, error)
	FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]Comment, error)
	Update(ctx context.Context, comment *Comment) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
