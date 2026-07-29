package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/devflow/devflow-backend/internal/api/dto/request"
	"github.com/devflow/devflow-backend/internal/api/dto/response"
	"github.com/devflow/devflow-backend/internal/api/middleware"
	"github.com/devflow/devflow-backend/internal/application/auth"
)

type authService interface {
	Register(ctx context.Context, email, password, displayName string) (*auth.AuthUser, error)
	Login(ctx context.Context, email, password string) (*auth.AuthUser, error)
	GetByID(ctx context.Context, id string) (*auth.AuthUser, error)
	Refresh(ctx context.Context, refreshToken string) (*auth.AuthUser, error)
	VerifyEmail(ctx context.Context, email string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, newPassword string) error
}

type AuthHandler struct {
	authService   authService
	rdb           *redis.Client
	jwtSecret     string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewAuthHandler(
	authService authService,
	rdb *redis.Client,
	jwtSecret string,
	refreshSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		rdb:           rdb,
		jwtSecret:     jwtSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		code := http.StatusInternalServerError
		msg := "registration failed"
		errCode := "INTERNAL_ERROR"
		c.JSON(code, response.NewErrorResponse(errCode, msg))
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("TOKEN_ERROR", "failed to generate tokens"))
		return
	}

	c.JSON(http.StatusCreated, response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(h.accessTTL.Seconds()),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("INVALID_CREDENTIALS", "invalid email or password"))
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("TOKEN_ERROR", "failed to generate tokens"))
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(h.accessTTL.Seconds()),
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	_, err := h.rdb.Get(c.Request.Context(), "blacklist:"+req.RefreshToken).Result()
	if err == nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("TOKEN_BLACKLISTED", "refresh token has been revoked"))
		return
	}

	user, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("INVALID_REFRESH_TOKEN", err.Error()))
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("TOKEN_ERROR", "failed to generate tokens"))
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(h.accessTTL.Seconds()),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, ok := userID.(string)
	if !ok || uid == "" {
		c.JSON(http.StatusUnauthorized, response.NewErrorResponse("UNAUTHORIZED", "user not found"))
		return
	}

	user, err := h.authService.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewErrorResponse("NOT_FOUND", "user not found"))
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   "",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	if err := h.rdb.Set(c.Request.Context(), "blacklist:"+req.RefreshToken, "revoked", h.refreshTTL).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to logout"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req request.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	email, err := h.rdb.Get(c.Request.Context(), "verify_email_token:"+req.Token).Result()
	if err != nil || email == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_TOKEN", "verification token is invalid or expired"))
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), email); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "email verification failed"))
		return
	}

	h.rdb.Del(c.Request.Context(), "verify_email_token:"+req.Token)
	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	resetToken := uuid.New().String()
	if err := h.rdb.Set(c.Request.Context(), "reset_token:"+resetToken, req.Email, 15*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "failed to create reset token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "password reset email sent",
		"reset_token": resetToken,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	email, err := h.rdb.Get(c.Request.Context(), "reset_token:"+req.Token).Result()
	if err != nil || email == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_TOKEN", "reset token is invalid or expired"))
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), email, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse("INTERNAL_ERROR", "password reset failed"))
		return
	}

	h.rdb.Del(c.Request.Context(), "reset_token:"+req.Token)
	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *AuthHandler) generateTokens(user *auth.AuthUser) (string, string, error) {
	now := time.Now()

	accessClaims := &middleware.Claims{
		UserID: user.ID,
		OrgID:  user.OrgID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(h.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &middleware.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(h.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString([]byte(h.refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}
