package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Remove existing database file if it exists
	_ = os.Remove("./pricewatcher.db")

	// Open SQLite database
	db, err := sql.Open("sqlite3", "./pricewatcher.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Enable foreign key constraints
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	fmt.Println("Database setup completed successfully!")
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
		);

		CREATE TABLE IF NOT EXISTS price_history (
			id TEXT PRIMARY KEY,
			product_id TEXT NOT NULL,
			price REAL NOT NULL,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_price_history_product_id ON price_history(product_id);
		CREATE INDEX IF NOT EXISTS idx_price_history_created_at ON price_history(created_at);

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
		);

		CREATE INDEX IF NOT EXISTS idx_alerts_product_id ON alerts(product_id);
		CREATE INDEX IF NOT EXISTS idx_alerts_is_active ON alerts(is_active);
	`)

	return err
}
