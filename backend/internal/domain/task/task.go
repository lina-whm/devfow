package task

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Type string

const (
	TypeTask  Type = "task"
	TypeBug   Type = "bug"
	TypeStory Type = "story"
	TypeEpic  Type = "epic"
)

func (t Type) IsValid() bool {
	switch t {
	case TypeTask, TypeBug, TypeStory, TypeEpic:
		return true
	}
	return false
}

type Priority string

const (
	PriorityNone     Priority = "none"
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "urgent"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityNone, PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	}
	return false
}

type Status string

const (
	StatusBacklog    Status = "backlog"
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusInReview   Status = "in_review"
	StatusDone       Status = "done"
	StatusCancelled  Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusBacklog, StatusTodo, StatusInProgress, StatusInReview, StatusDone, StatusCancelled:
		return true
	}
	return false
}

var validTransitions = map[Status][]Status{
	StatusBacklog:    {StatusTodo, StatusCancelled},
	StatusTodo:       {StatusInProgress, StatusCancelled},
	StatusInProgress: {StatusInReview, StatusTodo, StatusCancelled},
	StatusInReview:   {StatusDone, StatusInProgress, StatusCancelled},
	StatusDone:       {},
	StatusCancelled:  {StatusBacklog},
}

type Task struct {
	ID           uuid.UUID  `json:"id"`
	ProjectID    uuid.UUID  `json:"project_id"`
	ColumnID     *uuid.UUID `json:"column_id,omitempty"`
	SprintID     *uuid.UUID `json:"sprint_id,omitempty"`
	ParentTaskID *uuid.UUID `json:"parent_task_id,omitempty"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Type         Type       `json:"type"`
	Priority     Priority   `json:"priority"`
	Status       Status     `json:"status"`
	Position     float64    `json:"position"`
	AssigneeID   *uuid.UUID `json:"assignee_id,omitempty"`
	ReporterID   uuid.UUID  `json:"reporter_id"`
	VersionID    int        `json:"version_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewTask(projectID, reporterID uuid.UUID, title, description string, taskType Type) *Task {
	return &Task{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Type:        taskType,
		Priority:    PriorityNone,
		Status:      StatusBacklog,
		ReporterID:  reporterID,
		VersionID:   1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (t *Task) Validate() error {
	if t.ProjectID == uuid.Nil {
		return errors.New("project is required")
	}
	if t.Title == "" {
		return errors.New("title is required")
	}
	titleLen := utf8.RuneCountInString(t.Title)
	if titleLen < 1 || titleLen > 500 {
		return errors.New("title must be between 1 and 500 characters")
	}
	if !t.Type.IsValid() {
		return errors.New("invalid task type")
	}
	if !t.Priority.IsValid() {
		return errors.New("invalid task priority")
	}
	if !t.Status.IsValid() {
		return errors.New("invalid task status")
	}
	if t.ReporterID == uuid.Nil {
		return errors.New("reporter is required")
	}
	return nil
}

func (t *Task) ChangeStatus(newStatus Status) error {
	if !newStatus.IsValid() {
		return errors.New("invalid status")
	}
	if err := t.ValidateTransition(newStatus); err != nil {
		return err
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) ValidateTransition(newStatus Status) error {
	if newStatus == t.Status {
		return nil
	}
	allowed, ok := validTransitions[t.Status]
	if !ok {
		return errors.New("invalid current status")
	}
	for _, s := range allowed {
		if s == newStatus {
			return nil
		}
	}
	return errors.New("cannot transition from " + string(t.Status) + " to " + string(newStatus))
}

func (t *Task) MoveToColumn(columnID uuid.UUID) {
	t.ColumnID = &columnID
	t.UpdatedAt = time.Now()
}

func (t *Task) UpdatePriority(priority Priority) error {
	if !priority.IsValid() {
		return errors.New("invalid priority")
	}
	t.Priority = priority
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) Assign(userID uuid.UUID) {
	t.AssigneeID = &userID
	t.UpdatedAt = time.Now()
}

func (t *Task) Unassign() {
	t.AssigneeID = nil
	t.UpdatedAt = time.Now()
}

func (t *Task) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t *Task) IsDone() bool {
	return t.Status == StatusDone
}
