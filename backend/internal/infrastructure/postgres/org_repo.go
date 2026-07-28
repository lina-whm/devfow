package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/devflow/devflow-backend/internal/domain/organization"
)

type OrgRepository struct {
	pool *pgxpool.Pool
}

func NewOrgRepository(pool *pgxpool.Pool) *OrgRepository {
	return &OrgRepository{pool: pool}
}

func (r *OrgRepository) Create(ctx context.Context, org *organization.Organization) error {
	query := `INSERT INTO organizations (id, name, slug, description, owner_id, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, org.ID, org.Name, org.Slug, org.Description, org.OwnerID, org.CreatedAt, org.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	return nil
}

func (r *OrgRepository) FindByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error) {
	query := `SELECT id, name, slug, description, owner_id, created_at, updated_at, deleted_at
	          FROM organizations WHERE id = $1 AND deleted_at IS NULL`
	org := &organization.Organization{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.Slug, &org.Description, &org.OwnerID,
		&org.CreatedAt, &org.UpdatedAt, &org.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("org not found")
		}
		return nil, fmt.Errorf("find org: %w", err)
	}
	return org, nil
}

func (r *OrgRepository) FindBySlug(ctx context.Context, slug string) (*organization.Organization, error) {
	query := `SELECT id, name, slug, description, owner_id, created_at, updated_at, deleted_at
	          FROM organizations WHERE slug = $1 AND deleted_at IS NULL`
	org := &organization.Organization{}
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&org.ID, &org.Name, &org.Slug, &org.Description, &org.OwnerID,
		&org.CreatedAt, &org.UpdatedAt, &org.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("org not found")
		}
		return nil, fmt.Errorf("find by slug: %w", err)
	}
	return org, nil
}

func (r *OrgRepository) Update(ctx context.Context, org *organization.Organization) error {
	query := `UPDATE organizations SET name=$1, slug=$2, description=$3, updated_at=NOW() WHERE id=$4 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, org.Name, org.Slug, org.Description, org.ID)
	return err
}

func (r *OrgRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE organizations SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

type OrgMemberRepository struct {
	pool *pgxpool.Pool
}

func NewOrgMemberRepository(pool *pgxpool.Pool) *OrgMemberRepository {
	return &OrgMemberRepository{pool: pool}
}

func (r *OrgMemberRepository) AddMember(ctx context.Context, member *organization.OrganizationMember) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO organization_members (organization_id, user_id, role) VALUES ($1, $2, $3)`,
		member.OrganizationID, member.UserID, member.Role)
	return err
}

func (r *OrgMemberRepository) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]organization.OrganizationMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT organization_id, user_id, role, joined_at FROM organization_members WHERE organization_id = $1`,
		orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []organization.OrganizationMember
	for rows.Next() {
		var m organization.OrganizationMember
		if err := rows.Scan(&m.OrganizationID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *OrgMemberRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]organization.OrganizationMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT organization_id, user_id, role, joined_at FROM organization_members WHERE user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []organization.OrganizationMember
	for rows.Next() {
		var m organization.OrganizationMember
		if err := rows.Scan(&m.OrganizationID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *OrgMemberRepository) FindMember(ctx context.Context, orgID, userID uuid.UUID) (*organization.OrganizationMember, error) {
	var m organization.OrganizationMember
	err := r.pool.QueryRow(ctx,
		`SELECT organization_id, user_id, role, joined_at FROM organization_members WHERE organization_id = $1 AND user_id = $2`,
		orgID, userID).Scan(&m.OrganizationID, &m.UserID, &m.Role, &m.JoinedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("member not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *OrgMemberRepository) UpdateRole(ctx context.Context, orgID, userID uuid.UUID, role organization.Role) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE organization_members SET role = $1 WHERE organization_id = $2 AND user_id = $3`,
		role, orgID, userID)
	return err
}

func (r *OrgMemberRepository) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM organization_members WHERE organization_id = $1 AND user_id = $2`,
		orgID, userID)
	return err
}

type InvitationRepository struct {
	pool *pgxpool.Pool
}

func NewInvitationRepository(pool *pgxpool.Pool) *InvitationRepository {
	return &InvitationRepository{pool: pool}
}

func (r *InvitationRepository) Create(ctx context.Context, inv *organization.Invitation) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO invitations (id, organization_id, inviter_id, invitee_email, token_hash, role, status, expires_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		inv.ID, inv.OrganizationID, inv.InviterID, inv.InviteeEmail, inv.TokenHash,
		inv.Role, inv.Status, inv.ExpiresAt, inv.CreatedAt, inv.UpdatedAt)
	return err
}

func (r *InvitationRepository) FindByID(ctx context.Context, id uuid.UUID) (*organization.Invitation, error) {
	var inv organization.Invitation
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, inviter_id, invitee_email, token_hash, role, status, expires_at, accepted_at, created_at, updated_at
		 FROM invitations WHERE id = $1`, id).Scan(
		&inv.ID, &inv.OrganizationID, &inv.InviterID, &inv.InviteeEmail, &inv.TokenHash,
		&inv.Role, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	return &inv, nil
}

func (r *InvitationRepository) FindByOrgID(ctx context.Context, orgID uuid.UUID) ([]organization.Invitation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, organization_id, inviter_id, invitee_email, token_hash, role, status, expires_at, accepted_at, created_at, updated_at
		 FROM invitations WHERE organization_id = $1`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var invs []organization.Invitation
	for rows.Next() {
		var inv organization.Invitation
		if err := rows.Scan(&inv.ID, &inv.OrganizationID, &inv.InviterID, &inv.InviteeEmail, &inv.TokenHash,
			&inv.Role, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, err
		}
		invs = append(invs, inv)
	}
	return invs, nil
}

func (r *InvitationRepository) FindByEmail(ctx context.Context, email string) ([]organization.Invitation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, organization_id, inviter_id, invitee_email, token_hash, role, status, expires_at, accepted_at, created_at, updated_at
		 FROM invitations WHERE invitee_email = $1`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var invs []organization.Invitation
	for rows.Next() {
		var inv organization.Invitation
		if err := rows.Scan(&inv.ID, &inv.OrganizationID, &inv.InviterID, &inv.InviteeEmail, &inv.TokenHash,
			&inv.Role, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, err
		}
		invs = append(invs, inv)
	}
	return invs, nil
}

func (r *InvitationRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*organization.Invitation, error) {
	var inv organization.Invitation
	err := r.pool.QueryRow(ctx,
		`SELECT id, organization_id, inviter_id, invitee_email, token_hash, role, status, expires_at, accepted_at, created_at, updated_at
		 FROM invitations WHERE token_hash = $1`, tokenHash).Scan(
		&inv.ID, &inv.OrganizationID, &inv.InviterID, &inv.InviteeEmail, &inv.TokenHash,
		&inv.Role, &inv.Status, &inv.ExpiresAt, &inv.AcceptedAt, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	return &inv, nil
}

func (r *InvitationRepository) Update(ctx context.Context, inv *organization.Invitation) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE invitations SET status=$1, accepted_at=$2, updated_at=NOW() WHERE id=$3`,
		inv.Status, inv.AcceptedAt, inv.ID)
	return err
}

func (r *InvitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM invitations WHERE id = $1`, id)
	return err
}
