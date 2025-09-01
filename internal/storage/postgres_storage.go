package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/PedroM2626/PriceWatcher/internal/models"
)

// PostgresStorage implements Storage interface for PostgreSQL/Supabase
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage connects to Postgres and ensures schema exists
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := pgEnsureSchema(db); err != nil {
		db.Close()
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func pgEnsureSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			url TEXT NOT NULL UNIQUE,
			image_url TEXT,
			current_price DOUBLE PRECISION NOT NULL DEFAULT 0,
			currency TEXT NOT NULL DEFAULT 'BRL',
			is_available BOOLEAN NOT NULL DEFAULT TRUE,
			website TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS price_history (
			id UUID PRIMARY KEY,
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			price DOUBLE PRECISION NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_product_id ON price_history(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_price_history_created_at ON price_history(created_at)`,
		`CREATE TABLE IF NOT EXISTS alerts (
			id UUID PRIMARY KEY,
			product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
			target_price DOUBLE PRECISION NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			notification_type TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			notified_at TIMESTAMPTZ
		)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_product_id ON alerts(product_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_is_active ON alerts(is_active)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("schema init error: %w (stmt=%s)", err, firstLine(s))
		}
	}
	return nil
}

func firstLine(s string) string {
	ls := strings.Split(strings.TrimSpace(s), "\n")
	if len(ls) > 0 { return ls[0] }
	return s
}

// Close implements Storage.Close
func (s *PostgresStorage) Close() error { return s.db.Close() }

// CreateProduct implements Storage.CreateProduct
func (s *PostgresStorage) CreateProduct(ctx context.Context, p *models.Product) error {
	if p.ID == uuid.Nil { p.ID = uuid.New() }
	now := time.Now()
	if p.CreatedAt.IsZero() { p.CreatedAt = now }
	p.UpdatedAt = now
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO products (id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, p.ID, p.Name, p.URL, p.ImageURL, p.CurrentPrice, p.Currency, p.IsAvailable, p.Website, p.CreatedAt, p.UpdatedAt)
	return err
}

// GetProductByID implements Storage.GetProductByID
func (s *PostgresStorage) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var p models.Product
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at
		FROM products WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.URL, &p.ImageURL, &p.CurrentPrice, &p.Currency, &p.IsAvailable, &p.Website, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }
	return &p, nil
}

// UpdateProduct implements Storage.UpdateProduct
func (s *PostgresStorage) UpdateProduct(ctx context.Context, p *models.Product) error {
	p.UpdatedAt = time.Now()
	_, err := s.db.ExecContext(ctx, `
		UPDATE products SET name=$1, url=$2, image_url=$3, current_price=$4, currency=$5, is_available=$6, website=$7, updated_at=$8
		WHERE id=$9
	`, p.Name, p.URL, p.ImageURL, p.CurrentPrice, p.Currency, p.IsAvailable, p.Website, p.UpdatedAt, p.ID)
	return err
}

// ListProducts implements Storage.ListProducts
func (s *PostgresStorage) ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	query := `SELECT id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at FROM products ORDER BY created_at DESC`
	args := []any{}
	if limit > 0 { query += " LIMIT $1"; args = append(args, limit) }
	if offset > 0 {
		if len(args) == 0 { query += " OFFSET $1" } else { query += " OFFSET $2" }
		args = append(args, offset)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.URL, &p.ImageURL, &p.CurrentPrice, &p.Currency, &p.IsAvailable, &p.Website, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}

// DeleteProduct implements Storage.DeleteProduct
func (s *PostgresStorage) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}

