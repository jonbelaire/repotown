# Go Core Packages

Shared Go utilities and packages for microservices in the Repotown monorepo.

## Core Packages

### Server

HTTP server setup with sensible defaults:

```go
import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/jonbelaire/repotown/packages/go-core/logging"
    "github.com/jonbelaire/repotown/packages/go-core/server"
)

// Create logger
logger, _ := logging.New(logging.DefaultConfig())

// Create server with custom routes and middleware
srv := server.New(server.DefaultConfig(), logger, 
    server.WithRoutes(func(r chi.Router) {
        r.Get("/api/orders", func(w http.ResponseWriter, r *http.Request) {
            // Handle request
        })
    }),
)

// Start server
srv.Start()

// Graceful shutdown
srv.Shutdown()
```

### Database

Database connection and transaction management:

```go
import (
    "context"
    "github.com/jonbelaire/repotown/packages/go-core/database"
    "github.com/jonbelaire/repotown/packages/go-core/logging"
)

// Create logger and database connection
logger, _ := logging.New(logging.DefaultConfig())
db, _ := database.New(database.DefaultConfig(), logger)
defer db.Close()

// Use transaction with automatic commit/rollback
ctx := context.Background()
err := db.WithTransaction(ctx, func(tx *database.Transaction) error {
    // Perform database operations
    return nil // Will commit if no error
})
```

### Logging

Structured logging with zap:

```go
import "github.com/jonbelaire/repotown/packages/go-core/logging"

// Create logger with custom configuration
cfg := logging.DefaultConfig()
cfg.ServiceName = "order-service"
logger, _ := logging.New(cfg)

// Use structured logging
logger.Info("Application started", "version", "1.0.0")

// Add context fields
orderLogger := logging.WithField(logger, "order_id", "ORD-12345")
orderLogger.Info("Processing order")
```

### Config

Environment-based configuration:

```go
import "github.com/jonbelaire/repotown/packages/go-core/config"

// Define your configuration structure
type AppConfig struct {
    Server struct {
        Port int `envconfig:"SERVER_PORT" default:"8080"`
    }
    Database struct {
        Host string `envconfig:"DB_HOST" required:"true"`
    }
}

// Load configuration from environment
var cfg AppConfig
options := config.DefaultLoadOptions()
options.EnvFile = ".env"
config.LoadConfig(&cfg, options)
```

### Middleware

Common HTTP middleware components:

```go
import (
    "github.com/jonbelaire/repotown/packages/go-core/middleware"
    "go.uber.org/zap"
)

// Create router
r := chi.NewRouter()

// Add middleware
r.Use(middleware.RequestLogger(logger.Sugar()))
r.Use(middleware.Recoverer(logger))
r.Use(middleware.CORS(middleware.DefaultCORSConfig()))
r.Use(middleware.HealthCheck("/health"))
```

### Auth

JWT authentication utilities:

```go
import "github.com/jonbelaire/repotown/packages/go-core/auth"

// Create a JWT service
jwtService := auth.NewJWTService(
    "your-secret-key",
    24*time.Hour,
    "your-service-name",
)

// Generate a token
token, err := jwtService.GenerateToken("user123", "admin")

// Validate a token
claims, err := jwtService.ValidateToken(tokenString)
if err != nil {
    // Handle error
}
```

### HTTP Utils

Standard HTTP response formatting:

```go
import "github.com/jonbelaire/repotown/packages/go-core/httputils"

// Success response
httputils.JSON(w, http.StatusOK, data)

// Error response
err := httputils.NewError(
    http.StatusBadRequest,
    "INVALID_INPUT",
    "The input provided is invalid",
    details,
)
httputils.ErrorJSON(w, err)

// Common errors
httputils.ErrorJSON(w, httputils.ErrNotFound)
```

### Validator

Request validation with detailed error messages:

```go
import "github.com/jonbelaire/repotown/packages/go-core/validator"

type UserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

// Create a validator
v := validator.NewValidator()

// Validate a request
valid, errResp := v.ValidateJSON(request)
if !valid {
    httputils.ErrorJSON(w, *errResp)
    return
}
```

## Package Design Principles

1. **Flexibility**: Components can be used independently or together
2. **Sensible Defaults**: Provides reasonable defaults while allowing customization
3. **Pluggable**: Easy to replace components with custom implementations
4. **Production-Ready**: Includes features needed for production services
5. **Minimal Dependencies**: Only includes essential dependencies

## Usage in Microservices

To use these packages in your microservices, add the following to your `go.mod`:

```
require github.com/jonbelaire/repotown/packages/go-core v0.0.0
replace github.com/jonbelaire/repotown/packages/go-core => ../../packages/go-core
```

This ensures your services use the local version of these packages in the monorepo.