package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	// Ensure the directory exists
	dbPath := "./data/pricewatcher.db"
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	// Open SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	fmt.Printf("Successfully connected to SQLite database at %s\n", dbPath)
}

func createTables(db *sql.DB) error {
	// Create products table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			url TEXT NOT NULL UNIQUE,
			image_url TEXT,
			current_price REAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'BRL',
			is_available BOOLEAN NOT NULL DEFAULT 1,
			website TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	// Create price_history table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS price_history (
			id TEXT PRIMARY KEY,
			product_id TEXT NOT NULL,
			price REAL NOT NULL,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create price_history table: %w", err)
	}

	// Create alerts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS alerts (
			id TEXT PRIMARY KEY,
			product_id TEXT NOT NULL,
			target_price REAL NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT 1,
			notification_type TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			notified_at TIMESTAMP,
			CHECK (notification_type IN ('email', 'telegram')),
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create alerts table: %w", err)
	}

	// Create indexes
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_price_history_product_id ON price_history(product_id);
		CREATE INDEX IF NOT EXISTS idx_price_history_created_at ON price_history(created_at);
		CREATE INDEX IF NOT EXISTS idx_alerts_product_id ON alerts(product_id);
		CREATE INDEX IF NOT EXISTS idx_alerts_is_active ON alerts(is_active);
	`)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
