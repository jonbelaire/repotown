package database

import (
	"context"
	"database/sql"
	"time"

	// Import PostgreSQL driver
	_ "github.com/lib/pq"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
)

// Config holds database configuration
type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns sensible default configuration
func DefaultConfig() Config {
	return Config{
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

// DB wraps sql.DB with additional functionality
type DB struct {
	*sql.DB
	Logger logging.Logger
	Config Config
}

// New creates a new database connection
func New(cfg Config, logger logging.Logger) (*DB, error) {
	// Open connection
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Connected to database", "dsn", cfg.DSN)

	return &DB{
		DB:     db,
		Logger: logger,
		Config: cfg,
	}, nil
}

// Transaction represents a database transaction
type Transaction struct {
	*sql.Tx
	Committed bool
	Rolled    bool
}

// WithTransaction runs the given function within a transaction
func (db *DB) WithTransaction(ctx context.Context, fn func(*Transaction) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	t := &Transaction{
		Tx:       tx,
		Committed: false,
		Rolled:    false,
	}

	defer func() {
		if !t.Committed && !t.Rolled {
			if err := t.Rollback(); err != nil {
				db.Logger.Error("Failed to rollback transaction", "error", err)
			}
		}
	}()

	if err := fn(t); err != nil {
		if rollbackErr := t.Rollback(); rollbackErr != nil {
			db.Logger.Error("Failed to rollback transaction", "error", rollbackErr)
		}
		return err
	}

	return t.Commit()
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	if t.Committed || t.Rolled {
		return nil
	}
	
	t.Committed = true
	return t.Tx.Commit()
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	if t.Committed || t.Rolled {
		return nil
	}
	
	t.Rolled = true
	return t.Tx.Rollback()
}

// Close closes the database connection
func (db *DB) Close() error {
	db.Logger.Info("Closing database connection")
	return db.DB.Close()
}