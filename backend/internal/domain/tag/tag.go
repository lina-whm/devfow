package tag

import (
	"errors"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var hexColorRegex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

type Tag struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTag(name, color string) *Tag {
	return &Tag{
		ID:        uuid.New(),
		Name:      name,
		Color:     color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *Tag) Validate() error {
	if t.Name == "" {
		return errors.New("name is required")
	}
	nameLen := utf8.RuneCountInString(t.Name)
	if nameLen < 1 || nameLen > 30 {
		return errors.New("name must be between 1 and 30 characters")
	}
	if !hexColorRegex.MatchString(t.Color) {
		return errors.New("color must be a valid hex color (e.g., #FF0000)")
	}
	return nil
}
