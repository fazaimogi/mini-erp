package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// List godoc
// GET /api/categories?page=1&limit=20
func (h *CategoryHandler) List(c *gin.Context) {
	page, limit := pagination(c)

	categories, total, err := h.service.List(page, limit)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, categories, page, limit, total)
}

// GetByID godoc
// GET /api/categories/:id
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid category id")
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, category, "success")
}

// Create godoc
// POST /api/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var input services.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	category, err := h.service.Create(&input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusCreated, category, "category created")
}

// Update godoc
// PUT /api/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid category id")
		return
	}

	var input services.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	category, err := h.service.Update(id, &input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, category, "category updated")
}

// Delete godoc
// DELETE /api/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid category id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, nil, "category deleted")
}
