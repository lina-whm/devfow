package board

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	domain_board "github.com/devflow/devflow-backend/internal/domain/board"
	domain_task "github.com/devflow/devflow-backend/internal/domain/task"
)

var (
	ErrBoardNotFound = errors.New("board not found")
)

type BoardWithColumns struct {
	Board   *domain_board.Board
	Columns []ColumnWithTasks
}

type ColumnWithTasks struct {
	Column domain_board.Column
	Tasks  []domain_task.Task
}

type Service struct {
	boardRepo  domain_board.Repository
	columnRepo domain_board.ColumnRepository
	taskRepo   domain_task.Repository
}

func NewService(boardRepo domain_board.Repository, columnRepo domain_board.ColumnRepository, taskRepo domain_task.Repository) *Service {
	return &Service{boardRepo: boardRepo, columnRepo: columnRepo, taskRepo: taskRepo}
}

func (s *Service) GetOrCreate(ctx context.Context, projectID uuid.UUID) (*BoardWithColumns, error) {
	b, err := s.boardRepo.FindByProjectID(ctx, projectID)
	if err != nil || b == nil {
		b = domain_board.NewBoard(projectID, "Kanban")
		if err := s.boardRepo.Create(ctx, b); err != nil {
			return nil, fmt.Errorf("create board: %w", err)
		}

		defaultColumns := []string{"Backlog", "To Do", "In Progress", "Done"}
		for i, name := range defaultColumns {
			col := domain_board.NewColumn(b.ID, name, float64(i+1))
			if err := s.columnRepo.Create(ctx, col); err != nil {
				return nil, fmt.Errorf("create column %q: %w", name, err)
			}
		}
	}

	columns, err := s.columnRepo.FindByBoardID(ctx, b.ID)
	if err != nil {
		return nil, fmt.Errorf("find columns: %w", err)
	}

	result := BoardWithColumns{Board: b}
	for _, col := range columns {
		tasks, _, err := s.taskRepo.List(ctx, domain_task.TaskFilter{ColumnID: &col.ID, Limit: 100})
		if err != nil {
			return nil, fmt.Errorf("list tasks for column %q: %w", col.Name, err)
		}
		if tasks == nil {
			tasks = []domain_task.Task{}
		}
		result.Columns = append(result.Columns, ColumnWithTasks{
			Column: col,
			Tasks:  tasks,
		})
	}

	return &result, nil
}

type ColumnUpdate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Service) UpdateColumns(ctx context.Context, boardID uuid.UUID, updates []ColumnUpdate) error {
	for _, u := range updates {
		colID, err := uuid.Parse(u.ID)
		if err != nil {
			return fmt.Errorf("parse column id: %w", err)
		}
		col, err := s.columnRepo.FindByID(ctx, colID)
		if err != nil {
			return fmt.Errorf("find column %q: %w", u.ID, err)
		}
		if u.Name != "" {
			col.Name = u.Name
		}
		if err := s.columnRepo.Update(ctx, col); err != nil {
			return fmt.Errorf("update column %q: %w", u.ID, err)
		}
	}
	return nil
}
