package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

type ReportsHandler struct {
	service services.ReportsService
}

func NewReportsHandler(service services.ReportsService) *ReportsHandler {
	return &ReportsHandler{service: service}
}

// Dashboard godoc
// GET /api/reports/dashboard
// Admin & Staff only (PRD section 6.1).
// Responds with: total products, customers, sales, revenue, low stock count.
func (h *ReportsHandler) Dashboard(c *gin.Context) {
	metrics, err := h.service.Dashboard()
	if err != nil {
		respondServiceError(c, err)
		return
	}
	Success(c, http.StatusOK, metrics, "success")
}

// Sales godoc
// GET /api/reports/sales?period=today
// Admin & Staff only (PRD section 6.2, 16.1).
// Query params:
//   - period: today | week | month | custom (required)
//   - start_date: RFC3339 (required if period=custom)
//   - end_date: RFC3339 (required if period=custom)
func (h *ReportsHandler) Sales(c *gin.Context) {
	period := c.Query("period")
	if period == "" {
		Error(c, http.StatusUnprocessableEntity, "MISSING_PERIOD", "period parameter required")
		return
	}

	var startDate, endDate *time.Time

	if period == "custom" {
		startStr := c.Query("start_date")
		endStr := c.Query("end_date")

		if startStr == "" || endStr == "" {
			respondServiceError(c, services.ErrInvalidDateRange)
			return
		}

		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			Error(c, http.StatusUnprocessableEntity, "INVALID_DATE_FORMAT", "start_date must be RFC3339")
			return
		}

		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			Error(c, http.StatusUnprocessableEntity, "INVALID_DATE_FORMAT", "end_date must be RFC3339")
			return
		}

		startDate = &start
		endDate = &end
	}

	report, err := h.service.SalesReport(period, startDate, endDate)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	Success(c, http.StatusOK, report, "success")
}

// Products godoc
// GET /api/reports/products
// Admin & Staff only (PRD section 16.2).
// Responds with top selling and low stock products.
func (h *ReportsHandler) Products(c *gin.Context) {
	report, err := h.service.ProductsReport()
	if err != nil {
		respondServiceError(c, err)
		return
	}
	Success(c, http.StatusOK, report, "success")
}
