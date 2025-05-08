package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the common logging methods
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	WithField(key string, value interface{}) Logger
}

// ZapLogger implements the Logger interface using zap
type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// Debug logs a message at debug level
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Info logs a message at info level
func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warn logs a message at warn level
func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Error logs a message at error level
func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// WithField returns a logger with the added field
func (l *ZapLogger) WithField(key string, value interface{}) Logger {
	newLogger := WithField(l.logger, key, value)
	return &ZapLogger{
		logger: newLogger,
		sugar:  newLogger.Sugar(),
	}
}

// Config holds logger configuration
type Config struct {
	Environment string
	LogLevel    string
	Encoding    string
	ServiceName string
	Development bool
}

// DefaultConfig returns sensible default configuration
func DefaultConfig() Config {
	return Config{
		Environment: "development",
		LogLevel:    "info",
		Encoding:    "json",
		ServiceName: "service",
		Development: false,
	}
}

// New creates a new logger with the given configuration
func New(cfg Config) (Logger, error) {
	// Parse log level
	level, err := getLogLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	// Determine if we're in production
	isProduction := cfg.Environment == "production"

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configure global logger options
	zapOptions := []zap.Option{
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("service", cfg.ServiceName),
			zap.String("env", cfg.Environment),
		),
	}

	// Add stacktrace on error and above in production
	if isProduction {
		zapOptions = append(zapOptions, zap.AddStacktrace(zap.ErrorLevel))
	} else if cfg.Development {
		// In development, include more information
		zapOptions = append(zapOptions, zap.AddCaller(), zap.Development())
	}

	// Configure output format
	var encoder zapcore.Encoder
	if cfg.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger
	zapLogger := zap.New(core, zapOptions...)
	sugarLogger := zapLogger.Sugar()

	// Create and return our custom logger
	return &ZapLogger{
		logger: zapLogger,
		sugar:  sugarLogger,
	}, nil
}

// NewZapLogger creates a new zap logger with the given configuration
// This is for backward compatibility
func NewZapLogger(cfg Config) (*zap.Logger, error) {
	logger, err := New(cfg)
	if err != nil {
		return nil, err
	}
	
	// Type assert to get the underlying zap logger
	zapLogger, ok := logger.(*ZapLogger)
	if !ok {
		return nil, fmt.Errorf("failed to convert to ZapLogger")
	}
	
	return zapLogger.logger, nil
}

// getLogLevel converts string log level to zapcore.Level
func getLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "dpanic":
		return zapcore.DPanicLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

// WithField adds a field to the logger
func WithField(logger *zap.Logger, key string, value interface{}) *zap.Logger {
	return logger.With(zap.Any(key, value))
}

// WithFields adds multiple fields to the logger
func WithFields(logger *zap.Logger, fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return logger.With(zapFields...)
}