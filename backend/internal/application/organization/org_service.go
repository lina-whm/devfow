package organization

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/domain/organization"
	"github.com/devflow/devflow-backend/internal/domain/user"
)

var (
	ErrOrgNotFound     = errors.New("organization not found")
	ErrOrgSlugTaken    = errors.New("organization slug already taken")
	ErrMemberNotFound  = errors.New("member not found")
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrInvitationExpired  = errors.New("invitation has expired")
	ErrInsufficientPerms  = errors.New("insufficient permissions")
)

type Service struct {
	orgRepo  organization.Repository
	memRepo  organization.MemberRepository
	invRepo  organization.InvitationRepository
	userRepo user.Repository
}

func NewService(
	orgRepo organization.Repository,
	memRepo organization.MemberRepository,
	invRepo organization.InvitationRepository,
	userRepo user.Repository,
) *Service {
	return &Service{
		orgRepo:  orgRepo,
		memRepo:  memRepo,
		invRepo:  invRepo,
		userRepo: userRepo,
	}
}

func (s *Service) Create(ctx context.Context, name, slug string, ownerID uuid.UUID) (*organization.Organization, error) {
	existing, _ := s.orgRepo.FindBySlug(ctx, slug)
	if existing != nil {
		return nil, ErrOrgSlugTaken
	}

	org := organization.NewOrganization(name, slug, "", ownerID)
	if err := org.Validate(); err != nil {
		return nil, fmt.Errorf("validate org: %w", err)
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("create org: %w", err)
	}

	member, err := organization.NewOrganizationMember(org.ID, ownerID, organization.RoleOwner)
	if err != nil {
		return nil, fmt.Errorf("create owner member: %w", err)
	}
	if err := s.memRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("add owner: %w", err)
	}

	return org, nil
}

func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID) (*organization.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, ErrOrgNotFound
	}
	return org, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID uuid.UUID) ([]organization.Organization, error) {
	members, err := s.memRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find memberships: %w", err)
	}

	orgs := make([]organization.Organization, 0, len(members))
	for _, m := range members {
		org, err := s.orgRepo.FindByID(ctx, m.OrganizationID)
		if err != nil {
			continue
		}
		orgs = append(orgs, *org)
	}
	return orgs, nil
}

func (s *Service) Update(ctx context.Context, orgID uuid.UUID, name, slug string) (*organization.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, ErrOrgNotFound
	}

	if name != "" {
		org.Name = name
	}
	if slug != "" {
		existing, _ := s.orgRepo.FindBySlug(ctx, slug)
		if existing != nil && existing.ID != orgID {
			return nil, ErrOrgSlugTaken
		}
		org.Slug = slug
	}

	if err := org.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}
	return org, nil
}

func (s *Service) Delete(ctx context.Context, orgID uuid.UUID) error {
	if _, err := s.orgRepo.FindByID(ctx, orgID); err != nil {
		return ErrOrgNotFound
	}
	return s.orgRepo.SoftDelete(ctx, orgID)
}

func (s *Service) InviteMember(ctx context.Context, orgID, inviterID uuid.UUID, email string, role organization.Role) (*organization.Invitation, error) {
	member, err := s.memRepo.FindMember(ctx, orgID, inviterID)
	if err != nil || !member.Role.CanManageMembers() {
		return nil, ErrInsufficientPerms
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	rawToken := hex.EncodeToString(tokenBytes)
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	inv, err := organization.NewInvitation(orgID, inviterID, email, tokenHash, role, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return nil, err
	}

	// We store the raw token so we can return it once
	if err := s.invRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("create invitation: %w", err)
	}

	// The returned invitation's token is NOT stored, so the caller gets it once
	return inv, nil
}

func (s *Service) ListMembers(ctx context.Context, orgID uuid.UUID) ([]organization.OrganizationMember, error) {
	return s.memRepo.FindByOrgID(ctx, orgID)
}

func (s *Service) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	member, err := s.memRepo.FindMember(ctx, orgID, userID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.Role == organization.RoleOwner {
		return errors.New("cannot remove the organization owner")
	}
	return s.memRepo.RemoveMember(ctx, orgID, userID)
}

func (s *Service) ChangeMemberRole(ctx context.Context, orgID, userID uuid.UUID, newRole organization.Role) error {
	member, err := s.memRepo.FindMember(ctx, orgID, userID)
	if err != nil {
		return ErrMemberNotFound
	}
	if member.Role == organization.RoleOwner {
		return errors.New("cannot change the owner's role")
	}
	return member.ChangeRole(newRole)
}
