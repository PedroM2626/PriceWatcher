package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
	"github.com/PedroM2626/PriceWatcher/internal/models"
	"github.com/PedroM2626/PriceWatcher/internal/scraper"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

// Scheduler handles scheduling of price checks
type Scheduler struct {
	scheduler gocron.Scheduler
	scraper   *scraper.PriceScraper
	storage   storage.Storage
}

// NewScheduler creates a new scheduler instance
func NewScheduler(scraper *scraper.PriceScraper, storage storage.Storage) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		scheduler: s,
		scraper:   scraper,
		storage:   storage,
	}, nil
}

// Start starts the scheduler
func (s *Scheduler) Start() error {
	// Schedule price checks to run every hour
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(1*time.Hour),
		gocron.NewTask(s.CheckAllProducts),
	)
	if err != nil {
		return err
	}

	// Start the scheduler
	s.scheduler.Start()
	log.Info().Msg("Scheduler started")

	// Run initial check
	go s.CheckAllProducts()

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	if s.scheduler != nil {
		err := s.scheduler.Shutdown()
		if err != nil {
			log.Error().Err(err).Msg("Error shutting down scheduler")
		}
	}
}

// CheckAllProducts checks all products for price updates
func (s *Scheduler) CheckAllProducts() {
	log.Info().Msg("Starting scheduled price check")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Get all products from the database
	products, err := s.storage.ListProducts(ctx, 0, 0) // 0, 0 means get all products
	if err != nil {
		log.Error().Err(err).Msg("Failed to list products for price check")
		return
	}

	log.Info().Int("count", len(products)).Msg("Checking prices for products")

	// Check each product
	for _, product := range products {
		// Skip if the URL is empty
		if product.URL == "" {
			continue
		}

		// Scrape the product
		updatedProduct, err := s.scraper.Scrape(ctx, product.URL)
		if err != nil {
			log.Error().
				Err(err).
				Str("product_id", product.ID.String()).
				Str("url", product.URL).
				Msg("Failed to scrape product")
			continue
		}

		// Update the product in the database
		err = s.storage.UpdateProduct(ctx, updatedProduct)
		if err != nil {
			log.Error().
				Err(err).
				Str("product_id", updatedProduct.ID.String()).
				Msg("Failed to update product")
			continue
		}

		// Check if price changed
		if product.CurrentPrice != updatedProduct.CurrentPrice {
			log.Info().
				Str("product_id", updatedProduct.ID.String()).
				Float64("old_price", product.CurrentPrice).
				Float64("new_price", updatedProduct.CurrentPrice).
				Msg("Price updated")

			// Trigger price alerts
			s.checkPriceAlerts(ctx, product, updatedProduct)
		}
	}
}

// checkPriceAlerts checks if any price alerts should be triggered
func (s *Scheduler) checkPriceAlerts(ctx context.Context, oldProduct, newProduct *models.Product) {
	// Get all active alerts for this product
	alerts, err := s.storage.GetActiveAlertsForProduct(ctx, newProduct.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("product_id", newProduct.ID.String()).
			Msg("Failed to get alerts for product")
		return
	}

	for _, alert := range alerts {
		// Check if the price is below the target price
		if newProduct.CurrentPrice <= alert.TargetPrice {
			// Trigger the alert
			err := s.triggerAlert(ctx, alert, newProduct, oldProduct.CurrentPrice)
			if err != nil {
				log.Error().
					Err(err).
					Str("alert_id", alert.ID.String()).
					Msg("Failed to trigger alert")
			}

			// Mark the alert as notified
			alert.NotifiedAt = time.Now()
			err = s.storage.UpdateAlert(ctx, alert)
			if err != nil {
				log.Error().
					Err(err).
					Str("alert_id", alert.ID.String()).
					Msg("Failed to update alert")
			}
		}
	}
}

// triggerAlert triggers a price alert
func (s *Scheduler) triggerAlert(ctx context.Context, alert *models.Alert, product *models.Product, oldPrice float64) error {
	// TODO: Implement alert triggering logic
	// This would use the notifier package to send notifications
	// based on the alert's notification type (email, telegram, etc.)

	log.Info().
		Str("alert_id", alert.ID.String()).
		Str("product_id", product.ID.String()).
		Float64("target_price", alert.TargetPrice).
		Float64("current_price", product.CurrentPrice).
		Msg("Price alert triggered")

	return nil
}
