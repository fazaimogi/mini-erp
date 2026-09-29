package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterReportsRoutes registers the Reports API endpoints (PRD section 6, 16).
// All routes require admin or staff role.
func RegisterReportsRoutes(rg *gin.RouterGroup, h *handlers.ReportsHandler, tokens *services.TokenManager) {
	reports := rg.Group("/reports")
	reports.Use(middleware.RequireAuth(tokens))
	reports.Use(middleware.RequireRole("admin", "staff"))
	{
		reports.GET("/dashboard", h.Dashboard)
		reports.GET("/sales", h.Sales)
		reports.GET("/products", h.Products)
	}
}
