package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Scraper  ScraperConfig  `yaml:"scraper"`
	Notifier NotifierConfig `yaml:"notifier"`
	Server   ServerConfig   `yaml:"server"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver string `yaml:"driver"` // sqlite3, postgres, etc.
	DSN    string `yaml:"dsn"`    // Data Source Name (e.g., file path for SQLite, connection string for PostgreSQL)
	
	// PostgreSQL specific fields (optional, can be used to build DSN)
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	DBName   string `yaml:"dbname,omitempty"`
	SSLMode  string `yaml:"sslmode,omitempty"`
}

// ScraperConfig holds scraper configuration
type ScraperConfig struct {
	UserAgent      string        `yaml:"user_agent"`
	RequestDelay   time.Duration `yaml:"request_delay"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
	Workers        int           `yaml:"workers"`
}

// NotifierConfig holds notifier configuration
type NotifierConfig struct {
	Email    EmailConfig    `yaml:"email"`
	Telegram TelegramConfig `yaml:"telegram"`
}

// EmailConfig holds email notification configuration
type EmailConfig struct {
	Enabled  bool   `yaml:"enabled"`
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

// TelegramConfig holds Telegram notification configuration
type TelegramConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	ChatID  int64  `yaml:"chat_id"`
}

// ServerConfig holds web server configuration
type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Fallbacks from environment for secret-free config
	if cfg.Database.DSN == "" {
		if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
			cfg.Database.DSN = dsn
			if cfg.Database.Driver == "" { cfg.Database.Driver = "postgres" }
		}
	}

	return &cfg, nil
}
