package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/handlers"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/middleware"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize services
	redisService := services.NewRedisService(cfg)
	jwtService := services.NewJWTService(cfg)
	discordService := services.NewDiscordService(cfg)
	authService := services.NewAuthService(discordService, jwtService, redisService)

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := redisService.Ping(ctx); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
		log.Println("Continuing without Redis - sessions will not persist")
	} else {
		log.Println("Redis connection established")
	}

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check routes
	router.GET("/health", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		healthHandler.Health(w, r)
	}))
	router.GET("/healthz", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		healthHandler.Health(w, r)
	}))
	router.GET("/ready", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		healthHandler.Ready(w, r)
	}))
	router.GET("/live", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		healthHandler.Live(w, r)
	}))

	// Authentication routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.GET("/login", authHandler.Login)
		authRoutes.GET("/callback", authHandler.Callback)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.GET("/me", authMiddleware.RequireAuth(), authHandler.Me)
		authRoutes.POST("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
	}

	// API routes (protected)
	apiRoutes := router.Group("/api")
	apiRoutes.Use(authMiddleware.RequireAuth())
	{
		// Placeholder for future API endpoints
		apiRoutes.GET("/test", func(c *gin.Context) {
			user, _ := middleware.GetUserFromContext(c)
			c.JSON(http.StatusOK, gin.H{
				"message": "Authenticated API endpoint",
				"user":    user,
			})
		})
	}

	// Setup HTTP server
	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	// Channel to listen for interrupt signal to terminate server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close Redis connection
	if err := redisService.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Server exited")
}
