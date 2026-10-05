package routes

import (
	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterActivityLogRoutes(rg *gin.RouterGroup, handler *handlers.ActivityLogHandler, tokens *services.TokenManager) {
	logs := rg.Group("/activity-logs")
	logs.Use(middleware.RequireAuth(tokens))
	logs.Use(middleware.RequireRole("admin"))
	{
		logs.GET("", handler.List)
	}
}
