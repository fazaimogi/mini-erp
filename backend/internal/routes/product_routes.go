package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterProductRoutes registers the Product API endpoints (PRD section 7.3).
// Every route requires a token; writes are restricted to admin per the role
// matrix in PRD section 4 (staff has read-only access to products).
func RegisterProductRoutes(rg *gin.RouterGroup, h *handlers.ProductHandler, tokens *services.TokenManager) {
	products := rg.Group("/products")
	products.Use(middleware.RequireAuth(tokens))
	{
		// "/low-stock" is registered before "/:id"; gin rejects a wildcard that
		// collides with a static sibling.
		products.GET("/low-stock", h.LowStock)
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)

		adminOnly := products.Group("")
		adminOnly.Use(middleware.RequireRole("admin"))
		{
			adminOnly.POST("", h.Create)
			adminOnly.PUT("/:id", h.Update)
			adminOnly.DELETE("/:id", h.Delete)
		}
	}
}
