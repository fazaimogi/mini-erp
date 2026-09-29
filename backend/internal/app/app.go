// Package app is the composition root: it connects to PostgreSQL, applies
// migrations and wires the HTTP router. It exists so the exact wiring the
// server runs is the wiring the integration test exercises, instead of the test
// rebuilding an equivalent router by hand and drifting from production.
package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/fazasuny/erp-system/internal/config"
	"github.com/fazasuny/erp-system/internal/database"
	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/repositories"
	"github.com/fazasuny/erp-system/internal/routes"
	"github.com/fazasuny/erp-system/internal/services"
)

// New connects to the database, brings the schema up to date, then builds the
// router. Migrating here — before any repository is constructed — guarantees
// `roles` exists before the auth flow queries it.
func New(cfg config.Config) (*gin.Engine, *gorm.DB, error) {
	db, err := database.Connect(cfg)
	if err != nil {
		return nil, nil, err
	}

	if err := database.Migrate(db); err != nil {
		return nil, nil, fmt.Errorf("migrate: %w", err)
	}

	return Router(db, cfg), db, nil
}

// Router wires handlers, services and repositories onto a gin engine.
func Router(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ERP API is running"})
	})

	tokenManager := services.NewTokenManager(cfg.JWTSecret, cfg.JWTExpiryTime)

	userRepo := repositories.NewGormUserRepository(db)
	authHandler := handlers.NewAuthHandler(services.NewAuthService(userRepo, tokenManager))
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	productRepo := repositories.NewGormProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	inventoryRepo := repositories.NewGormInventoryRepository(db)
	inventoryHandler := handlers.NewInventoryHandler(
		services.NewInventoryService(productService, inventoryRepo),
	)
	reportsRepo := repositories.NewGormReportsRepository(db)
	reportsHandler := handlers.NewReportsHandler(
		services.NewReportsService(reportsRepo),
	)
	categoryRepo := repositories.NewGormCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	customerRepo := repositories.NewGormCustomerRepository(db)
	customerService := services.NewCustomerService(customerRepo)
	customerHandler := handlers.NewCustomerHandler(customerService)

	saleRepo := repositories.NewGormSaleRepository(db)
	saleService := services.NewSaleService(saleRepo, productRepo, inventoryRepo, db)
	saleHandler := handlers.NewSaleHandler(saleService)

	api := router.Group("/api")
	routes.RegisterAuthRoutes(api, authHandler, tokenManager)
	routes.RegisterProductRoutes(api, productHandler, tokenManager)
	routes.RegisterInventoryRoutes(api, inventoryHandler, tokenManager)
	routes.RegisterSalesRoutes(api, saleHandler, tokenManager)
	routes.RegisterReportsRoutes(api, reportsHandler, tokenManager)
	routes.RegisterCategoryRoutes(api, categoryHandler, tokenManager)
	routes.RegisterUserRoutes(api, userHandler, tokenManager)
	routes.RegisterCustomerRoutes(api, customerHandler, tokenManager)

	return router
}