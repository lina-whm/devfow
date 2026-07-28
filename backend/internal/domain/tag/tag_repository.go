package tag

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, tag *Tag) error
	FindByID(ctx context.Context, id uuid.UUID) (*Tag, error)
	FindAll(ctx context.Context) ([]Tag, error)
	Update(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
}
