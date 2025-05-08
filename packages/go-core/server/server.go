package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
)

// ServerOption allows customizing the server
type ServerOption func(*Server)

// Server represents a generic HTTP server with common functionalities
type Server struct {
	*http.Server
	Router          *chi.Mux
	Logger          logging.Logger
	ShutdownTimeout time.Duration
	MiddlewareHooks []func(r chi.Router)
	RouteHooks      []func(r chi.Router)
}

// Config holds server configuration
type Config struct {
	Address             string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	ShutdownTimeout     time.Duration
	ReadTimeoutSecs     int
	WriteTimeoutSecs    int
	IdleTimeoutSecs     int
	ShutdownTimeoutSecs int
}

// DefaultConfig returns sensible default configuration
func DefaultConfig() Config {
	return Config{
		Address:      ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}
}

// New creates a new server with the given options
func New(cfg Config, logger logging.Logger, options ...ServerOption) *Server {
	r := chi.NewRouter()
	
	s := &Server{
		Server: &http.Server{
			Addr:         cfg.Address,
			Handler:      r,
			ReadTimeout:  getDuration(cfg.ReadTimeout, cfg.ReadTimeoutSecs),
			WriteTimeout: getDuration(cfg.WriteTimeout, cfg.WriteTimeoutSecs),
			IdleTimeout:  getDuration(cfg.IdleTimeout, cfg.IdleTimeoutSecs),
		},
		Router:         r,
		Logger:         logger,
		ShutdownTimeout: getDuration(cfg.ShutdownTimeout, cfg.ShutdownTimeoutSecs),
		MiddlewareHooks: []func(r chi.Router){},
		RouteHooks:      []func(r chi.Router){},
	}
	
	// Apply options
	for _, option := range options {
		option(s)
	}
	
	// Setup default middleware
	s.setupDefaultMiddleware()
	
	// Apply middleware hooks
	for _, hook := range s.MiddlewareHooks {
		hook(s.Router)
	}
	
	// Apply route hooks
	for _, hook := range s.RouteHooks {
		hook(s.Router)
	}
	
	return s
}

// getDuration returns the duration from either a set duration or seconds
func getDuration(duration time.Duration, seconds int) time.Duration {
	if duration != 0 {
		return duration
	}
	if seconds != 0 {
		return time.Duration(seconds) * time.Second
	}
	return 0
}

// setupDefaultMiddleware adds basic middleware that most services need
func (s *Server) setupDefaultMiddleware() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))
	s.Router.Use(requestLogger(s.Logger))
}

// requestLogger is a middleware that logs HTTP requests
func requestLogger(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			next.ServeHTTP(ww, r)
			
			logger.Info("HTTP request",
				"status", ww.Status(),
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}

// WithMiddleware adds custom middleware to the server
func WithMiddleware(middlewareFunc func(chi.Router)) ServerOption {
	return func(s *Server) {
		s.MiddlewareHooks = append(s.MiddlewareHooks, middlewareFunc)
	}
}

// WithRouter adds a router to the server
func WithRouter(router chi.Router) ServerOption {
	return func(s *Server) {
		s.Router = router.(*chi.Mux)
	}
}

// WithRoutes adds custom routes to the server
func WithRoutes(routeFunc func(chi.Router)) ServerOption {
	return func(s *Server) {
		s.RouteHooks = append(s.RouteHooks, routeFunc)
	}
}

// Start starts the server
func (s *Server) Start() error {
	s.Logger.Info("Starting server", "address", s.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("Server error", "error", err)
		}
	}()
	return nil
}

// ListenAndServe starts the server and blocks until it shuts down
func (s *Server) ListenAndServe() error {
	return s.Server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Shutting down server")
	
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()
	}
	
	return s.Server.Shutdown(ctx)
}

// Close immediately closes the server
func (s *Server) Close() error {
	s.Logger.Info("Closing server")
	return s.Server.Close()
}