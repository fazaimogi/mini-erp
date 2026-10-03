package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/fazasuny/erp-system/internal/services"
)

type VoucherHandler struct {
	service services.VoucherService
}

func NewVoucherHandler(service services.VoucherService) *VoucherHandler {
	return &VoucherHandler{service: service}
}

// List godoc
// GET /api/vouchers?page=1&limit=20
func (h *VoucherHandler) List(c *gin.Context) {
	page, limit := pagination(c)

	vouchers, total, err := h.service.List(page, limit)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, vouchers, page, limit, total)
}

// GetByID godoc
// GET /api/vouchers/:id
func (h *VoucherHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid voucher id")
		return
	}

	voucher, err := h.service.GetByID(id)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, voucher, "success")
}

// Create godoc
// POST /api/vouchers
func (h *VoucherHandler) Create(c *gin.Context) {
	var input services.VoucherInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	voucher, err := h.service.Create(&input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusCreated, voucher, "voucher created")
}

// Update godoc
// PUT /api/vouchers/:id
func (h *VoucherHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid voucher id")
		return
	}

	var input services.VoucherInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	voucher, err := h.service.Update(id, &input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, voucher, "voucher updated")
}

// Delete godoc
// DELETE /api/vouchers/:id
func (h *VoucherHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid voucher id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, nil, "voucher deleted")
}

type validateVoucherRequest struct {
	Code        string          `json:"code" binding:"required"`
	TotalAmount decimal.Decimal `json:"total_amount"`
}

// Validate godoc
// POST /api/vouchers/validate
func (h *VoucherHandler) Validate(c *gin.Context) {
	var req validateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	if req.TotalAmount.IsNegative() {
		Error(c, http.StatusUnprocessableEntity, "INVALID_INPUT", "total_amount cannot be negative")
		return
	}

	validation, err := h.service.Validate(req.Code, req.TotalAmount)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, validation, "voucher valid")
}