package activity_log

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID         uuid.UUID              `json:"id"`
	ActorID    uuid.UUID              `json:"actor_id"`
	EntityType string                 `json:"entity_type"`
	EntityID   uuid.UUID              `json:"entity_id"`
	Action     string                 `json:"action"`
	Changes    map[string]interface{} `json:"changes,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

func NewActivityLog(actorID uuid.UUID, entityType string, entityID uuid.UUID, action string, changes map[string]interface{}) *ActivityLog {
	return &ActivityLog{
		ID:         uuid.New(),
		ActorID:    actorID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Changes:    changes,
		CreatedAt:  time.Now(),
	}
}

func (a *ActivityLog) Validate() error {
	if a.ActorID == uuid.Nil {
		return errors.New("actor is required")
	}
	if a.EntityType == "" {
		return errors.New("entity type is required")
	}
	if a.EntityID == uuid.Nil {
		return errors.New("entity id is required")
	}
	if a.Action == "" {
		return errors.New("action is required")
	}
	return nil
}
