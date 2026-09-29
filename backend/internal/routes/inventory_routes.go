package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterInventoryRoutes registers the Inventory API (PRD section 23).
//
// PRD section 4 grants staff "Update" on inventory, so stock movements are open
// to both roles while browsing is read-only for them. There is no adjustment
// endpoint: PRD section 23 does not list one.
func RegisterInventoryRoutes(rg *gin.RouterGroup, h *handlers.InventoryHandler, tokens *services.TokenManager) {
	inventory := rg.Group("/inventory")
	inventory.Use(middleware.RequireAuth(tokens))
	{
		// "/history" is registered before "/:id"; gin rejects a wildcard that
		// collides with a static sibling.
		inventory.GET("/history", h.History)
		inventory.GET("", h.List)
		inventory.GET("/:id", h.GetByID)

		mutations := inventory.Group("")
		mutations.Use(middleware.RequireRole("admin", "staff"))
		{
			mutations.POST("/stock-in", h.StockIn)
			mutations.POST("/stock-out", h.StockOut)
		}
	}
}
