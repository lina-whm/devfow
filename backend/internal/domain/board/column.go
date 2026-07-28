package board

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Column struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	Name      string    `json:"name"`
	Position  float64   `json:"position"`
	WIPLimit  int       `json:"wip_limit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewColumn(boardID uuid.UUID, name string, position float64) *Column {
	return &Column{
		ID:        uuid.New(),
		BoardID:   boardID,
		Name:      name,
		Position:  position,
		WIPLimit:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c *Column) Validate() error {
	if c.BoardID == uuid.Nil {
		return errors.New("board is required")
	}
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.WIPLimit < 0 {
		return errors.New("wip limit must be zero or positive")
	}
	return nil
}

func (c *Column) SetWIPLimit(limit int) error {
	if limit < 0 {
		return errors.New("wip limit must be zero or positive")
	}
	c.WIPLimit = limit
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Column) MoveToPosition(position float64) {
	c.Position = position
	c.UpdatedAt = time.Now()
}
