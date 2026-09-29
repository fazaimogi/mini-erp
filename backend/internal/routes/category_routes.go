package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterCategoryRoutes registers the Category API endpoints (PRD section 8).
// All routes require auth; writes are restricted to admin.
func RegisterCategoryRoutes(rg *gin.RouterGroup, h *handlers.CategoryHandler, tokens *services.TokenManager) {
	categories := rg.Group("/categories")
	categories.Use(middleware.RequireAuth(tokens))
	{
		categories.GET("", h.List)
		categories.GET("/:id", h.GetByID)

		adminOnly := categories.Group("")
		adminOnly.Use(middleware.RequireRole("admin"))
		{
			adminOnly.POST("", h.Create)
			adminOnly.PUT("/:id", h.Update)
			adminOnly.DELETE("/:id", h.Delete)
		}
	}
}
