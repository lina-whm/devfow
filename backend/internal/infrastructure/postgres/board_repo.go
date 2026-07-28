package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devflow/devflow-backend/internal/domain/board"
)

type BoardRepository struct {
	pool *pgxpool.Pool
}

func NewBoardRepository(pool *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{pool: pool}
}

func (r *BoardRepository) Create(ctx context.Context, b *board.Board) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO boards (id, project_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		b.ID, b.ProjectID, b.Name, b.CreatedAt, b.UpdatedAt)
	return err
}

func (r *BoardRepository) FindByID(ctx context.Context, id uuid.UUID) (*board.Board, error) {
	b := &board.Board{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, created_at, updated_at FROM boards WHERE id = $1`, id).
		Scan(&b.ID, &b.ProjectID, &b.Name, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("board not found")
		}
		return nil, err
	}
	return b, nil
}

func (r *BoardRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) (*board.Board, error) {
	b := &board.Board{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, created_at, updated_at FROM boards WHERE project_id = $1`, projectID).
		Scan(&b.ID, &b.ProjectID, &b.Name, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("board not found")
		}
		return nil, err
	}
	return b, nil
}

func (r *BoardRepository) Update(ctx context.Context, b *board.Board) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE boards SET name=$1, updated_at=NOW() WHERE id=$2`, b.Name, b.ID)
	return err
}

func (r *BoardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM boards WHERE id = $1`, id)
	return err
}

type ColumnRepository struct {
	pool *pgxpool.Pool
}

func NewColumnRepository(pool *pgxpool.Pool) *ColumnRepository {
	return &ColumnRepository{pool: pool}
}

func (r *ColumnRepository) Create(ctx context.Context, c *board.Column) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO columns (id, board_id, name, position, wip_limit, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.ID, c.BoardID, c.Name, c.Position, c.WIPLimit, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *ColumnRepository) FindByID(ctx context.Context, id uuid.UUID) (*board.Column, error) {
	c := &board.Column{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, board_id, name, position, wip_limit, created_at, updated_at FROM columns WHERE id = $1`, id).
		Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &c.WIPLimit, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("column not found")
		}
		return nil, err
	}
	return c, nil
}

func (r *ColumnRepository) FindByBoardID(ctx context.Context, boardID uuid.UUID) ([]board.Column, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, board_id, name, position, wip_limit, created_at, updated_at FROM columns WHERE board_id = $1 ORDER BY position`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []board.Column
	for rows.Next() {
		var c board.Column
		if err := rows.Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &c.WIPLimit, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	return cols, nil
}

func (r *ColumnRepository) Update(ctx context.Context, c *board.Column) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE columns SET name=$1, position=$2, wip_limit=$3, updated_at=NOW() WHERE id=$4`,
		c.Name, c.Position, c.WIPLimit, c.ID)
	return err
}

func (r *ColumnRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM columns WHERE id = $1`, id)
	return err
}

func (r *ColumnRepository) Reorder(ctx context.Context, boardID uuid.UUID, columnIDs []uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, colID := range columnIDs {
		if _, err := tx.Exec(ctx,
			`UPDATE columns SET position = $1, updated_at = NOW() WHERE id = $2 AND board_id = $3`,
			i, colID, boardID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
