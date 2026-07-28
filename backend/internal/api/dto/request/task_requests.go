// Package request contains HTTP request DTO types.
package request

type CreateTaskRequest struct {
	ProjectID   string   `json:"-"`
	Title       string   `json:"title" binding:"required,min=1,max=255"`
	Description string   `json:"description"`
	Type        string   `json:"type" binding:"omitempty,oneof=task bug story epic"`
	Priority    string   `json:"priority" binding:"omitempty,oneof=none low medium high critical"`
	AssigneeID  string   `json:"assignee_id"`
	TagIDs      []string `json:"tag_ids"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title" binding:"omitempty,min=1,max=255"`
	Description string  `json:"description"`
	Type        string  `json:"type" binding:"omitempty,oneof=task bug story epic"`
	Priority    string  `json:"priority" binding:"omitempty,oneof=none low medium high critical"`
	Status      string  `json:"status" binding:"omitempty,oneof=backlog todo in_progress in_review done cancelled"`
	AssigneeID  *string `json:"assignee_id"`
	VersionID   int     `json:"version_id"`
}

type MoveTaskRequest struct {
	ColumnID  string  `json:"column_id" binding:"required"`
	Position  float64 `json:"position"`
	VersionID int     `json:"version_id"`
}

type ListTasksQuery struct {
	Statuses   string `form:"status"`
	Priorities string `form:"priority"`
	AssigneeID string `form:"assignee_id"`
	TagIDs     string `form:"tag_ids"`
	Query      string `form:"q"`
	Cursor     string `form:"cursor"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Sort       string `form:"sort"`
}
