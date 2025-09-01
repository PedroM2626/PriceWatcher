package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/PedroM2626/PriceWatcher/internal/config"
)

// NewStorage creates a new storage instance based on the provided configuration
func NewStorage(cfg config.DatabaseConfig) (Storage, error) {
	switch cfg.Driver {
	case "sqlite3", "sqlite":
		if cfg.DSN == "" {
			return nil, fmt.Errorf("DSN is required for SQLite")
		}
		// Ensure the directory exists
		dir := filepath.Dir(cfg.DSN)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, fmt.Errorf("failed to create directory for SQLite database: %w", err)
			}
		}
		log.Printf("Connecting to SQLite database at: %s", cfg.DSN)
		return NewSQLiteStorage(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}
