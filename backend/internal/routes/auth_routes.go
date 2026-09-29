package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/middleware"
	"github.com/fazasuny/erp-system/internal/services"
)

// RegisterAuthRoutes registers the Authentication API (PRD section 5).
// register and login are public; logout and me require a valid token.
func RegisterAuthRoutes(rg *gin.RouterGroup, h *handlers.AuthHandler, tokens *services.TokenManager) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)

		protected := auth.Group("")
		protected.Use(middleware.RequireAuth(tokens))
		{
			protected.POST("/logout", h.Logout)
			protected.GET("/me", h.Me)
		}
	}
}
