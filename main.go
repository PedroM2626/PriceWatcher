package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PedroM2626/PriceWatcher/internal/config"
	"github.com/PedroM2626/PriceWatcher/internal/scraper"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage
	db, err := storage.NewPostgresStorage(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer db.Close()

	// Initialize scraper
	scraper := scraper.NewScraper(db, cfg.Scraper)

	// Start scraping in a separate goroutine
	go func() {
		if err := scraper.Run(); err != nil {
			log.Fatalf("Scraper error: %v", err)
		}
	}()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down PriceWatcher...")
}
