package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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

	// Create an initial minimal router with health check so Railway
	// can connect to the port immediately, before DB is ready.
	initialRouter := gin.Default()
	initialRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	initialRouter.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Service is starting up",
		})
	})

	// Use atomic handler swap so the server can start immediately
	// and upgrade to the full router once dependencies are ready.
	var handler atomic.Value
	handler.Store(http.Handler(initialRouter))

	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.Load().(http.Handler).ServeHTTP(w, r)
		}),
	}

	// Start HTTP server immediately so health checks pass
	go func() {
		log.Printf("Server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Initialize dependencies (DB connection, migrations — may be slow)
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Setup full router and swap it in
	fullRouter := routes.SetupRouter(container)
	handler.Store(http.Handler(fullRouter))
	log.Println("Application fully initialized")

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if err := server.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	cancel()
	log.Println("Server exiting")
}
