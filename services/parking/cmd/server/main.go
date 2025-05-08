package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonbelaire/repotown/packages/go-core/config"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/api"
	parkingConfig "github.com/jonbelaire/repotown/services/parking/internal/config"
)

func main() {
	// Load configuration
	cfg := parkingConfig.Config{}
	opts := config.DefaultLoadOptions()
	opts.EnvFile = ".env"
	if err := config.LoadConfig(&cfg, opts); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create logger
	logger, err := logging.New(logging.Config{
		ServiceName: "parking-service",
		Development: cfg.Environment != "production",
		LogLevel:    cfg.LogLevel,
	})
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		os.Exit(1)
	}

	// Create and start server
	server, err := api.NewServer(cfg, logger)
	if err != nil {
		logger.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}