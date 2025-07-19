package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pteronimbus/pteronimbus/apps/controller/internal/client"
	"github.com/pteronimbus/pteronimbus/apps/controller/internal/handlers"
	"github.com/pteronimbus/pteronimbus/apps/controller/internal/heartbeat"
)

func main() {
	// Configuration
	backendURL := getEnv("BACKEND_URL", "http://localhost:8080")
	clusterID := getEnv("CLUSTER_ID", "default-cluster")
	clusterName := getEnv("CLUSTER_NAME", "Default Cluster")
	version := getEnv("CONTROLLER_VERSION", "0.1.0")

	// Create backend client
	backendClient := client.NewBackendClient(backendURL, clusterID, clusterName, version)

	// Perform initial handshake
	log.Println("Performing handshake with backend...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := backendClient.Handshake(ctx); err != nil {
		log.Fatalf("Handshake failed: %v", err)
	}

	log.Printf("Handshake successful! Controller ID: %s", backendClient.GetControllerID())

	// Create heartbeat manager
	heartbeatInterval := time.Duration(backendClient.GetHeartbeatTTL()) * time.Second
	if heartbeatInterval == 0 {
		heartbeatInterval = 5 * time.Second // Default to 5 seconds
	}

	heartbeatManager := heartbeat.NewManager(backendClient, heartbeatInterval)

	// Start heartbeat manager
	heartbeatCtx, heartbeatCancel := context.WithCancel(context.Background())
	defer heartbeatCancel()

	if err := heartbeatManager.Start(heartbeatCtx); err != nil {
		log.Fatalf("Failed to start heartbeat manager: %v", err)
	}

	// Initialize handlers
	h := handlers.NewHealthHandler()

	// Setup HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h.Routes(),
	}

	// Channel to listen for interrupt signal to terminate server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit
	log.Println("Shutting down server...")

	// Stop heartbeat manager
	heartbeatManager.Stop()

	// Give outstanding requests 30 seconds to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
