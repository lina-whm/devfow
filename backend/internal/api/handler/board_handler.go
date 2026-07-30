package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/devflow/devflow-backend/internal/api/dto/response"
	application_board "github.com/devflow/devflow-backend/internal/application/board"

	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	service *application_board.Service
}

func NewBoardHandler(service *application_board.Service) *BoardHandler {
	return &BoardHandler{service: service}
}

type columnResponse struct {
	ID       string                  `json:"id"`
	BoardID  string                  `json:"board_id"`
	Name     string                  `json:"name"`
	Position int                     `json:"position"`
	WIPLimit int                     `json:"wip_limit"`
	Tasks    []response.TaskResponse `json:"tasks"`
}

type boardResponse struct {
	ID        string           `json:"id"`
	ProjectID string           `json:"project_id"`
	Name      string           `json:"name"`
	Columns   []columnResponse `json:"columns"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

func (h *BoardHandler) GetByProject(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid project id"))
		return
	}

	result, err := h.service.GetOrCreate(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to get board"))
		return
	}

	colResponses := make([]columnResponse, len(result.Columns))
	for i, cwt := range result.Columns {
		taskResponses := make([]response.TaskResponse, len(cwt.Tasks))
		for j := range cwt.Tasks {
			t := &cwt.Tasks[j]
			taskResponses[j] = toTaskResponse(t)
		}
		colResponses[i] = columnResponse{
			ID:       cwt.Column.ID.String(),
			BoardID:  cwt.Column.BoardID.String(),
			Name:     cwt.Column.Name,
			Position: int(cwt.Column.Position),
			WIPLimit: cwt.Column.WIPLimit,
			Tasks:    taskResponses,
		}
	}

	c.JSON(http.StatusOK, boardResponse{
		ID:        result.Board.ID.String(),
		ProjectID: result.Board.ProjectID.String(),
		Name:      result.Board.Name,
		Columns:   colResponses,
		CreatedAt: result.Board.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: result.Board.UpdatedAt.UTC().Format(time.RFC3339),
	})
}

type updateColumnsRequest struct {
	Columns []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"columns" binding:"required"`
}

func (h *BoardHandler) UpdateColumns(c *gin.Context) {
	boardIDStr := c.Param("boardId")

	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid board id"))
		return
	}

	var req updateColumnsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	updates := make([]application_board.ColumnUpdate, len(req.Columns))
	for i, col := range req.Columns {
		updates[i] = application_board.ColumnUpdate{
			ID:   col.ID,
			Name: col.Name,
		}
	}

	if err := h.service.UpdateColumns(c.Request.Context(), boardID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to update columns"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
