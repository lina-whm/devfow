package response

type TaskResponse struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	ColumnID    string    `json:"column_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	Position    float64   `json:"position"`
	Assignee    *TaskUser `json:"assignee,omitempty"`
	Tags        []TaskTag `json:"tags,omitempty"`
	VersionID   int       `json:"version_id"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type TaskUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

type TaskTag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TaskListResponse struct {
	Tasks      []TaskResponse `json:"data"`
	Total      int            `json:"total"`
	NextCursor string         `json:"next_cursor,omitempty"`
	HasMore    bool           `json:"has_more"`
}

type TaskPositionResponse struct {
	TaskID    string  `json:"task_id"`
	ColumnID  string  `json:"column_id"`
	Position  float64 `json:"position"`
	VersionID string  `json:"version_id"`
}
