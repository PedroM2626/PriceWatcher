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
	// Open database connection (driver name is "sqlite" for glebarez/go-sqlite)
	db, err := sql.Open("sqlite", dsn)
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
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			url TEXT NOT NULL UNIQUE,
			image_url TEXT,
			current_price REAL NOT NULL DEFAULT 0,
			currency TEXT NOT NULL DEFAULT 'BRL',
			is_available INTEGER NOT NULL DEFAULT 1,
			website TEXT NOT NULL DEFAULT '',
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
			is_active INTEGER NOT NULL DEFAULT 1,
			notification_type TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			notified_at TIMESTAMP,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_alerts_product_id ON alerts(product_id);
		CREATE INDEX IF NOT EXISTS idx_alerts_is_active ON alerts(is_active);
	`)
	return err
}

// Close implements Storage.Close
func (s *SQLiteStorage) Close() error { return s.db.Close() }

// CreateProduct implements Storage.CreateProduct
func (s *SQLiteStorage) CreateProduct(ctx context.Context, product *models.Product) error {
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}
	now := time.Now()
	if product.CreatedAt.IsZero() { product.CreatedAt = now }
	product.UpdatedAt = now

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO products (id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, product.ID.String(), product.Name, product.URL, product.ImageURL, product.CurrentPrice, product.Currency,
		boolToInt(product.IsAvailable), product.Website, product.CreatedAt, product.UpdatedAt)
	return err
}

// GetProductByID implements Storage.GetProductByID
func (s *SQLiteStorage) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var p models.Product
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at
		FROM products WHERE id = ?
	`, id.String()).Scan(&p.ID, &p.Name, &p.URL, &p.ImageURL, &p.CurrentPrice, &p.Currency, &p.IsAvailable, &p.Website, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }
	return &p, nil
}

// UpdateProduct implements Storage.UpdateProduct
func (s *SQLiteStorage) UpdateProduct(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()
	_, err := s.db.ExecContext(ctx, `
		UPDATE products SET name = ?, url = ?, image_url = ?, current_price = ?, currency = ?, is_available = ?, website = ?, updated_at = ?
		WHERE id = ?
	`, product.Name, product.URL, product.ImageURL, product.CurrentPrice, product.Currency, boolToInt(product.IsAvailable), product.Website, product.UpdatedAt, product.ID.String())
	return err
}

// ListProducts implements Storage.ListProducts
func (s *SQLiteStorage) ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	query := `SELECT id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at FROM products ORDER BY created_at DESC`
	args := []any{}
	if limit > 0 { query += " LIMIT ?"; args = append(args, limit) }
	if offset > 0 { query += " OFFSET ?"; args = append(args, offset) }

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.URL, &p.ImageURL, &p.CurrentPrice, &p.Currency, &p.IsAvailable, &p.Website, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, &p)
	}
	return items, rows.Err()
}

// DeleteProduct implements Storage.DeleteProduct
func (s *SQLiteStorage) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id.String())
	return err
}

// AddPriceHistory implements Storage.AddPriceHistory
func (s *SQLiteStorage) AddPriceHistory(ctx context.Context, productID uuid.UUID, price float64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO price_history (id, product_id, price, created_at)
		VALUES (?, ?, ?, ?)
	`, uuid.New().String(), productID.String(), price, time.Now())
	return err
}

// GetPriceHistory implements Storage.GetPriceHistory
func (s *SQLiteStorage) GetPriceHistory(ctx context.Context, productID uuid.UUID, sinceDays int) ([]*models.PriceHistory, error) {
	query := `SELECT id, product_id, price, created_at FROM price_history WHERE product_id = ?`
	args := []any{productID.String()}
	if sinceDays > 0 {
		query += " AND created_at >= ?"
		args = append(args, time.Now().Add(-time.Duration(sinceDays)*24*time.Hour))
	}
	query += " ORDER BY created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []*models.PriceHistory
	for rows.Next() {
		var ph models.PriceHistory
		if err := rows.Scan(&ph.ID, &ph.ProductID, &ph.Price, &ph.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &ph)
	}
	return items, rows.Err()
}

// CreateAlert implements Storage.CreateAlert
func (s *SQLiteStorage) CreateAlert(ctx context.Context, alert *models.Alert) error {
	if alert.ID == uuid.Nil { alert.ID = uuid.New() }
	if alert.CreatedAt.IsZero() { alert.CreatedAt = time.Now() }
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO alerts (id, product_id, target_price, is_active, notification_type, created_at, notified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, alert.ID.String(), alert.ProductID.String(), alert.TargetPrice, boolToInt(alert.IsActive), alert.NotificationType, alert.CreatedAt, nullTime(alert.NotifiedAt))
	return err
}

// GetAlertByID implements Storage.GetAlertByID
func (s *SQLiteStorage) GetAlertByID(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	var a models.Alert
	err := s.db.QueryRowContext(ctx, `
		SELECT id, product_id, target_price, is_active, notification_type, created_at, notified_at
		FROM alerts WHERE id = ?
	`, id.String()).Scan(&a.ID, &a.ProductID, &a.TargetPrice, &a.IsActive, &a.NotificationType, &a.CreatedAt, &a.NotifiedAt)
	if err == sql.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }
	return &a, nil
}

// ListAlerts implements Storage.ListAlerts
func (s *SQLiteStorage) ListAlerts(ctx context.Context, limit, offset int) ([]*models.Alert, error) {
	query := `SELECT id, product_id, target_price, is_active, notification_type, created_at, notified_at FROM alerts ORDER BY created_at DESC`
	args := []any{}
	if limit > 0 { query += " LIMIT ?"; args = append(args, limit) }
	if offset > 0 { query += " OFFSET ?"; args = append(args, offset) }

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var items []*models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.ProductID, &a.TargetPrice, &a.IsActive, &a.NotificationType, &a.CreatedAt, &a.NotifiedAt); err != nil {
			return nil, err
		}
		items = append(items, &a)
	}
	return items, rows.Err()
}

// GetActiveAlertsForProduct implements Storage.GetActiveAlertsForProduct
func (s *SQLiteStorage) GetActiveAlertsForProduct(ctx context.Context, productID uuid.UUID) ([]*models.Alert, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, product_id, target_price, is_active, notification_type, created_at, notified_at
		FROM alerts WHERE product_id = ? AND is_active = 1
	`, productID.String())
	if err != nil { return nil, err }
	defer rows.Close()

	var items []*models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.ProductID, &a.TargetPrice, &a.IsActive, &a.NotificationType, &a.CreatedAt, &a.NotifiedAt); err != nil {
			return nil, err
		}
		items = append(items, &a)
	}
	return items, rows.Err()
}

// UpdateAlert implements Storage.UpdateAlert
func (s *SQLiteStorage) UpdateAlert(ctx context.Context, alert *models.Alert) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE alerts SET product_id = ?, target_price = ?, is_active = ?, notification_type = ?, created_at = ?, notified_at = ?
		WHERE id = ?
	`, alert.ProductID.String(), alert.TargetPrice, boolToInt(alert.IsActive), alert.NotificationType, alert.CreatedAt, nullTime(alert.NotifiedAt), alert.ID.String())
	return err
}

// DeleteAlert implements Storage.DeleteAlert
func (s *SQLiteStorage) DeleteAlert(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM alerts WHERE id = ?`, id.String())
	return err
}

func boolToInt(b bool) int { if b { return 1 }; return 0 }

// nullTime returns either the given time or NULL if zero-value
func nullTime(t time.Time) any { if t.IsZero() { return nil }; return t }
