package config

import (
	"time"
)

// Config holds application configuration
type Config struct {
	// Server settings
	ServerAddress   string        `envconfig:"SERVER_ADDRESS" default:":8080"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"15s"`

	// Environment settings
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`

	// Database settings
	DatabaseURL          string        `envconfig:"DATABASE_URL" required:"true"`
	DatabaseMaxOpenConns int           `envconfig:"DATABASE_MAX_OPEN_CONNS" default:"25"`
	DatabaseMaxIdleConns int           `envconfig:"DATABASE_MAX_IDLE_CONNS" default:"25"`
	DatabaseMaxLifetime  time.Duration `envconfig:"DATABASE_MAX_LIFETIME" default:"5m"`
	DatabaseMaxIdleTime  time.Duration `envconfig:"DATABASE_MAX_IDLE_TIME" default:"5m"`

	// Auth settings
	JWTSecret string        `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiry time.Duration `envconfig:"JWT_EXPIRY" default:"24h"`
	JWTIssuer string        `envconfig:"JWT_ISSUER" default:"banking-service"`
}

// DatabaseConfig returns the database configuration
func (c *Config) DatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		URL:          c.DatabaseURL,
		MaxOpenConns: c.DatabaseMaxOpenConns,
		MaxIdleConns: c.DatabaseMaxIdleConns,
		MaxLifetime:  c.DatabaseMaxLifetime,
		MaxIdleTime:  c.DatabaseMaxIdleTime,
	}
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	URL          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
}
