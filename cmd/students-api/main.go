package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rishavvrajj/golang-rest-api/internal/config"
	"github.com/rishavvrajj/golang-rest-api/internal/http/router"
	"github.com/rishavvrajj/golang-rest-api/internal/service"
	"github.com/rishavvrajj/golang-rest-api/internal/storage/sqlite"
)

func main() {
	// 1. Load configuration (YAML / Env vars)
	cfg := config.MustLoad()

	// 2. Setup structured logging (slog)
	var logger *slog.Logger
	if cfg.Env == "local" || cfg.Env == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(logger)

	slog.Info("starting students api server", slog.String("env", cfg.Env), slog.String("address", cfg.HttpServer.Address))

	// 3. Initialize storage repository (SQLite)
	dbStorage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		slog.Error("failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("storage initialized successfully", slog.String("path", cfg.StoragePath))

	// 4. Initialize business logic service layer
	studentService := service.NewStudentService(dbStorage)

	// 5. Setup HTTP router and middleware
	handler := router.NewRouter(studentService)

	// 6. Setup HTTP server with production-grade timeouts
	server := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Graceful shutdown handling
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("HTTP server is listening", slog.String("address", cfg.HttpServer.Address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down HTTP server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown of HTTP server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("HTTP server stopped successfully")
}
