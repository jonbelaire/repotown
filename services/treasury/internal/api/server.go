package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jonbelaire/repotown/packages/go-core/database"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/packages/go-core/server"
	"github.com/jonbelaire/repotown/services/treasury/internal/config"
	"github.com/jonbelaire/repotown/services/treasury/internal/repository"
	"github.com/jonbelaire/repotown/services/treasury/internal/service"
)

// Server represents the API server
type Server struct {
	*server.Server
	config       config.Config
	router       chi.Router
	logger       logging.Logger
	db           *database.DB
	repositories *repository.Repositories
	services     *service.Services
}

// NewServer creates a new API server
func NewServer(cfg config.Config, logger logging.Logger) (*Server, error) {
	// Connect to database
	db, err := repository.Connect(cfg.DatabaseConfig(), logger)
	if err != nil {
		return nil, err
	}

	// Create repositories
	repos := repository.NewRepositories(db, logger)

	// Create services
	services := service.NewServices(repos, logger)

	// Create server
	serverCfg := server.Config{
		Address:         cfg.ServerAddress,
		ShutdownTimeout: cfg.ShutdownTimeout,
	}

	// Create router
	r := chi.NewRouter()

	// Create server instance
	s := &Server{
		config:       cfg,
		router:       r,
		logger:       logger,
		db:           db,
		repositories: repos,
		services:     services,
	}

	// Create server with routes
	s.Server = server.New(serverCfg, logger,
		server.WithRouter(r),
		server.WithRoutes(s.setupRoutes),
	)

	return s, nil
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes(r chi.Router) {
	// Basic middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Create handlers
		taxRateHandler := NewTaxRateHandler(s.services.TaxRate, s.logger)
		taxpayerHandler := NewTaxpayerHandler(s.services.Taxpayer, s.logger)
		taxFilingHandler := NewTaxFilingHandler(s.services.TaxFiling, s.logger)
		taxPaymentHandler := NewTaxPaymentHandler(s.services.TaxPayment, s.logger)
		taxReportHandler := NewTaxReportHandler(s.services.TaxReport, s.logger)

		// Register routes
		taxRateHandler.RegisterRoutes(r)
		taxpayerHandler.RegisterRoutes(r)
		taxFilingHandler.RegisterRoutes(r)
		taxPaymentHandler.RegisterRoutes(r)
		taxReportHandler.RegisterRoutes(r)
	})
}

// Shutdown closes all resources gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	// Close database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			s.logger.Error("Failed to close database connection", "error", err)
		}
	}

	// Shutdown server
	return s.Server.Shutdown(ctx)
}