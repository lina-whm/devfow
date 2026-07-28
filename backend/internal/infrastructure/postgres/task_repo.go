package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devflow/devflow-backend/internal/domain/task"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tasks (id, project_id, column_id, sprint_id, parent_task_id, title, description, type, priority, status, position, assignee_id, reporter_id, version_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		t.ID, t.ProjectID, t.ColumnID, t.SprintID, t.ParentTaskID, t.Title, t.Description,
		t.Type, t.Priority, t.Status, t.Position, t.AssigneeID, t.ReporterID, t.VersionID,
		t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	t := &task.Task{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, column_id, sprint_id, parent_task_id, title, description, type, priority, status, position, assignee_id, reporter_id, version_id, created_at, updated_at, deleted_at
		 FROM tasks WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&t.ID, &t.ProjectID, &t.ColumnID, &t.SprintID, &t.ParentTaskID, &t.Title, &t.Description,
		&t.Type, &t.Priority, &t.Status, &t.Position, &t.AssigneeID, &t.ReporterID, &t.VersionID,
		&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("find task: %w", err)
	}
	return t, nil
}

func (r *TaskRepository) List(ctx context.Context, filter task.TaskFilter) ([]task.Task, *time.Time, error) {
	query := `SELECT id, project_id, column_id, sprint_id, parent_task_id, title, description, type, priority, status, position, assignee_id, reporter_id, version_id, created_at, updated_at, deleted_at
	          FROM tasks WHERE deleted_at IS NULL`
	args := make([]interface{}, 0)
	argIdx := 1

	if filter.ProjectID != nil {
		query += fmt.Sprintf(" AND project_id = $%d", argIdx)
		args = append(args, *filter.ProjectID)
		argIdx++
	}
	if filter.SprintID != nil {
		query += fmt.Sprintf(" AND sprint_id = $%d", argIdx)
		args = append(args, *filter.SprintID)
		argIdx++
	}
	if filter.ColumnID != nil {
		query += fmt.Sprintf(" AND column_id = $%d", argIdx)
		args = append(args, *filter.ColumnID)
		argIdx++
	}
	if filter.AssigneeID != nil {
		query += fmt.Sprintf(" AND assignee_id = $%d", argIdx)
		args = append(args, *filter.AssigneeID)
		argIdx++
	}
	if filter.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, *filter.Type)
		argIdx++
	}
	if filter.Priority != nil {
		query += fmt.Sprintf(" AND priority = $%d", argIdx)
		args = append(args, *filter.Priority)
		argIdx++
	}
	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	if filter.Cursor != nil {
		query += fmt.Sprintf(" AND created_at < $%d", argIdx)
		args = append(args, *filter.Cursor)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.ColumnID, &t.SprintID, &t.ParentTaskID, &t.Title, &t.Description,
			&t.Type, &t.Priority, &t.Status, &t.Position, &t.AssigneeID, &t.ReporterID, &t.VersionID,
			&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt); err != nil {
			return nil, nil, err
		}
		tasks = append(tasks, t)
	}

	var nextCursor *time.Time
	if len(tasks) > 0 {
		nextCursor = &tasks[len(tasks)-1].CreatedAt
	}

	return tasks, nextCursor, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tasks SET title=$1, description=$2, type=$3, priority=$4, status=$5, position=$6, column_id=$7, assignee_id=$8, updated_at=NOW() WHERE id=$9 AND deleted_at IS NULL`,
		t.Title, t.Description, t.Type, t.Priority, t.Status, t.Position, t.ColumnID, t.AssigneeID, t.ID)
	return err
}

func (r *TaskRepository) UpdateWithVersion(ctx context.Context, t *task.Task, currentVersion int) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE tasks SET title=$1, description=$2, type=$3, priority=$4, status=$5, position=$6, column_id=$7, assignee_id=$8, version_id=version_id+1, updated_at=NOW()
		 WHERE id=$9 AND version_id=$10 AND deleted_at IS NULL`,
		t.Title, t.Description, t.Type, t.Priority, t.Status, t.Position, t.ColumnID, t.AssigneeID, t.ID, currentVersion)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task was modified by another operation")
	}
	return nil
}

func (r *TaskRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE tasks SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *TaskRepository) Restore(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE tasks SET deleted_at=NULL WHERE id=$1`, id)
	return err
}
