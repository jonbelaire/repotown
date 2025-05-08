package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadOptions defines options for loading configuration
type LoadOptions struct {
	EnvFile       string
	EnvPrefix     string
	RequiredEnvs  []string
	DefaultValues map[string]interface{}
}

// DefaultLoadOptions returns sensible default loading options
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		EnvFile:       ".env",
		EnvPrefix:     "",
		RequiredEnvs:  []string{},
		DefaultValues: make(map[string]interface{}),
	}
}

// LoadConfig loads configuration from environment variables into the given struct
func LoadConfig(cfg interface{}, opts LoadOptions) error {
	// Load .env file if it exists
	if opts.EnvFile != "" {
		if _, err := os.Stat(opts.EnvFile); err == nil {
			if err := godotenv.Load(opts.EnvFile); err != nil {
				return fmt.Errorf("error loading env file: %w", err)
			}
		}

		// Try loading environment-specific .env files
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = os.Getenv("ENV")
		}
		if env == "" {
			env = os.Getenv("ENVIRONMENT")
		}
		if env == "" {
			env = "development"
		}

		envFile := fmt.Sprintf("%s.%s", strings.TrimSuffix(opts.EnvFile, filepath.Ext(opts.EnvFile)), env)
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return fmt.Errorf("error loading env file %s: %w", envFile, err)
			}
		}
	}

	// Apply default values
	if len(opts.DefaultValues) > 0 {
		applyDefaultValues(cfg, opts.DefaultValues)
	}

	// Process environment variables
	if err := envconfig.Process(opts.EnvPrefix, cfg); err != nil {
		return fmt.Errorf("error processing env vars: %w", err)
	}

	// Check for required environment variables
	for _, env := range opts.RequiredEnvs {
		if os.Getenv(env) == "" {
			return fmt.Errorf("required environment variable %s is not set", env)
		}
	}

	return nil
}

// applyDefaultValues sets default values for the configuration struct
func applyDefaultValues(cfg interface{}, defaults map[string]interface{}) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return
	}

	for key, val := range defaults {
		setStructField(v.Elem(), key, val)
	}
}

// setStructField sets a field value in a struct
func setStructField(v reflect.Value, name string, value interface{}) {
	field := v.FieldByName(name)
	if !field.IsValid() {
		return
	}

	if field.CanSet() {
		fieldVal := reflect.ValueOf(value)
		if fieldVal.Type().AssignableTo(field.Type()) {
			field.Set(fieldVal)
		}
	}
}