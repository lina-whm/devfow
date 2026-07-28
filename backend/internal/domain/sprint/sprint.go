package sprint

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPlanned   Status = "planned"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPlanned, StatusActive, StatusCompleted:
		return true
	}
	return false
}

type Sprint struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	Name      string     `json:"name"`
	Goal      *string    `json:"goal,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewSprint(projectID uuid.UUID, name string, startDate, endDate *time.Time) *Sprint {
	return &Sprint{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      name,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    StatusPlanned,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Sprint) Validate() error {
	if s.ProjectID == uuid.Nil {
		return errors.New("project is required")
	}
	if s.Name == "" {
		return errors.New("name is required")
	}
	nameLen := utf8.RuneCountInString(s.Name)
	if nameLen < 1 || nameLen > 100 {
		return errors.New("name must be between 1 and 100 characters")
	}
	if !s.Status.IsValid() {
		return errors.New("invalid sprint status")
	}
	if s.StartDate != nil && s.EndDate != nil && s.EndDate.Before(*s.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

func (s *Sprint) Start() error {
	if s.Status != StatusPlanned {
		return errors.New("only planned sprints can be started")
	}
	now := time.Now()
	s.Status = StatusActive
	s.StartDate = &now
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Sprint) Complete() error {
	if s.Status != StatusActive {
		return errors.New("only active sprints can be completed")
	}
	now := time.Now()
	s.Status = StatusCompleted
	s.EndDate = &now
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Sprint) IsActive() bool {
	return s.Status == StatusActive
}

func (s *Sprint) IsDeleted() bool {
	return s.DeletedAt != nil
}
