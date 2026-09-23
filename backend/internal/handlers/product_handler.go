package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// List godoc
// GET /api/products?page=1&limit=20
func (h *ProductHandler) List(c *gin.Context) {
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	limit := parsePositiveInt(c.DefaultQuery("limit", "20"), 20)
	if limit > 100 {
		limit = 20
	}

	products, total, err := h.service.List(page, limit)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, products, page, limit, total)
}

// GetByID godoc
// GET /api/products/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid product id")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, product, "success")
}

// Create godoc
// POST /api/products
func (h *ProductHandler) Create(c *gin.Context) {
	var input services.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	product, err := h.service.Create(&input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusCreated, product, "product created")
}

// Update godoc
// PUT /api/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid product id")
		return
	}

	var input services.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	product, err := h.service.Update(id, &input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, product, "product updated")
}

// Delete godoc
// DELETE /api/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid product id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, nil, "product deleted")
}

func parseIDParam(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func parsePositiveInt(value string, fallback int) int {
	number, err := strconv.Atoi(value)
	if err != nil || number < 1 {
		return fallback
	}
	return number
}

func respondServiceError(c *gin.Context, err error) {
	var appErr *services.AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
		return
	}
	respondInternalError(c, err)
}

func respondInternalError(c *gin.Context, err error) {
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
}
