package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/fazasuny/erp-system/internal/app"
	"github.com/fazasuny/erp-system/internal/config"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	engine, _, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
