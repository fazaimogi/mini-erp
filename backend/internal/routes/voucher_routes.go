package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterVoucherRoutes registers the Voucher API.
// Management (CRUD) is admin-only; validation is open to admin and staff
// because the sales checkout calls it.
func RegisterVoucherRoutes(rg *gin.RouterGroup, h *handlers.VoucherHandler, tokens *services.TokenManager) {
	vouchers := rg.Group("/vouchers")
	vouchers.Use(middleware.RequireAuth(tokens))
	{
		checkout := vouchers.Group("")
		checkout.Use(middleware.RequireRole("admin", "staff"))
		{
			checkout.POST("/validate", h.Validate)
		}

		adminOnly := vouchers.Group("")
		adminOnly.Use(middleware.RequireRole("admin"))
		{
			adminOnly.GET("", h.List)
			adminOnly.GET("/:id", h.GetByID)
			adminOnly.POST("", h.Create)
			adminOnly.PUT("/:id", h.Update)
			adminOnly.DELETE("/:id", h.Delete)
		}
	}
}