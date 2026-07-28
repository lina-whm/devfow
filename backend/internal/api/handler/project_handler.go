package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/api/dto/response"
	application_project "github.com/devflow/devflow-backend/internal/application/project"
	domain_project "github.com/devflow/devflow-backend/internal/domain/project"

	"github.com/gin-gonic/gin"
)

var (
	ErrProjectNotFound  = errors.New("project not found")
	ErrProjectForbidden = errors.New("access to project denied")
)

type ProjectHandler struct {
	service *application_project.Service
}

func NewProjectHandler(service *application_project.Service) *ProjectHandler {
	return &ProjectHandler{service: service}
}

type createProjectRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Key         string  `json:"key" binding:"required,min=2,max=10"`
	Description string  `json:"description"`
}

type updateProjectRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	orgIDStr := c.Param("orgId")

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid organization id"))
		return
	}

	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	project, err := h.service.Create(c.Request.Context(), req.Name, req.Key, req.Description, orgID, nil)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to create project"
		if errors.Is(err, application_project.ErrProjectKeyTaken) {
			code = http.StatusConflict
			errCode = "CONFLICT"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusCreated, toProjectResponse(project))
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid project id"))
		return
	}

	project, err := h.service.GetByID(c.Request.Context(), projectID)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to get project"
		if errors.Is(err, application_project.ErrProjectNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toProjectResponse(project))
}

func (h *ProjectHandler) List(c *gin.Context) {
	orgIDStr := c.Param("orgId")

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid organization id"))
		return
	}

	projects, err := h.service.ListByOrgID(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to list projects"))
		return
	}

	resp := make([]projectResponse, len(projects))
	for i, p := range projects {
		resp[i] = toProjectResponse(&p)
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid project id"))
		return
	}

	var req updateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	name := ""
	description := ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}

	project, err := h.service.Update(c.Request.Context(), projectID, name, description, "")
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to update project"
		if errors.Is(err, application_project.ErrProjectNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toProjectResponse(project))
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	projectIDStr := c.Param("projectId")

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid project id"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), projectID); err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to delete project"
		if errors.Is(err, application_project.ErrProjectNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project deleted successfully"})
}

type projectResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Key            string `json:"key"`
	Status         string `json:"status"`
	Description    string `json:"description,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func toProjectResponse(p *domain_project.Project) projectResponse {
	desc := ""
	if p.Description != nil {
		desc = *p.Description
	}
	return projectResponse{
		ID:             p.ID.String(),
		OrganizationID: p.OrganizationID.String(),
		Name:           p.Name,
		Key:            p.Key,
		Status:         string(p.Status),
		Description:    desc,
		CreatedAt:      p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
