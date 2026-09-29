package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterCustomerRoutes registers the Customer API endpoints (PRD section 11).
// GET list/search and GET detail require auth for any role.
// POST, PUT require auth for admin or staff.
// DELETE requires auth for admin only.
func RegisterCustomerRoutes(rg *gin.RouterGroup, h *handlers.CustomerHandler, tokens *services.TokenManager) {
	customers := rg.Group("/customers")
	customers.Use(middleware.RequireAuth(tokens))
	{
		customers.GET("", h.List)
		customers.GET("/:id", h.GetByID)

		// admin or staff can create/update
		staffUp := customers.Group("")
		staffUp.Use(middleware.RequireRole("admin", "staff"))
		{
			staffUp.POST("", h.Create)
			staffUp.PUT("/:id", h.Update)
		}

		// admin only can delete
		adminOnly := customers.Group("")
		adminOnly.Use(middleware.RequireRole("admin"))
		{
			adminOnly.DELETE("/:id", h.Delete)
		}
	}
}
