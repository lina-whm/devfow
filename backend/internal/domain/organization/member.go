package organization

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember:
		return true
	}
	return false
}

func (r Role) CanManageMembers() bool {
	return r == RoleOwner || r == RoleAdmin
}

func (r Role) CanDeleteOrganization() bool {
	return r == RoleOwner
}

type OrganizationMember struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
	Role           Role      `json:"role"`
	JoinedAt       time.Time `json:"joined_at"`
}

func NewOrganizationMember(organizationID, userID uuid.UUID, role Role) (*OrganizationMember, error) {
	if !role.IsValid() {
		return nil, errors.New("invalid role")
	}
	return &OrganizationMember{
		OrganizationID: organizationID,
		UserID:         userID,
		Role:           role,
		JoinedAt:       time.Now(),
	}, nil
}

func (m *OrganizationMember) Validate() error {
	if m.OrganizationID == uuid.Nil {
		return errors.New("organization is required")
	}
	if m.UserID == uuid.Nil {
		return errors.New("user is required")
	}
	if !m.Role.IsValid() {
		return errors.New("invalid role")
	}
	return nil
}

func (m *OrganizationMember) ChangeRole(newRole Role) error {
	if !newRole.IsValid() {
		return errors.New("invalid role")
	}
	m.Role = newRole
	return nil
}
