package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PedroM2626/PriceWatcher/internal/config"
	"github.com/PedroM2626/PriceWatcher/internal/scraper"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage (SQLite expected)
	db, err := storage.NewStorage(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer db.Close()

	// Initialize scraper service
	s := scraper.NewScraper(db, scraper.ScraperConfig{
		UserAgent:      cfg.Scraper.UserAgent,
		RequestDelay:   cfg.Scraper.RequestDelay,
		RequestTimeout: cfg.Scraper.RequestTimeout,
		Workers:        cfg.Scraper.Workers,
	})

	// Start scraper in background (no-op for now)
	go func() {
		if err := s.Run(); err != nil {
			log.Printf("Scraper error: %v", err)
		}
	}()

	log.Println("PriceWatcher service started. Press Ctrl+C to stop.")

	<-sigChan
	log.Println("Shutting down PriceWatcher...")
	_ = ctx
}
