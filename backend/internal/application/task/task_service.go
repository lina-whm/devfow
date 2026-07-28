package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/domain/task"
	"github.com/devflow/devflow-backend/internal/domain/user"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrVersionConflict   = errors.New("task was modified by another operation")
)

type Service struct {
	taskRepo task.Repository
	userRepo user.Repository
}

func NewService(taskRepo task.Repository, userRepo user.Repository) *Service {
	return &Service{taskRepo: taskRepo, userRepo: userRepo}
}

func (s *Service) Create(ctx context.Context, projectID, reporterID uuid.UUID, title, description string, taskType task.Type, assigneeID *uuid.UUID) (*task.Task, error) {
	if assigneeID != nil && *assigneeID == uuid.Nil {
		assigneeID = nil
	}

	t := task.NewTask(projectID, reporterID, title, description, taskType)
	if assigneeID != nil {
		t.AssigneeID = assigneeID
	}

	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if err := s.taskRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return t, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return t, nil
}

func (s *Service) List(ctx context.Context, projectID uuid.UUID, filter task.TaskFilter) ([]task.Task, *time.Time, error) {
	filter.ProjectID = &projectID
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	return s.taskRepo.List(ctx, filter)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, title, description string, taskType task.Type, priority task.Priority, status task.Status, assigneeID *uuid.UUID, versionID int) (*task.Task, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	if taskType != "" && taskType.IsValid() {
		t.Type = taskType
	}
	if priority != "" && priority.IsValid() {
		t.Priority = priority
	}
	if status != "" && status.IsValid() {
		if err := t.ChangeStatus(status); err != nil {
			return nil, ErrInvalidTransition
		}
	}
	if assigneeID != nil {
		if *assigneeID == uuid.Nil {
			t.Unassign()
		} else {
			t.Assign(*assigneeID)
		}
	}

	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if versionID > 0 {
		if err := s.taskRepo.UpdateWithVersion(ctx, t, versionID); err != nil {
			return nil, ErrVersionConflict
		}
	} else {
		if err := s.taskRepo.Update(ctx, t); err != nil {
			return nil, fmt.Errorf("update: %w", err)
		}
	}
	return t, nil
}

func (s *Service) Move(ctx context.Context, id uuid.UUID, columnID uuid.UUID, position float64, versionID int) (*task.Task, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	t.MoveToColumn(columnID)
	t.Position = position

	if versionID > 0 {
		if err := s.taskRepo.UpdateWithVersion(ctx, t, versionID); err != nil {
			return nil, ErrVersionConflict
		}
	} else {
		if err := s.taskRepo.Update(ctx, t); err != nil {
			return nil, fmt.Errorf("move: %w", err)
		}
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.taskRepo.FindByID(ctx, id); err != nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.SoftDelete(ctx, id)
}

func (s *Service) Restore(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	t, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	if err := s.taskRepo.Restore(ctx, id); err != nil {
		return nil, fmt.Errorf("restore: %w", err)
	}
	return t, nil
}
