package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/devflow/devflow-backend/internal/api/dto/response"
	application_comment "github.com/devflow/devflow-backend/internal/application/comment"
	domain_comment "github.com/devflow/devflow-backend/internal/domain/comment"

	"github.com/gin-gonic/gin"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

type CommentHandler struct {
	service *application_comment.Service
}

func NewCommentHandler(service *application_comment.Service) *CommentHandler {
	return &CommentHandler{service: service}
}

type createCommentRequest struct {
	Body string `json:"body" binding:"required,min=1,max=10000"`
}

type updateCommentRequest struct {
	Body *string `json:"body" binding:"omitempty,min=1,max=10000"`
}

func (h *CommentHandler) Create(c *gin.Context) {
	taskIDStr := c.Param("taskId")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid task id"))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", err.Error()))
		return
	}

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	comment, err := h.service.Create(c.Request.Context(), taskID, userID, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to create comment"))
		return
	}

	c.JSON(http.StatusCreated, toCommentResponse(comment))
}

func (h *CommentHandler) ListByTask(c *gin.Context) {
	taskIDStr := c.Param("taskId")

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid task id"))
		return
	}

	comments, err := h.service.ListByTaskID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to list comments"))
		return
	}

	resp := make([]commentResponse, len(comments))
	for i, c := range comments {
		resp[i] = toCommentResponse(&c)
	}

	c.JSON(http.StatusOK, resp)
}

func (h *CommentHandler) GetByID(c *gin.Context) {
	commentIDStr := c.Param("commentId")

	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid comment id"))
		return
	}

	comment, err := h.service.GetByID(c.Request.Context(), commentID)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to get comment"
		if errors.Is(err, application_comment.ErrCommentNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toCommentResponse(comment))
}

func (h *CommentHandler) Update(c *gin.Context) {
	commentIDStr := c.Param("commentId")

	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid comment id"))
		return
	}

	var req updateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	body := ""
	if req.Body != nil {
		body = *req.Body
	}

	comment, err := h.service.Update(c.Request.Context(), commentID, body)
	if err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to update comment"
		if errors.Is(err, application_comment.ErrCommentNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, toCommentResponse(comment))
}

func (h *CommentHandler) Delete(c *gin.Context) {
	commentIDStr := c.Param("commentId")

	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", "invalid comment id"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), commentID); err != nil {
		code := http.StatusInternalServerError
		errCode := "INTERNAL_ERROR"
		msg := "failed to delete comment"
		if errors.Is(err, application_comment.ErrCommentNotFound) {
			code = http.StatusNotFound
			errCode = "NOT_FOUND"
			msg = err.Error()
		}
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}

type commentResponse struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	AuthorID  string `json:"author_id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toCommentResponse(c *domain_comment.Comment) commentResponse {
	return commentResponse{
		ID:        c.ID.String(),
		TaskID:    c.TaskID.String(),
		AuthorID:  c.AuthorID.String(),
		Body:      c.Body,
		CreatedAt: c.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
