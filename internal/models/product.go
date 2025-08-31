package models

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product being tracked
type Product struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	URL          string    `json:"url" db:"url"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	CurrentPrice float64   `json:"current_price" db:"current_price"`
	Currency     string    `json:"currency" db:"currency"`
	IsAvailable  bool      `json:"is_available" db:"is_available"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Website      string    `json:"website" db:"website"`
}

// PriceHistory represents the price history of a product
type PriceHistory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Alert represents a price alert set by the user
type Alert struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ProductID    uuid.UUID `json:"product_id" db:"product_id"`
	TargetPrice  float64   `json:"target_price" db:"target_price"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	NotifiedAt   time.Time `json:"notified_at,omitempty" db:"notified_at"`
	NotificationType string `json:"notification_type" db:"notification_type"` // email, telegram, etc.
}
