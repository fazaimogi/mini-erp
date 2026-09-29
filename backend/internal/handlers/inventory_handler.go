package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/services"
)

type InventoryHandler struct {
	service services.InventoryService
}

func NewInventoryHandler(service services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// List godoc
// GET /api/inventory
func (h *InventoryHandler) List(c *gin.Context) {
	page, limit := pagination(c)

	products, total, err := h.service.List(page, limit)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Collection(c, products, page, limit, total)
}

// GetByID godoc
// GET /api/inventory/:id
// The id is the product id, since stock is a product attribute.
func (h *InventoryHandler) GetByID(c *gin.Context) {
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

// StockIn godoc
// POST /api/inventory/stock-in
func (h *InventoryHandler) StockIn(c *gin.Context) {
	h.applyMovement(c, h.service.StockIn)
}

// StockOut godoc
// POST /api/inventory/stock-out
func (h *InventoryHandler) StockOut(c *gin.Context) {
	h.applyMovement(c, h.service.StockOut)
}

// History godoc
// GET /api/inventory/history?product_id=&page=&limit=
func (h *InventoryHandler) History(c *gin.Context) {
	page, limit := pagination(c)

	var productID uint
	if raw := c.Query("product_id"); raw != "" {
		parsed, err := parseUint(raw)
		if err != nil {
			Error(c, http.StatusBadRequest, "INVALID_ID", "invalid product id")
			return
		}
		productID = parsed
	}

	records, total, err := h.service.History(productID, page, limit)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Collection(c, records, page, limit, total)
}

// applyMovement runs the shared body of stock-in and stock-out.
func (h *InventoryHandler) applyMovement(
	c *gin.Context,
	run func(*services.StockMovementInput, uint) (*models.InventoryTransaction, error),
) {
	var input services.StockMovementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		Error(c, http.StatusUnauthorized, services.ErrUnauthorized.Code, services.ErrUnauthorized.Message)
		return
	}

	record, err := run(&input, userID)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, record, "stock updated")
}
