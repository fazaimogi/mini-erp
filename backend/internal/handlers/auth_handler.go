package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// POST /api/auth/register
// New accounts receive the staff role; admin accounts are created through user
// management (PRD section 4).
func (h *AuthHandler) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.service.Register(&input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusCreated, user, "user registered")
}

// Login godoc
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	user, token, err := h.service.Login(&input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	}, "login success")
}

// Logout godoc
// POST /api/auth/logout
// JWTs are stateless, so logout is a client-side token discard acknowledged by
// the API for a symmetric frontend flow (PRD section 5).
func (h *AuthHandler) Logout(c *gin.Context) {
	Success(c, http.StatusOK, nil, "logout success")
}

// Me godoc
// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		Error(c, http.StatusUnauthorized, services.ErrUnauthorized.Code, services.ErrUnauthorized.Message)
		return
	}

	user, err := h.service.Me(userID)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, user, "success")
}
