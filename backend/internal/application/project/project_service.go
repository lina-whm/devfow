package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/domain/project"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrProjectKeyTaken = errors.New("project key already taken in this organization")
)

type Service struct {
	repo project.Repository
}

func NewService(repo project.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name, key, description string, orgID uuid.UUID, leadID *uuid.UUID) (*project.Project, error) {
	if leadID != nil && *leadID == uuid.Nil {
		leadID = nil
	}

	existing, _ := s.repo.FindByKey(ctx, orgID, key)
	if existing != nil {
		return nil, ErrProjectKeyTaken
	}

	p := project.NewProject(name, key, description, orgID, leadID)
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return p, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}
	return p, nil
}

func (s *Service) ListByOrgID(ctx context.Context, orgID uuid.UUID) ([]project.Project, error) {
	return s.repo.FindByOrganizationID(ctx, orgID)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name, description string, status project.Status) (*project.Project, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}

	if name != "" {
		p.Name = name
	}
	if description != "" {
		p.Description = &description
	}
	if status != "" && status.IsValid() {
		p.Status = status
	}

	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return p, nil
}

func (s *Service) Archive(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}
	p.Archive()
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("archive: %w", err)
	}
	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return ErrProjectNotFound
	}
	return s.repo.SoftDelete(ctx, id)
}
