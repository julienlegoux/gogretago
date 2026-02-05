package main

import (
	"fmt"
	"log"

	"github.com/lgxju/gogretago/config"
	"github.com/lgxju/gogretago/internal/infrastructure/di"
	"github.com/lgxju/gogretago/internal/presentation/routes"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting server in %s mode", cfg.AppEnv)

	// Initialize dependency container
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Setup router
	router := routes.SetupRouter(container)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
