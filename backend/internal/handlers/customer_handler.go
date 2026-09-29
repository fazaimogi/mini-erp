package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

type CustomerHandler struct {
	service services.CustomerService
}

func NewCustomerHandler(service services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// List godoc
// GET /api/customers?page=1&limit=20&q=search
func (h *CustomerHandler) List(c *gin.Context) {
	query := c.Query("q")
	page, limit := pagination(c)

	var customers []interface{}
	var total int64
	var err error

	if query != "" {
		custs, t, e := h.service.Search(query, page, limit)
		for i := range custs {
			customers = append(customers, custs[i])
		}
		total = t
		err = e
	} else {
		custs, t, e := h.service.List(page, limit)
		for i := range custs {
			customers = append(customers, custs[i])
		}
		total = t
		err = e
	}

	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, customers, page, limit, total)
}

// GetByID godoc
// GET /api/customers/:id
func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid customer id")
		return
	}

	customer, err := h.service.GetByID(id)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, customer, "success")
}

// Create godoc
// POST /api/customers
func (h *CustomerHandler) Create(c *gin.Context) {
	var input services.CustomerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	customer, err := h.service.Create(&input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusCreated, customer, "customer created")
}

// Update godoc
// PUT /api/customers/:id
func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid customer id")
		return
	}

	var input services.CustomerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}

	customer, err := h.service.Update(id, &input)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, customer, "customer updated")
}

// Delete godoc
// DELETE /api/customers/:id
func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_ID", "invalid customer id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, nil, "customer deleted")
}
