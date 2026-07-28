package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/api/dto/response"
	application_org "github.com/devflow/devflow-backend/internal/application/organization"
	domain_org "github.com/devflow/devflow-backend/internal/domain/organization"

	"github.com/gin-gonic/gin"
)

var (
	ErrOrgNotFound   = errors.New("organization not found")
	ErrOrgForbidden  = errors.New("access to organization denied")
	ErrInvalidUserID = errors.New("invalid user id in context")
)

type OrgHandler struct {
	service *application_org.Service
}

func NewOrgHandler(service *application_org.Service) *OrgHandler {
	return &OrgHandler{service: service}
}

type createOrgRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
	Slug string `json:"slug" binding:"required,min=2,max=50"`
}

type updateOrgRequest struct {
	Name *string `json:"name" binding:"omitempty,min=1,max=100"`
	Slug *string `json:"slug" binding:"omitempty,min=2,max=50"`
}

func (h *OrgHandler) Create(c *gin.Context) {
	var req createOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", err.Error()))
		return
	}

	org, err := h.service.Create(c.Request.Context(), req.Name, req.Slug, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to create organization"))
		return
	}

	c.JSON(http.StatusCreated, toOrgResponse(org))
}

func (h *OrgHandler) GetByID(c *gin.Context) {
	orgIDStr := c.Param("orgId")

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid organization id"))
		return
	}

	org, err := h.service.GetByID(c.Request.Context(), orgID)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to get organization"
		if errors.Is(err, application_org.ErrOrgNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toOrgResponse(org))
}

func (h *OrgHandler) List(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", err.Error()))
		return
	}

	orgs, err := h.service.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to list organizations"))
		return
	}

	resp := make([]orgResponse, len(orgs))
	for i, org := range orgs {
		resp[i] = toOrgResponse(&org)
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrgHandler) Update(c *gin.Context) {
	orgIDStr := c.Param("orgId")

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid organization id"))
		return
	}

	var req updateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	name := ""
	slug := ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.Slug != nil {
		slug = *req.Slug
	}

	org, err := h.service.Update(c.Request.Context(), orgID, name, slug)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to update organization"
		if errors.Is(err, application_org.ErrOrgNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toOrgResponse(org))
}

func (h *OrgHandler) Delete(c *gin.Context) {
	orgIDStr := c.Param("orgId")

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid organization id"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), orgID); err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to delete organization"
		if errors.Is(err, application_org.ErrOrgNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "organization deleted successfully"})
}

func getUserID(c *gin.Context) (uuid.UUID, error) {
	uid, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, ErrInvalidUserID
	}
	id, ok := uid.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidUserID
	}
	return id, nil
}

type orgResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toOrgResponse(o *domain_org.Organization) orgResponse {
	return orgResponse{
		ID:        o.ID.String(),
		Name:      o.Name,
		Slug:      o.Slug,
		CreatedAt: o.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: o.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
