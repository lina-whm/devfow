package project

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusArchived Status = "archived"
	StatusFrozen   Status = "frozen"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusArchived, StatusFrozen:
		return true
	}
	return false
}

type Project struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Key            string     `json:"key"`
	Description    *string    `json:"description,omitempty"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	LeadID         *uuid.UUID `json:"lead_id,omitempty"`
	Status         Status     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func NewProject(name, key, description string, organizationID uuid.UUID, leadID *uuid.UUID) *Project {
	var desc *string
	if description != "" {
		desc = &description
	}
	return &Project{
		ID:             uuid.New(),
		Name:           name,
		Key:            key,
		Description:    desc,
		OrganizationID: organizationID,
		LeadID:         leadID,
		Status:         StatusActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	nameLen := utf8.RuneCountInString(p.Name)
	if nameLen < 2 || nameLen > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	if p.Key == "" {
		return errors.New("key is required")
	}
	if len(p.Key) < 2 || len(p.Key) > 10 {
		return errors.New("key must be between 2 and 10 characters")
	}
	if p.OrganizationID == uuid.Nil {
		return errors.New("organization is required")
	}
	if !p.Status.IsValid() {
		return errors.New("invalid project status")
	}
	return nil
}

func (p *Project) Archive() {
	p.Status = StatusArchived
	p.UpdatedAt = time.Now()
}

func (p *Project) IsActive() bool {
	return p.Status == StatusActive
}

func (p *Project) IsDeleted() bool {
	return p.DeletedAt != nil
}
