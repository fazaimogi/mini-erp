package handlers

import (
	"net/http"

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
	page, limit := pagination(c)

	products, total, err := h.service.List(page, limit)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, products, page, limit, total)
}

// LowStock godoc
// GET /api/products/low-stock?page=1&limit=20
// Products where stock <= minimum_stock (PRD section 6.3).
func (h *ProductHandler) LowStock(c *gin.Context) {
	page, limit := pagination(c)

	products, total, err := h.service.ListLowStock(page, limit)
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
