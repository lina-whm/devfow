package board

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Board struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBoard(projectID uuid.UUID, name string) *Board {
	return &Board{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (b *Board) Validate() error {
	if b.ProjectID == uuid.Nil {
		return errors.New("project is required")
	}
	if b.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