// AddPriceHistory implements Storage.AddPriceHistory
func (s *PostgresStorage) AddPriceHistory(ctx context.Context, productID uuid.UUID, price float64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO price_history (id, product_id, price, created_at) VALUES ($1,$2,$3,$4)
	`, uuid.New(), productID, price, time.Now())
	return err
}

// GetPriceHistory implements Storage.GetPriceHistory
func (s *PostgresStorage) GetPriceHistory(ctx context.Context, productID uuid.UUID, sinceDays int) ([]*models.PriceHistory, error) {
	query := `SELECT id, product_id, price, created_at FROM price_history WHERE product_id=$1`
	args := []any{productID}
	if sinceDays > 0 { query += " AND created_at >= $2"; args = append(args, time.Now().Add(-time.Duration(sinceDays)*24*time.Hour)) }
	query += " ORDER BY created_at DESC"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*models.PriceHistory
	for rows.Next() {
		var ph models.PriceHistory
		if err := rows.Scan(&ph.ID, &ph.ProductID, &ph.Price, &ph.CreatedAt); err != nil { return nil, err }
		out = append(out, &ph)
	}
	return out, rows.Err()
}

// CreateAlert implements Storage.CreateAlert
func (s *PostgresStorage) CreateAlert(ctx context.Context, a *models.Alert) error {
	if a.ID == uuid.Nil { a.ID = uuid.New() }
	if a.CreatedAt.IsZero() { a.CreatedAt = time.Now() }
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO alerts (id, product_id, target_price, is_active, notification_type, created_at, notified_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, a.ID, a.ProductID, a.TargetPrice, a.IsActive, a.NotificationType, a.CreatedAt, nullPGTime(a.NotifiedAt))
	return err
}

// GetAlertByID implements Storage.GetAlertByID
func (s *PostgresStorage) GetAlertByID(ctx context.Context, id uuid.UUID) (*models.Alert, error) {
	var a models.Alert
	err := s.db.QueryRowContext(ctx, `
		SELECT id, product_id, target_price, is_active, notification_type, created_at, notified_at
		FROM alerts WHERE id=$1
	`, id).Scan(&a.ID, &a.ProductID, &a.TargetPrice, &a.IsActive, &a.NotificationType, &a.CreatedAt, &a.NotifiedAt)
	if err == sql.ErrNoRows { return nil, nil }
	if err != nil { return nil, err }
	return &a, nil
}

// ListAlerts implements Storage.ListAlerts
func (s *PostgresStorage) ListAlerts(ctx context.Context, limit, offset int) ([]*models.Alert, error) {
	query := `SELECT id, product_id, target_price, is_active, notification_type, created_at, notified_at FROM alerts ORDER BY created_at DESC`
	args := []any{}
	if limit > 0 { query += " LIMIT $1"; args = append(args, limit) }
	if offset > 0 {
		if len(args) == 0 { query += " OFFSET $1" } else { query += " OFFSET $2" }
		args = append(args, offset)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.ProductID, &a.TargetPrice, &a.IsActive, &a.NotificationType, &a.CreatedAt, &a.NotifiedAt); err != nil { return nil, err }
		out = append(out, &a)
	}
	return out, rows.Err()
}

// GetActiveAlertsForProduct implements Storage.GetActiveAlertsForProduct
func (s *PostgresStorage) GetActiveAlertsForProduct(ctx context.Context, productID uuid.UUID) ([]*models.Alert, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, product_id, target_price, is_active, notification_type, created_at, notified_at
		FROM alerts WHERE product_id=$1 AND is_active=TRUE
	`, productID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []*models.Alert
	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(&a.ID, &a.ProductID, &a.TargetPrice, &a.IsActive, &a.NotificationType, &a.CreatedAt, &a.NotifiedAt); err != nil { return nil, err }
		out = append(out, &a)
	}
	return out, rows.Err()
}

// UpdateAlert implements Storage.UpdateAlert
func (s *PostgresStorage) UpdateAlert(ctx context.Context, a *models.Alert) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE alerts SET product_id=$1, target_price=$2, is_active=$3, notification_type=$4, created_at=$5, notified_at=$6
		WHERE id=$7
	`, a.ProductID, a.TargetPrice, a.IsActive, a.NotificationType, a.CreatedAt, nullPGTime(a.NotifiedAt), a.ID)
	return err
}

// DeleteAlert implements Storage.DeleteAlert
func (s *PostgresStorage) DeleteAlert(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM alerts WHERE id=$1`, id)
	return err
}

func nullPGTime(t time.Time) any { if t.IsZero() { return nil }; return t }
