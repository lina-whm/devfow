package comment

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/domain/comment"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

type Service struct {
	repo comment.Repository
}

func NewService(repo comment.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, taskID, authorID uuid.UUID, body string) (*comment.Comment, error) {
	c := comment.NewComment(taskID, authorID, body)
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return c, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrCommentNotFound
	}
	return c, nil
}

func (s *Service) ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]comment.Comment, error) {
	return s.repo.FindByTaskID(ctx, taskID)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, body string) (*comment.Comment, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrCommentNotFound
	}
	if err := c.Edit(body); err != nil {
		return nil, fmt.Errorf("edit: %w", err)
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return c, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return ErrCommentNotFound
	}
	return s.repo.SoftDelete(ctx, id)
}
