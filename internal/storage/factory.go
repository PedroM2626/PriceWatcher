package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/PedroM2626/PriceWatcher/internal/config"
)

// NewStorage creates a new storage instance based on the provided configuration
func NewStorage(cfg config.DatabaseConfig) (Storage, error) {
	switch cfg.Driver {
	case "postgres":
		// For PostgreSQL, build the connection string from individual parameters
		if cfg.DSN == "" {
			dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
			return NewPostgresStorage(dsn)
		}
		return NewPostgresStorage(cfg.DSN)
		
	case "sqlite3", "sqlite":
		if cfg.DSN == "" {
			return nil, fmt.Errorf("DSN is required for SQLite")
		}
		
		// Ensure the directory exists
		dir := filepath.Dir(cfg.DSN)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory for SQLite database: %w", err)
			}
		}
		
		// Add error logging for SQLite connection
		log.Printf("Connecting to SQLite database at: %s", cfg.DSN)
		storage, err := NewSQLiteStorage(cfg.DSN)
		if err != nil {
			log.Printf("Failed to initialize SQLite storage: %v", err)
			return nil, fmt.Errorf("failed to initialize SQLite storage: %w", err)
		}
		return storage, nil
		
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

// NewPostgresStorage creates a new PostgreSQL storage instance with a connection string
func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresStorage{db: db}, nil
}
