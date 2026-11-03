package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/kittu533/reputap-be/internal/auth"
	"github.com/kittu533/reputap-be/internal/config"
	"github.com/kittu533/reputap-be/internal/database"
	"github.com/kittu533/reputap-be/internal/middleware"
	"github.com/kittu533/reputap-be/internal/response"
)

// 💡 Go Concept: main()
// Every Go program starts at main(). This is the entry point.
// There's no framework "bootstrapping" — you wire everything yourself.
// This gives you full control and makes the code explicit.
func main() {
	// Setup structured logging
	// 💡 Go Concept: slog (structured logging)
	// slog is Go's built-in structured logger (since Go 1.21).
	// It outputs JSON logs which are great for production (parseable by log tools).
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()
	slog.Info("configuration loaded", "env", cfg.AppEnv, "port", cfg.AppPort)

	// Validate critical config
	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-256-bit-secret-change-this-in-production" {
		slog.Warn("JWT_SECRET is not set or using default value. Please set a secure secret for production.")
	}

	// Connect to database
	db, err := database.Connect(cfg.DSN())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// ========================================
	// Initialize modules
	// 💡 Go Concept: Manual dependency injection
	// We create each layer manually and pass dependencies.
	// No magic DI container — everything is explicit.
	// This makes the code easy to trace: "where does this come from?"
	// ========================================

	// Auth module
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService)

	// ========================================
	// Setup router
	// 💡 Go Concept: Chi Router
	// Chi is a lightweight HTTP router for Go.
	// It's composable: you can nest routers, apply middleware per-route.
	// r.Route("/path", func(r chi.Router) { ... }) creates a sub-router.
	// r.Use(middleware) applies middleware to all routes in that group.
	// r.Get/Post/Put/Patch/Delete register routes.
	// ========================================
	r := chi.NewRouter()

	// Global middleware (applied to ALL routes)
	r.Use(cors.Handler(middleware.CORS()))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)

	// Health check (no auth needed)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Public redirect endpoint (no auth needed)
	// Will be registered when redirect module is built
	// r.Get("/r/{code}", redirectHandler.Redirect)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// --- Public routes (no auth) ---
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)

		// --- Protected routes (auth required) ---
		// 💡 Go Concept: Route groups with middleware
		// r.Group creates a new group that inherits the parent's middleware
		// and adds its own. Here, all routes inside require JWT auth.
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware(authService))

			// Auth
			r.Get("/auth/me", authHandler.Me)
			r.Post("/auth/logout", authHandler.Logout)

			// Module routes will be registered here as they're built:
			// r.Route("/customers", func(r chi.Router) { ... })
			// r.Route("/stores", func(r chi.Router) { ... })
			// r.Route("/products", func(r chi.Router) { ... })
			// r.Route("/orders", func(r chi.Router) { ... })
			// etc.
		})
	})

	// Start server
	addr := ":" + cfg.AppPort
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
