package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/PedroM2626/PriceWatcher/internal/models"
)

// Storage defines the interface for database operations
type Storage interface {
	// Product operations
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductByURL(ctx context.Context, url string) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Price history operations
	AddPriceHistory(ctx context.Context, productID uuid.UUID, price float64) error
	GetPriceHistory(ctx context.Context, productID uuid.UUID, days int) ([]*models.PriceHistory, error)

	// Alert operations
	CreateAlert(ctx context.Context, alert *models.Alert) error
	GetActiveAlertsForProduct(ctx context.Context, productID uuid.UUID) ([]*models.Alert, error)
	UpdateAlert(ctx context.Context, alert *models.Alert) error
	DeleteAlert(ctx context.Context, id uuid.UUID) error

	// Close closes the database connection
	Close() error
}

// PostgresStorage implements Storage interface for PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new PostgreSQL storage instance
func NewPostgresStorage(cfg config.DatabaseConfig) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresStorage{db: db}, nil
}

// Close implements Storage.Close
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

// CreateProduct implements Storage.CreateProduct
func (s *PostgresStorage) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	err := s.db.QueryRowContext(
		ctx,
		query,
		product.ID,
		product.Name,
		product.URL,
		product.ImageURL,
		product.CurrentPrice,
		product.Currency,
		product.IsAvailable,
		product.Website,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetProductByID implements Storage.GetProductByID
func (s *PostgresStorage) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, name, url, image_url, current_price, currency, is_available, website, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.URL,
		&product.ImageURL,
		&product.CurrentPrice,
		&product.Currency,
		&product.IsAvailable,
		&product.Website,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// Implement other required methods (UpdateProduct, ListProducts, DeleteProduct, etc.)
// ...

// AddPriceHistory implements Storage.AddPriceHistory
func (s *PostgresStorage) AddPriceHistory(ctx context.Context, productID uuid.UUID, price float64) error {
	query := `
		INSERT INTO price_history (id, product_id, price, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		uuid.New(),
		productID,
		price,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to add price history: %w", err)
	}

	return nil
}

// CreateAlert implements Storage.CreateAlert
func (s *PostgresStorage) CreateAlert(ctx context.Context, alert *models.Alert) error {
	query := `
		INSERT INTO alerts (id, product_id, target_price, is_active, notification_type, created_at, notified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	alert.CreatedAt = now

	_, err := s.db.ExecContext(
		ctx,
		query,
		alert.ID,
		alert.ProductID,
		alert.TargetPrice,
		alert.IsActive,
		alert.NotificationType,
		alert.CreatedAt,
		alert.NotifiedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}

	return nil
}
