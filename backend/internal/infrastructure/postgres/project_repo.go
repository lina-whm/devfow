package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devflow/devflow-backend/internal/domain/project"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) Create(ctx context.Context, p *project.Project) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO projects (id, organization_id, name, key, description, lead_id, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		p.ID, p.OrganizationID, p.Name, p.Key, p.Description, p.LeadID, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	p := &project.Project{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, name, key, description, lead_id, status, created_at, updated_at, deleted_at
		 FROM projects WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&p.ID, &p.OrganizationID, &p.Name, &p.Key, &p.Description, &p.LeadID,
		&p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("find project: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) FindByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]project.Project, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, organization_id, name, key, description, lead_id, status, created_at, updated_at, deleted_at
		 FROM projects WHERE organization_id = $1 AND deleted_at IS NULL`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []project.Project
	for rows.Next() {
		var p project.Project
		if err := rows.Scan(&p.ID, &p.OrganizationID, &p.Name, &p.Key, &p.Description, &p.LeadID,
			&p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) FindByKey(ctx context.Context, orgID uuid.UUID, key string) (*project.Project, error) {
	p := &project.Project{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, name, key, description, lead_id, status, created_at, updated_at, deleted_at
		 FROM projects WHERE organization_id = $1 AND key = $2 AND deleted_at IS NULL`, orgID, key).Scan(
		&p.ID, &p.OrganizationID, &p.Name, &p.Key, &p.Description, &p.LeadID,
		&p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("find by key: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, p *project.Project) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE projects SET name=$1, description=$2, lead_id=$3, status=$4, updated_at=NOW() WHERE id=$5 AND deleted_at IS NULL`,
		p.Name, p.Description, p.LeadID, p.Status, p.ID)
	return err
}

func (r *ProjectRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE projects SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}
