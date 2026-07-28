package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devflow/devflow-backend/internal/domain/comment"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (r *CommentRepository) Create(ctx context.Context, c *comment.Comment) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO comments (id, task_id, author_id, body, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID, c.TaskID, c.AuthorID, c.Body, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *CommentRepository) FindByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	c := &comment.Comment{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, task_id, author_id, body, created_at, updated_at, deleted_at FROM comments WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&c.ID, &c.TaskID, &c.AuthorID, &c.Body, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, err
	}
	return c, nil
}

func (r *CommentRepository) FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]comment.Comment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, task_id, author_id, body, created_at, updated_at, deleted_at FROM comments WHERE task_id = $1 AND deleted_at IS NULL ORDER BY created_at`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []comment.Comment
	for rows.Next() {
		var c comment.Comment
		if err := rows.Scan(&c.ID, &c.TaskID, &c.AuthorID, &c.Body, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepository) Update(ctx context.Context, c *comment.Comment) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE comments SET body=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL`, c.Body, c.ID)
	return err
}

func (r *CommentRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE comments SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}
