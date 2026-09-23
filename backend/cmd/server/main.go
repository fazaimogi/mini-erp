package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/fazasuny/erp-system/internal/config"
	"github.com/fazasuny/erp-system/internal/database"
	"github.com/fazasuny/erp-system/internal/handlers"
	"github.com/fazasuny/erp-system/internal/repositories"
	"github.com/fazasuny/erp-system/internal/routes"
	"github.com/fazasuny/erp-system/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using environment variables")
	}

	cfg := config.Load()
	if missing := missingDBConfig(cfg); len(missing) > 0 {
		log.Fatalf("fatal: missing required database configuration: %s", strings.Join(missing, ", "))
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ERP API is running",
		})
	})

	productRepo := repositories.NewGormProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	api := router.Group("/api")
	routes.RegisterProductRoutes(api, productHandler)

	log.Println("server running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// missingDBConfig reports database settings that are absent from the
// environment. Failing fast here prevents a silently broken startup where the
// process runs but every request dies inside the repository layer.
func missingDBConfig(cfg config.Config) []string {
	required := map[string]string{
		"DB_HOST":     cfg.DBHost,
		"DB_PORT":     cfg.DBPort,
		"DB_USER":     cfg.DBUser,
		"DB_PASSWORD": cfg.DBPassword,
		"DB_NAME":     cfg.DBName,
	}

	missing := make([]string, 0, len(required))
	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}

	return missing
}
