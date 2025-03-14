package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ladderseeker/gin-crud-starter/configs"
	"github.com/ladderseeker/gin-crud-starter/internal/api/routes"
	"github.com/ladderseeker/gin-crud-starter/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *configs.Config
	db     *gorm.DB
}

// NewServer creates a new server instance
func NewServer(config *configs.Config, db *gorm.DB) *Server {
	// Set Gin mode
	gin.SetMode(config.Server.Mode)

	// Create router
	router := gin.New()

	return &Server{
		router: router,
		config: config,
		db:     db,
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Setup routes
	routes.SetupRoutes(s.router, s.db)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + s.config.Server.Port,
		Handler:      s.router,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", s.config.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error starting server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Server shutting down...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
		return err
	}

	logger.Info("Server exited gracefully")
	return nil
}
