package comment

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID  `json:"id"`
	TaskID    uuid.UUID  `json:"task_id"`
	AuthorID  uuid.UUID  `json:"author_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewComment(taskID, authorID uuid.UUID, body string) *Comment {
	return &Comment{
		ID:        uuid.New(),
		TaskID:    taskID,
		AuthorID:  authorID,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c *Comment) Validate() error {
	if c.TaskID == uuid.Nil {
		return errors.New("task is required")
	}
	if c.AuthorID == uuid.Nil {
		return errors.New("author is required")
	}
	if c.Body == "" {
		return errors.New("body is required")
	}
	bodyLen := utf8.RuneCountInString(c.Body)
	if bodyLen > 10000 {
		return errors.New("body must not exceed 10000 characters")
	}
	return nil
}

func (c *Comment) Edit(body string) error {
	if body == "" {
		return errors.New("body is required")
	}
	if utf8.RuneCountInString(body) > 10000 {
		return errors.New("body must not exceed 10000 characters")
	}
	c.Body = body
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Comment) IsDeleted() bool {
	return c.DeletedAt != nil
}
