package notification

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (*Notification, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Notification, error)
	CountUnreadByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}
