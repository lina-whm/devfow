package team

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Team struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Description    *string    `json:"description,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func NewTeam(name string, organizationID uuid.UUID) *Team {
	return &Team{
		ID:             uuid.New(),
		Name:           name,
		OrganizationID: organizationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (t *Team) Validate() error {
	if t.Name == "" {
		return errors.New("name is required")
	}
	nameLen := utf8.RuneCountInString(t.Name)
	if nameLen < 2 || nameLen > 50 {
		return errors.New("name must be between 2 and 50 characters")
	}
	if t.OrganizationID == uuid.Nil {
		return errors.New("organization is required")
	}
	return nil
}

func (t *Team) IsDeleted() bool {
	return t.DeletedAt != nil
}
