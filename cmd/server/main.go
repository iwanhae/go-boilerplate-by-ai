package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"gosuda.org/boilerplate/api"
	"gosuda.org/boilerplate/internal/application"
	"gosuda.org/boilerplate/internal/config"
	"gosuda.org/boilerplate/internal/infrastructure"
	"gosuda.org/boilerplate/internal/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := infrastructure.NewLogger(&cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize storage
	baseStore := infrastructure.NewMemoryStore()
	metrics := infrastructure.NewMetricsCollector()
	store := infrastructure.NewMetricsStore(baseStore, metrics)

	// Initialize services
	postService := application.NewPostService(store, metrics)
	debugService := application.NewDebugService(logger, store, metrics)

	// Initialize middleware
	requestIDMiddleware := middleware.NewRequestIDMiddleware()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	recoveryMiddleware := middleware.NewRecoveryMiddleware(logger)
	corsMiddleware := middleware.NewCORSMiddleware(&cfg.CORS)
	errorHandlerMiddleware := middleware.NewErrorHandlerMiddleware(logger)
	metricsMiddleware := middleware.NewMetricsMiddleware(metrics)

	// Initialize handlers
	handlers := api.NewHandlers(postService, debugService, errorHandlerMiddleware)

	// Create router
	r := chi.NewRouter()

	// Add middleware in order
	r.Use(recoveryMiddleware.Handler)
	r.Use(requestIDMiddleware.Handler)
	r.Use(loggingMiddleware.Handler)
	r.Use(corsMiddleware.Handler)
	r.Use(metricsMiddleware.Handler)

	// Add Chi middleware
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// API routes
	r.Route("/debug", func(r chi.Router) {
		r.Get("/metrics", handlers.GetMetrics)
		r.Post("/logs", handlers.SetLogLevel)
		
		// Pprof routes
		r.Route("/pprof", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/profile", pprof.Profile)
			r.Post("/symbol", pprof.Symbol)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
			r.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
			r.Get("/block", pprof.Handler("block").ServeHTTP)
			r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
			r.Get("/heap", pprof.Handler("heap").ServeHTTP)
			r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
			r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		})
	})

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", handlers.ListPosts)
		r.Post("/", handlers.CreatePost)
		r.Get("/{id}", handlers.GetPost)
		r.Put("/{id}", handlers.UpdatePost)
		r.Delete("/{id}", handlers.DeletePost)
	})

	// Health check
	r.Get("/health", handlers.GetHealth)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Log startup
	logger.LogStartup("1.0.0", "development", cfg)

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.LogShutdown("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	// Close storage
	if err := store.Close(); err != nil {
		logger.Error("Storage close error", "error", err)
	}

	logger.LogShutdown("Server stopped")
}