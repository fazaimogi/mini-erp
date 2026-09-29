package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

type SaleHandler struct {
	service services.SaleService
}

func NewSaleHandler(service services.SaleService) *SaleHandler {
	return &SaleHandler{service: service}
}

// GetSales returns paginated sales list.
// GET /api/sales?page=1&limit=20
func (h *SaleHandler) GetSales(c *gin.Context) {
	page, limit := pagination(c)

	sales, total, err := h.service.List(page, limit)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, sales, page, limit, total)
}

// GetSaleByID returns single sale with items.
// GET /api/sales/:id
func (h *SaleHandler) GetSaleByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid sale id")
		return
	}

	sale, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "sale not found" {
			Error(c, http.StatusNotFound, "SALE_NOT_FOUND", "sale not found")
			return
		}
		respondInternalError(c, err)
		return
	}

	Success(c, http.StatusOK, sale, "success")
}

// CreateSale creates new sale with items and deducts stock atomically.
// POST /api/sales
func (h *SaleHandler) CreateSale(c *gin.Context) {
	var input services.CreateSaleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "user not found in context")
		return
	}

	sale, err := h.service.Create(&input, userID.(uint))
	if err != nil {
		errMsg := err.Error()

		if len(errMsg) > len("insufficient") && errMsg[:len("insufficient")] == "insufficient" {
			Error(c, http.StatusConflict, "INSUFFICIENT_STOCK", errMsg)
			return
		}
		if len(errMsg) > len("invalid") && errMsg[:len("invalid")] == "invalid" {
			Error(c, http.StatusBadRequest, "INVALID_INPUT", errMsg)
			return
		}

		respondInternalError(c, err)
		return
	}

	Success(c, http.StatusCreated, sale, "sale created")
}
