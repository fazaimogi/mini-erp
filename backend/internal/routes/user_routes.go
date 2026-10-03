package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

func RegisterUserRoutes(api *gin.RouterGroup, handler *handlers.UserHandler, tokenManager *services.TokenManager) {
	users := api.Group("/users")
	users.Use(middleware.RequireAuth(tokenManager))
	users.Use(middleware.RequireRole("admin"))

	users.GET("", handler.List)
	users.GET("/:id", handler.GetByID)
	users.POST("", handler.Create)
	users.PUT("/:id", handler.Update)
	users.DELETE("/:id", handler.Delete)
	users.PUT("/:id/role", handler.ChangeRole)
	users.PUT("/:id/status", handler.ChangeStatus)

	// Roles feed the New User form dropdown (frontend GETs /api/roles).
	api.GET("/roles", middleware.RequireAuth(tokenManager), middleware.RequireRole("admin"), handler.ListRoles)
}
