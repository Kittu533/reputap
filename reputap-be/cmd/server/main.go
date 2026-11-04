package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/kittu533/reputap-be/internal/activation"
	"github.com/kittu533/reputap-be/internal/analytics"
	"github.com/kittu533/reputap-be/internal/audit"
	"github.com/kittu533/reputap-be/internal/auth"
	"github.com/kittu533/reputap-be/internal/config"
	"github.com/kittu533/reputap-be/internal/customers"
	"github.com/kittu533/reputap-be/internal/dashboard"
	"github.com/kittu533/reputap-be/internal/database"
	"github.com/kittu533/reputap-be/internal/design"
	"github.com/kittu533/reputap-be/internal/middleware"
	"github.com/kittu533/reputap-be/internal/orders"
	"github.com/kittu533/reputap-be/internal/productinstances"
	"github.com/kittu533/reputap-be/internal/products"
	"github.com/kittu533/reputap-be/internal/production"
	"github.com/kittu533/reputap-be/internal/qc"
	"github.com/kittu533/reputap-be/internal/qrcodes"
	"github.com/kittu533/reputap-be/internal/redirect"
	"github.com/kittu533/reputap-be/internal/response"
	"github.com/kittu533/reputap-be/internal/shipping"
	"github.com/kittu533/reputap-be/internal/stores"

	_ "github.com/kittu533/reputap-be/docs"
)

