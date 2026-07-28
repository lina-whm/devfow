package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/domain/team"
)

var (
	ErrTeamNotFound = errors.New("team not found")
)

type Service struct {
	teamRepo team.Repository
	memRepo  team.MemberRepository
}

func NewService(teamRepo team.Repository, memRepo team.MemberRepository) *Service {
	return &Service{teamRepo: teamRepo, memRepo: memRepo}
}

func (s *Service) Create(ctx context.Context, name string, orgID uuid.UUID) (*team.Team, error) {
	t := team.NewTeam(name, orgID)
	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	if err := s.teamRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return t, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*team.Team, error) {
	t, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return t, nil
}

func (s *Service) ListByOrgID(ctx context.Context, orgID uuid.UUID) ([]team.Team, error) {
	return s.teamRepo.FindByOrganizationID(ctx, orgID)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, name string) (*team.Team, error) {
	t, err := s.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	t.Name = name
	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	if err := s.teamRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.teamRepo.FindByID(ctx, id); err != nil {
		return ErrTeamNotFound
	}
	return s.teamRepo.SoftDelete(ctx, id)
}

func (s *Service) AddMember(ctx context.Context, teamID, userID uuid.UUID) error {
	member := &team.TeamMember{
		TeamID: teamID,
		UserID: userID,
	}
	return s.memRepo.AddMember(ctx, member)
}

func (s *Service) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	return s.memRepo.RemoveMember(ctx, teamID, userID)
}
