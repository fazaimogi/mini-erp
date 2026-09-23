package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
)

// RegisterProductRoutes registers the Product API endpoints (PRD section 7.3).
func RegisterProductRoutes(rg *gin.RouterGroup, h *handlers.ProductHandler) {
	products := rg.Group("/products")
	{
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)
		products.POST("", h.Create)
		products.PUT("/:id", h.Update)
		products.DELETE("/:id", h.Delete)
	}
}
