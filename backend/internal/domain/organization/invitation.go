package organization

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type InvitationStatus string

const (
	InvitationStatusPending   InvitationStatus = "pending"
	InvitationStatusAccepted  InvitationStatus = "accepted"
	InvitationStatusExpired   InvitationStatus = "expired"
	InvitationStatusCancelled InvitationStatus = "cancelled"
)

type Invitation struct {
	ID             uuid.UUID        `json:"id"`
	OrganizationID uuid.UUID        `json:"organization_id"`
	InviterID      uuid.UUID        `json:"inviter_id"`
	InviteeEmail   string           `json:"invitee_email"`
	TokenHash      string           `json:"token_hash"`
	Role           Role             `json:"role"`
	Status         InvitationStatus `json:"status"`
	ExpiresAt      time.Time        `json:"expires_at"`
	AcceptedAt     *time.Time       `json:"accepted_at,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

func NewInvitation(organizationID, inviterID uuid.UUID, inviteeEmail, tokenHash string, role Role, expiresAt time.Time) (*Invitation, error) {
	if !role.IsValid() {
		return nil, errors.New("invalid role")
	}
	return &Invitation{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		InviterID:      inviterID,
		InviteeEmail:   inviteeEmail,
		TokenHash:      tokenHash,
		Role:           role,
		Status:         InvitationStatusPending,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (i *Invitation) Validate() error {
	if i.OrganizationID == uuid.Nil {
		return errors.New("organization is required")
	}
	if i.InviterID == uuid.Nil {
		return errors.New("inviter is required")
	}
	if i.InviteeEmail == "" {
		return errors.New("invitee email is required")
	}
	if i.TokenHash == "" {
		return errors.New("token hash is required")
	}
	if !i.Role.IsValid() {
		return errors.New("invalid role")
	}
	if i.ExpiresAt.Before(time.Now()) {
		return errors.New("invitation has already expired")
	}
	return nil
}

func (i *Invitation) Accept() error {
	if i.Status != InvitationStatusPending {
		return errors.New("only pending invitations can be accepted")
	}
	if time.Now().After(i.ExpiresAt) {
		i.Status = InvitationStatusExpired
		return errors.New("invitation has expired")
	}
	now := time.Now()
	i.Status = InvitationStatusAccepted
	i.AcceptedAt = &now
	i.UpdatedAt = now
	return nil
}

func (i *Invitation) Cancel() error {
	if i.Status != InvitationStatusPending {
		return errors.New("only pending invitations can be cancelled")
	}
	i.Status = InvitationStatusCancelled
	i.UpdatedAt = time.Now()
	return nil
}

func (i *Invitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

func (i *Invitation) IsActive() bool {
	return i.Status == InvitationStatusPending && !i.IsExpired()
}
