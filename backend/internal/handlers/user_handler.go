package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// List godoc
// GET /api/users?page=1&limit=20
// Admin only. Returns paginated list of all users.
func (h *UserHandler) List(c *gin.Context) {
	page, limit := pagination(c)

	users, total, err := h.service.List(page, limit)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, users, page, limit, total)
}

// GetByID godoc
// GET /api/users/:id
// Admin only. Returns user detail.
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid user id")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, user, "success")
}

// Create godoc
// POST /api/users
// Admin only. Creates new user with name, email, password, role_id.
func (h *UserHandler) Create(c *gin.Context) {
	var input services.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.service.Create(&input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusCreated, user, "user created")
}

// Update godoc
// PUT /api/users/:id
// Admin only. Updates user name, email, role_id, status.
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid user id")
		return
	}

	var input services.UserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.service.Update(id, &input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, user, "user updated")
}

// Delete godoc
// DELETE /api/users/:id
// Admin only. Deletes user.
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid user id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, nil, "user deleted")
}

// ListRoles godoc
// GET /api/roles
// Admin only. Returns all roles for the New User role dropdown.
func (h *UserHandler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Success(c, http.StatusOK, roles, "success")
}

// ChangeRole godoc
// PUT /api/users/:id/role
// Admin only. Changes user role.
func (h *UserHandler) ChangeRole(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid user id")
		return
	}

	var input struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.service.ChangeRole(id, input.RoleID)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, user, "role changed")
}

// ChangeStatus godoc
// PUT /api/users/:id/status
// Admin only. Activates or deactivates user.
func (h *UserHandler) ChangeStatus(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid user id")
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.service.ChangeStatus(id, input.Status)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, user, "status changed")
}
