package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/api/dto/request"
	"github.com/devflow/devflow-backend/internal/api/dto/response"
	application_task "github.com/devflow/devflow-backend/internal/application/task"
	domain_task "github.com/devflow/devflow-backend/internal/domain/task"

	"github.com/gin-gonic/gin"
)

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrTaskForbidden  = errors.New("access to task denied")
)

type TaskHandler struct {
	service *application_task.Service
}

func NewTaskHandler(service *application_task.Service) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) List(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid project id"))
		return
	}

	var q request.ListTasksQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}
	if q.Limit <= 0 {
		q.Limit = 50
	}

	filter := domain_task.TaskFilter{Limit: q.Limit}

	if q.AssigneeID != "" {
		if uid, err := uuid.Parse(q.AssigneeID); err == nil {
			filter.AssigneeID = &uid
		}
	}
	if q.Query != "" {
		filter.Search = q.Query
	}
	if q.Cursor != "" {
		if t, err := time.Parse(time.RFC3339, q.Cursor); err == nil {
			filter.Cursor = &t
		}
	}
	if q.Statuses != "" {
		status := domain_task.Status(strings.Split(q.Statuses, ",")[0])
		filter.Status = &status
	}
	if q.Priorities != "" {
		priority := domain_task.Priority(strings.Split(q.Priorities, ",")[0])
		filter.Priority = &priority
	}

	tasks, nextCursor, err := h.service.List(c.Request.Context(), projectID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to list tasks"))
		return
	}

	resp := make([]response.TaskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toTaskResponse(&t)
	}

	total := len(tasks)
	cursorStr := ""
	if nextCursor != nil {
		cursorStr = nextCursor.UTC().Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, response.TaskListResponse{
		Tasks:      resp,
		Total:      total,
		NextCursor: cursorStr,
		HasMore:    len(tasks) == q.Limit,
	})
}

func (h *TaskHandler) Create(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid project id"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", err.Error()))
		return
	}

	var req request.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	taskType := domain_task.TypeTask
	if req.Type != "" {
		taskType = domain_task.Type(req.Type)
	}

	var assigneeID *uuid.UUID
	if req.AssigneeID != "" {
		if uid, err := uuid.Parse(req.AssigneeID); err == nil {
			assigneeID = &uid
		}
	}

	t, err := h.service.Create(c.Request.Context(), projectID, userID, req.Title, req.Description, taskType, assigneeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to create task"))
		return
	}

	c.JSON(http.StatusCreated, toTaskResponse(t))
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	taskIDStr := c.Param("taskId")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid task id"))
		return
	}

	t, err := h.service.GetByID(c.Request.Context(), taskID)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to get task"
		if errors.Is(err, application_task.ErrTaskNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(t))
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskIDStr := c.Param("taskId")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid task id"))
		return
	}

	var req request.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	taskType := domain_task.Type("")
	if req.Type != "" {
		taskType = domain_task.Type(req.Type)
	}
	priority := domain_task.Priority("")
	if req.Priority != "" {
		priority = domain_task.Priority(req.Priority)
	}
	status := domain_task.Status("")
	if req.Status != "" {
		status = domain_task.Status(req.Status)
	}

	var assigneeID *uuid.UUID
	if req.AssigneeID != nil {
		if *req.AssigneeID == "" {
			empty := uuid.Nil
			assigneeID = &empty
		} else {
			if uid, err := uuid.Parse(*req.AssigneeID); err == nil {
				assigneeID = &uid
			}
		}
	}

	t, err := h.service.Update(c.Request.Context(), taskID, req.Title, req.Description, taskType, priority, status, assigneeID, req.VersionID)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to update task"
		if errors.Is(err, application_task.ErrTaskNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		} else if errors.Is(err, application_task.ErrVersionConflict) {
			code = http.StatusConflict
			errCode = "CONFLICT"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(t))
}

func (h *TaskHandler) Move(c *gin.Context) {
	taskIDStr := c.Param("taskId")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid task id"))
		return
	}

	var req request.MoveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	columnID, err := uuid.Parse(req.ColumnID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid column id"))
		return
	}

	t, err := h.service.Move(c.Request.Context(), taskID, columnID, req.Position, req.VersionID)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to move task"
		if errors.Is(err, application_task.ErrTaskNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		} else if errors.Is(err, application_task.ErrVersionConflict) {
			code = http.StatusConflict
			errCode = "CONFLICT"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, response.TaskPositionResponse{
		TaskID:    t.ID.String(),
		ColumnID:  t.ColumnID.String(),
		Position:  t.Position,
		VersionID: "ok",
	})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskIDStr := c.Param("taskId")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid task id"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), taskID); err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to delete task"
		if errors.Is(err, application_task.ErrTaskNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}

func (h *TaskHandler) Restore(c *gin.Context) {
	taskIDStr := c.Param("taskId")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid task id"))
		return
	}

	t, err := h.service.Restore(c.Request.Context(), taskID)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to restore task"
		if errors.Is(err, application_task.ErrTaskNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(t))
}

func toTaskResponse(t *domain_task.Task) response.TaskResponse {
	r := response.TaskResponse{
		ID:          t.ID.String(),
		ProjectID:   t.ProjectID.String(),
		Title:       t.Title,
		Description: t.Description,
		Type:        string(t.Type),
		Priority:    string(t.Priority),
		Status:      string(t.Status),
		Position:    t.Position,
		Tags:        []response.TaskTag{},
		VersionID:   t.VersionID,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if t.ColumnID != nil {
		r.ColumnID = t.ColumnID.String()
	}
	return r
}
