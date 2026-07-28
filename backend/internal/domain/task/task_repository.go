package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskFilter struct {
	ProjectID  *uuid.UUID
	SprintID   *uuid.UUID
	ColumnID   *uuid.UUID
	AssigneeID *uuid.UUID
	Type       *Type
	Priority   *Priority
	Status     *Status
	Search     string
	Cursor     *time.Time
	Limit      int
}

type Repository interface {
	Create(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	List(ctx context.Context, filter TaskFilter) ([]Task, *time.Time, error)
	Update(ctx context.Context, task *Task) error
	UpdateWithVersion(ctx context.Context, task *Task, currentVersion int) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}
