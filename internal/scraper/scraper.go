package scraper

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"github.com/PedroM2626/PriceWatcher/internal/models"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

// ScraperConfig holds configuration for the scraper
type ScraperConfig struct {
	UserAgent      string
	RequestDelay   time.Duration
	RequestTimeout time.Duration
	Workers        int
}

// PriceScraper scrapes product pages periodically
type PriceScraper struct {
	storage storage.Storage
	config  ScraperConfig
}

// NewScraper creates a new instance of PriceScraper
func NewScraper(storage storage.Storage, cfg ScraperConfig) *PriceScraper {
	return &PriceScraper{storage: storage, config: cfg}
}

// Run starts the scraping process (no-op placeholder loop for now)
func (s *PriceScraper) Run() error {
	return nil
}

// Scrape extracts product information from the given URL
func (s *PriceScraper) Scrape(ctx context.Context, productURL string) (*models.Product, error) {
	u, err := url.Parse(productURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	c := colly.NewCollector(colly.UserAgent(s.config.UserAgent))
	c.SetRequestTimeout(s.config.RequestTimeout)

	product := &models.Product{
		ID:        uuid.New(),
		URL:       productURL,
		Website:   u.Hostname(),
		Currency:  "BRL",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Try to get the title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		if product.Name == "" {
			product.Name = e.Text
		}
	})

	// Best-effort visit (errors are returned)
	if err := c.Visit(productURL); err != nil {
		return nil, fmt.Errorf("failed to visit url: %w", err)
	}

	return product, nil
}
