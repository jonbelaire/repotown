package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
)

// RequestLogger logs information about each HTTP request
func RequestLogger(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			
			next.ServeHTTP(ww, r)
			
			logger.Info("HTTP request",
				"status", ww.Status(),
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}

// Recoverer is a middleware that recovers from panics
func Recoverer(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					logger.Error("Recovered from panic",
						"panic_value", rvr,
						"request_id", middleware.GetReqID(r.Context()),
					)
					
					httputils.ErrorJSON(w, httputils.ErrInternal)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// CORSConfig defines configuration for CORS middleware
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns sensible default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
}

// CORS returns a CORS middleware with the provided configuration
func CORS(cfg CORSConfig) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		ExposedHeaders:   cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge,
	})
}

// HealthCheck is a simple middleware that adds a health check endpoint
func HealthCheck(path string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == path {
				httputils.JSON(w, http.StatusOK, map[string]string{
					"status": "ok",
					"time":   time.Now().Format(time.RFC3339),
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}