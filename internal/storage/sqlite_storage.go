package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"github.com/PedroM2626/PriceWatcher/internal/models"
)

// SQLiteStorage implements Storage interface for SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dsn string) (*SQLiteStorage, error) {
	// Open database connection
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign key constraints
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

// createTables creates the necessary tables if they don't exist
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

// Close implements Storage.Close
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// CreateProduct implements Storage.CreateProduct
func (s *SQLiteStorage) CreateProduct(ctx context.Context, product *models.Product) error {
	product.ID = uuid.New()
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO products (id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, product.ID, product.Name, product.URL, product.ImageURL, product.CurrentPrice, product.Currency, 
	   product.IsAvailable, product.Website, product.CreatedAt, product.UpdatedAt)

	return err
}

// GetProductByID implements Storage.GetProductByID
func (s *SQLiteStorage) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at
		FROM products
		WHERE id = ?
	`, id).Scan(
		&product.ID, &product.Name, &product.URL, &product.ImageURL, &product.CurrentPrice,
		&product.Currency, &product.IsAvailable, &product.Website, &product.CreatedAt, &product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Implement other required methods with SQLite-specific queries
// ...

// Note: The rest of the methods (GetProductByURL, UpdateProduct, ListProducts, etc.)
// should be implemented following the same pattern as above, with SQLite-compatible SQL.
