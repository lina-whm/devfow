package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devflow/devflow-backend/internal/domain/team"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) Create(ctx context.Context, t *team.Team) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO teams (id, organization_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.OrganizationID, t.Name, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *TeamRepository) FindByID(ctx context.Context, id uuid.UUID) (*team.Team, error) {
	t := &team.Team{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, name, created_at, updated_at, deleted_at FROM teams WHERE id = $1 AND deleted_at IS NULL`,
		id).Scan(&t.ID, &t.OrganizationID, &t.Name, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}
	return t, nil
}

func (r *TeamRepository) FindByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]team.Team, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, organization_id, name, created_at, updated_at, deleted_at FROM teams WHERE organization_id = $1 AND deleted_at IS NULL`,
		orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []team.Team
	for rows.Next() {
		var t team.Team
		if err := rows.Scan(&t.ID, &t.OrganizationID, &t.Name, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (r *TeamRepository) Update(ctx context.Context, t *team.Team) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE teams SET name=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL`,
		t.Name, t.ID)
	return err
}

func (r *TeamRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE teams SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

type TeamMemberRepository struct {
	pool *pgxpool.Pool
}

func NewTeamMemberRepository(pool *pgxpool.Pool) *TeamMemberRepository {
	return &TeamMemberRepository{pool: pool}
}

func (r *TeamMemberRepository) AddMember(ctx context.Context, m *team.TeamMember) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO team_members (team_id, user_id) VALUES ($1, $2)`,
		m.TeamID, m.UserID)
	return err
}

func (r *TeamMemberRepository) FindByTeamID(ctx context.Context, teamID uuid.UUID) ([]team.TeamMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT team_id, user_id FROM team_members WHERE team_id = $1`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []team.TeamMember
	for rows.Next() {
		var m team.TeamMember
		if err := rows.Scan(&m.TeamID, &m.UserID); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *TeamMemberRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]team.TeamMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT team_id, user_id FROM team_members WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []team.TeamMember
	for rows.Next() {
		var m team.TeamMember
		if err := rows.Scan(&m.TeamID, &m.UserID); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *TeamMemberRepository) FindMember(ctx context.Context, teamID, userID uuid.UUID) (*team.TeamMember, error) {
	var m team.TeamMember
	err := r.pool.QueryRow(ctx,
		`SELECT team_id, user_id FROM team_members WHERE team_id = $1 AND user_id = $2`,
		teamID, userID).Scan(&m.TeamID, &m.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("member not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *TeamMemberRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`,
		teamID, userID)
	return err
}
