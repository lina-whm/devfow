package organization

import (
	"errors"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Organization struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description,omitempty"`
	OwnerID     uuid.UUID  `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewOrganization(name, slug, description string, ownerID uuid.UUID) *Organization {
	var desc *string
	if description != "" {
		desc = &description
	}
	return &Organization{
		ID:          uuid.New(),
		Name:        name,
		Slug:        slug,
		Description: desc,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (o *Organization) Validate() error {
	if o.Name == "" {
		return errors.New("name is required")
	}
	nameLen := utf8.RuneCountInString(o.Name)
	if nameLen < 2 || nameLen > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	if o.Slug == "" {
		return errors.New("slug is required")
	}
	if !slugRegex.MatchString(o.Slug) {
		return errors.New("slug must contain only lowercase letters, numbers, and hyphens")
	}
	slugLen := utf8.RuneCountInString(o.Slug)
	if slugLen < 2 || slugLen > 80 {
		return errors.New("slug must be between 2 and 80 characters")
	}
	if o.OwnerID == uuid.Nil {
		return errors.New("owner is required")
	}
	return nil
}

func (o *Organization) IsDeleted() bool {
	return o.DeletedAt != nil
}
