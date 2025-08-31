package scraper

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/PedroM2626/PriceWatcher/internal/models"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

// Scraper defines the interface for product scrapers
type Scraper interface {
	// Scrape extracts product information from the given URL
	Scrape(ctx context.Context, productURL string) (*models.Product, error)
}

// ScraperConfig holds configuration for the scraper
type ScraperConfig struct {
	UserAgent      string
	RequestDelay   time.Duration
	RequestTimeout time.Duration
	Workers        int
}

// PriceScraper implements the Scraper interface
type PriceScraper struct {
	storage storage.Storage
	config  ScraperConfig
	mu      sync.Mutex
}

// NewScraper creates a new instance of PriceScraper
func NewScraper(storage storage.Storage, cfg ScraperConfig) *PriceScraper {
	return &PriceScraper{
		storage: storage,
		config:  cfg,
	}
}

// Run starts the scraping process
func (s *PriceScraper) Run() error {
	// TODO: Implement the main scraping loop
	// 1. Fetch products from the database
	// 2. Scrape each product's page
	// 3. Update product information and price history
	// 4. Check for price alerts
	// 5. Wait for the next interval

	return nil
}

// Scrape extracts product information from the given URL
func (s *PriceScraper) Scrape(ctx context.Context, productURL string) (*models.Product, error) {
	// Parse the URL to determine which scraper to use
	u, err := url.Parse(productURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create a new collector
	c := colly.NewCollector(
		colly.UserAgent(s.config.UserAgent),
		colly.AllowURLRevisit(),
	)

	// Set request timeout
	c.SetRequestTimeout(s.config.RequestTimeout)

	// Create a new product
	product := &models.Product{
		ID:        uuid.New(),
		URL:       productURL,
		Website:   u.Hostname(),
		CreatedAt: time.Now(),
	}

	// Get the appropriate scraper for the website
	scraper, err := s.getScraperForDomain(u.Hostname())
	if err != nil {
		return nil, err
	}

	// Scrape the product page
	if err := scraper.Scrape(c, product); err != nil {
		return nil, err
	}

	return product, nil
}

// getScraperForDomain returns the appropriate scraper for the given domain
func (s *PriceScraper) getScraperForDomain(domain string) (Scraper, error) {
	switch {
	case strings.Contains(domain, "amazon"):
		return &AmazonScraper{}, nil
	case strings.Contains(domain, "mercadolivre"):
		return &MercadoLivreScraper{}, nil
	// Add more scrapers for other websites
	default:
		return &GenericScraper{}, nil
	}
}

// AmazonScraper implements Scraper for Amazon
// TODO: Implement Amazon scraper

// MercadoLivreScraper implements Scraper for Mercado Livre
// TODO: Implement Mercado Livre scraper

// GenericScraper implements Scraper for generic websites
// TODO: Implement generic scraper
