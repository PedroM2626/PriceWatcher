package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PedroM2626/PriceWatcher/internal/api"
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
	db, err := storage.NewStorage(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer db.Close()

	// Initialize and start API server
	apiCfg := api.Config{
		Address:        ":8080",
		Environment:    "development",
		AllowedOrigins: []string{"http://localhost:3000"},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
	}

	server := api.NewServer(apiCfg, db)

	// Start API server in a goroutine
	go func() {
		log.Printf("Starting API server on %s", apiCfg.Address)
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// Initialize scraper
	scraper := scraper.NewScraper(db, cfg.Scraper)

	// Start scraping in a separate goroutine
	go func() {
		log.Println("Starting price scraper...")
		if err := scraper.Run(); err != nil {
			log.Fatalf("Scraper error: %v", err)
		}
	}()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down PriceWatcher...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("PriceWatcher has been shut down")
}
