package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/PedroM2626/PriceWatcher/internal/models"
)

// Storage defines the interface for database operations
// Implementations: SQLite (default). You can add Postgres or others by implementing this interface.
type Storage interface {
	// Product operations
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	ListProducts(ctx context.Context, limit, offset int) ([]*models.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Price history operations
	AddPriceHistory(ctx context.Context, productID uuid.UUID, price float64) error
	GetPriceHistory(ctx context.Context, productID uuid.UUID, sinceDays int) ([]*models.PriceHistory, error)

	// Alert operations
	CreateAlert(ctx context.Context, alert *models.Alert) error
	GetAlertByID(ctx context.Context, id uuid.UUID) (*models.Alert, error)
	ListAlerts(ctx context.Context, limit, offset int) ([]*models.Alert, error)
	GetActiveAlertsForProduct(ctx context.Context, productID uuid.UUID) ([]*models.Alert, error)
	UpdateAlert(ctx context.Context, alert *models.Alert) error
	DeleteAlert(ctx context.Context, id uuid.UUID) error

	// Close closes the database connection
	Close() error
}
