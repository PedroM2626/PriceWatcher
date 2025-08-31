package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PedroM2626/PriceWatcher/internal/config"
	"github.com/PedroM2626/PriceWatcher/internal/notifier"
	"github.com/PedroM2626/PriceWatcher/internal/scraper"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

func main() {
	// Create a context that will be cancelled on interrupt signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

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

	// Initialize notifier
	notifierSvc, err := notifier.NewNotificationService(notifier.NotificationConfig{
		Email:    notifier.EmailConfig(cfg.Notifier.Email),
		Telegram: notifier.TelegramConfig(cfg.Notifier.Telegram),
	})
	if err != nil {
		log.Fatalf("Failed to initialize notifier: %v", err)
	}

	// Initialize scraper
	scraperCfg := scraper.ScraperConfig{
		UserAgent:      cfg.Scraper.UserAgent,
		RequestDelay:   cfg.Scraper.RequestDelay,
		RequestTimeout: cfg.Scraper.RequestTimeout,
		Workers:        cfg.Scraper.Workers,
	}

	scraperSvc := scraper.NewScraper(db, scraperCfg, notifierSvc)

	// Start the scraper in a goroutine
	go func() {
		if err := scraperSvc.Run(ctx); err != nil {
			log.Printf("Scraper error: %v", err)
		}
	}()

	log.Println("PriceWatcher started. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down PriceWatcher...")
}
