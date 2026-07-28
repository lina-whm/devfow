package notification

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Type string

const (
	TypeTaskAssigned      Type = "task_assigned"
	TypeTaskUpdated       Type = "task_updated"
	TypeCommentAdded      Type = "comment_added"
	TypeInvitationReceived Type = "invitation_received"
	TypeSprintStarted     Type = "sprint_started"
	TypeMentioned         Type = "mentioned"
)

type Notification struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	Type      Type                   `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ReadAt    *time.Time             `json:"read_at,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

func NewNotification(userID uuid.UUID, notifType Type, title, body string, metadata map[string]interface{}) *Notification {
	return &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}
}

func (n *Notification) Validate() error {
	if n.UserID == uuid.Nil {
		return errors.New("user is required")
	}
	if n.Title == "" {
		return errors.New("title is required")
	}
	titleLen := utf8.RuneCountInString(n.Title)
	if titleLen > 200 {
		return errors.New("title must not exceed 200 characters")
	}
	if utf8.RuneCountInString(n.Body) > 5000 {
		return errors.New("body must not exceed 5000 characters")
	}
	return nil
}

func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.ReadAt = &now
}

func (n *Notification) IsRead() bool {
	return n.ReadAt != nil
}
