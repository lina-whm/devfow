package board

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, board *Board) error
	FindByID(ctx context.Context, id uuid.UUID) (*Board, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID) (*Board, error)
	Update(ctx context.Context, board *Board) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ColumnRepository interface {
	Create(ctx context.Context, column *Column) error
	FindByID(ctx context.Context, id uuid.UUID) (*Column, error)
	FindByBoardID(ctx context.Context, boardID uuid.UUID) ([]Column, error)
	Update(ctx context.Context, column *Column) error
	Delete(ctx context.Context, id uuid.UUID) error
	Reorder(ctx context.Context, boardID uuid.UUID, columnIDs []uuid.UUID) error
}