//	@title						REPUTAP API
//	@version					1.0
//	@description				Internal API for REPUTAP — Smart Review Stand management system.
//	@description				Manages customers, stores, orders, QR codes, NFC, production, QC, shipping, and scan analytics.
//
//	@contact.name				REPUTAP Team
//	@contact.email				admin@reputap.id
//
//	@host						localhost:8081
//	@BasePath					/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your JWT token with the `Bearer ` prefix. Example: `Bearer eyJhbGci...`
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()
	slog.Info("configuration loaded", "env", cfg.AppEnv, "port", cfg.AppPort)

	if cfg.JWTSecret == "" || cfg.JWTSecret == "your-256-bit-secret-change-this-in-production" {
		slog.Warn("JWT_SECRET is not set or using default value. Please set a secure secret for production.")
	}

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// ========================================
	// Initialize modules
	// ========================================

	// Auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService)

	// Customers
	customerRepo := customers.NewRepository(db)
	customerService := customers.NewService(customerRepo)
	customerHandler := customers.NewHandler(customerService)

	// Stores
	storeRepo := stores.NewRepository(db)
	storeService := stores.NewService(storeRepo)
	storeHandler := stores.NewHandler(storeService)

	// Products
	productRepo := products.NewRepository(db)
	productService := products.NewService(productRepo)
	productHandler := products.NewHandler(productService)

	// Orders
	orderRepo := orders.NewRepository(db)
	orderService := orders.NewService(orderRepo)
	orderHandler := orders.NewHandler(orderService)

	// Product Instances
	instanceRepo := productinstances.NewRepository(db)
	instanceService := productinstances.NewService(instanceRepo)
	instanceHandler := productinstances.NewHandler(instanceService)

	// QR Codes
	qrRepo := qrcodes.NewRepository(db)
	qrService := qrcodes.NewService(qrRepo, cfg.RedirectBaseURL)
	qrHandler := qrcodes.NewHandler(qrService)

	// Redirect Engine
	redirectRepo := redirect.NewRepository(db)
	redirectService := redirect.NewService(redirectRepo)
	redirectHandler := redirect.NewHandler(redirectService, cfg.FrontendURL)

	// Activation
	activationRepo := activation.NewRepository(db)
	activationService := activation.NewService(activationRepo, cfg.RedirectBaseURL)
	activationHandler := activation.NewHandler(activationService)

	// Design Assets
	designRepo := design.NewRepository(db)
	designService := design.NewService(designRepo, cfg.UploadDir, cfg.MaxUploadSize)
	designHandler := design.NewHandler(designService, cfg.MaxUploadSize)

	// Production
	productionRepo := production.NewRepository(db)
	productionService := production.NewService(productionRepo)
	productionHandler := production.NewHandler(productionService)

	// Quality Control
	qcRepo := qc.NewRepository(db)
	qcService := qc.NewService(qcRepo)
	qcHandler := qc.NewHandler(qcService)

	// Shipping
	shippingRepo := shipping.NewRepository(db)
	shippingService := shipping.NewService(shippingRepo)
	shippingHandler := shipping.NewHandler(shippingService)

	// Audit Logs
	auditRepo := audit.NewRepository(db)
	auditService := audit.NewService(auditRepo)
	auditHandler := audit.NewHandler(auditService)

	// Dashboard
	dashboardRepo := dashboard.NewRepository(db)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	// Analytics
	analyticsRepo := analytics.NewRepository(db)
	analyticsService := analytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)

	// ========================================
	// Setup router
	// ========================================
	r := chi.NewRouter()

	// Global middleware
	r.Use(cors.Handler(middleware.CORS()))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Public redirect engine (no auth)
	r.Get("/r/{code}", redirectHandler.Redirect)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Get("/activate/check", activationHandler.Check)
		r.Post("/activate", activationHandler.Activate)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware(authService))

			// Auth
			r.Get("/auth/me", authHandler.Me)
			r.Post("/auth/logout", authHandler.Logout)

			// Customers
			r.Route("/customers", func(r chi.Router) {
				r.Get("/", customerHandler.List)
				r.Post("/", customerHandler.Create)
				r.Get("/{id}", customerHandler.GetByID)
				r.Put("/{id}", customerHandler.Update)
				r.Delete("/{id}", customerHandler.Delete)
			})

			// Stores
			r.Route("/stores", func(r chi.Router) {
				r.Get("/", storeHandler.List)
				r.Post("/", storeHandler.Create)
				r.Get("/{id}", storeHandler.GetByID)
				r.Put("/{id}", storeHandler.Update)
				r.Patch("/{id}/status", storeHandler.UpdateStatus)
				r.Put("/{id}/review-url", storeHandler.UpdateReviewURL)
			})

			// Products
			r.Route("/products", func(r chi.Router) {
				r.Get("/", productHandler.List)
				r.Post("/", productHandler.Create)
				r.Get("/{id}", productHandler.GetByID)
				r.Put("/{id}", productHandler.Update)
			})

			// Orders
			r.Route("/orders", func(r chi.Router) {
				r.Get("/", orderHandler.List)
				r.Post("/", orderHandler.Create)
				r.Get("/{id}", orderHandler.GetByID)
				r.Patch("/{id}/status", orderHandler.UpdateStatus)
				r.Patch("/{id}/payment-status", orderHandler.UpdatePaymentStatus)
				r.Get("/{id}/history", orderHandler.GetStatusHistory)
				r.Post("/{orderId}/instances/generate", instanceHandler.GenerateForOrder)
				r.Get("/{orderId}/instances", instanceHandler.ListByOrderID)
				r.Get("/{orderId}/designs", designHandler.ListByOrderID)
				r.Get("/{orderId}/shipment", shippingHandler.GetByOrderID)
			})

			// Product Instances
			r.Route("/product-instances", func(r chi.Router) {
				r.Get("/", instanceHandler.List)
				r.Get("/{id}", instanceHandler.GetByID)
				r.Patch("/{id}/nfc-status", instanceHandler.UpdateNFCStatus)
			})

			// QR Codes
			r.Route("/qr-codes", func(r chi.Router) {
				r.Get("/", qrHandler.List)
				r.Get("/{id}", qrHandler.GetByID)
				r.Patch("/{id}/status", qrHandler.UpdateStatus)
				r.Get("/{id}/download", qrHandler.Download)
			})

			// Design Assets
			r.Route("/design-assets", func(r chi.Router) {
				r.Get("/", designHandler.List)
				r.Post("/upload", designHandler.Upload)
				r.Get("/{id}", designHandler.GetByID)
				r.Patch("/{id}/approval", designHandler.UpdateApprovalStatus)
				r.Delete("/{id}", designHandler.Delete)
				r.Get("/{id}/download", designHandler.Download)
			})

			// Production Jobs
			r.Route("/production/jobs", func(r chi.Router) {
				r.Get("/", productionHandler.List)
				r.Post("/", productionHandler.Create)
				r.Get("/{id}", productionHandler.GetByID)
				r.Patch("/{id}/status", productionHandler.UpdateStatus)
			})

			// Quality Control
			r.Route("/quality-controls", func(r chi.Router) {
				r.Get("/", qcHandler.List)
				r.Post("/", qcHandler.Create)
				r.Get("/{id}", qcHandler.GetByID)
				r.Get("/instance/{instanceId}", qcHandler.ListByInstanceID)
			})

			// Shipments
			r.Route("/shipments", func(r chi.Router) {
				r.Get("/", shippingHandler.List)
				r.Post("/", shippingHandler.Create)
				r.Get("/{id}", shippingHandler.GetByID)
				r.Patch("/{id}/status", shippingHandler.UpdateStatus)
			})

			// Audit Logs
			r.Get("/audit-logs", auditHandler.List)

			// Dashboard
			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/summary", dashboardHandler.Summary)
				r.Get("/recent-activities", dashboardHandler.RecentActivities)
			})

			// Analytics
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/overview", analyticsHandler.Overview)
				r.Get("/scans", analyticsHandler.Scans)
				r.Get("/stores/{id}", analyticsHandler.StoreAnalytics)
			})
		})
	})

	// Start server
	addr := ":" + cfg.AppPort
	slog.Info("server starting", "addr", addr, "swagger", "http://localhost:"+cfg.AppPort+"/swagger/index.html")
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
