package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterSalesRoutes registers the Sales API endpoints (PRD section 9).
// All routes require admin or staff role.
func RegisterSalesRoutes(
	api *gin.RouterGroup,
	handler *handlers.SaleHandler,
	tokenManager *services.TokenManager,
) {
	sales := api.Group("/sales")
	sales.Use(middleware.RequireAuth(tokenManager))
	sales.Use(middleware.RequireRole("admin", "staff"))
	{
		sales.GET("", handler.GetSales)
		sales.GET("/:id", handler.GetSaleByID)
		sales.POST("", handler.CreateSale)
	}
}
