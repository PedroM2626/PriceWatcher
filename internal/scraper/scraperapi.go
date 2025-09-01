package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/PedroM2626/PriceWatcher/internal/models"
)

// ScraperAPIClient handles communication with the ScraperAPI
// Documentation: https://www.scraperapi.com/documentation/
type ScraperAPIClient struct {
	apiKey string
	client *http.Client
}

// NewScraperAPIClient creates a new ScraperAPI client
func NewScraperAPIClient(apiKey string) *ScraperAPIClient {
	return &ScraperAPIClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ScrapeProduct scrapes product information using ScraperAPI
func (c *ScraperAPIClient) ScrapeProduct(ctx context.Context, productURL string) (*models.Product, error) {
	// Build the ScraperAPI URL
	baseURL := "http://api.scraperapi.com"
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("url", productURL)
	params.Add("render", "true") // Enable JavaScript rendering

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Send the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("scraperapi request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	// Note: In a real implementation, you would parse the HTML response
	// to extract product information. This is a simplified example.
	product := &models.Product{
		ID:        uuid.New(),
		URL:       productURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// TODO: Parse the HTML response to extract product details
	// This would involve using goquery or similar to parse the HTML
	// and extract the product name, price, etc.

	return product, nil
}

// ScraperAPIResponse represents the response from ScraperAPI
type ScraperAPIResponse struct {
	Status  string `json:"status"`
	Request struct {
		URL         string `json:"url"`
		Success     bool   `json:"success"`
		StatusCode  int    `json:"status_code"`
		ContentType string `json:"content_type"`
	} `json:"request"`
	HTML string `json:"html"`
}
