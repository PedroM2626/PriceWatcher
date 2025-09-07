# PriceWatcher

PriceWatcher is a Go application that monitors product prices across various e-commerce websites and notifies users when prices drop below a specified threshold.

## Features

- 🛒 Monitor multiple products from different websites
- 📊 Track price history and trends
- 🔔 Get notified via email or Telegram when prices drop
- ⚡ Fast and efficient web scraping
- 🏗 Extensible architecture for adding new websites
- 🐳 Docker support for easy deployment

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- (Optional) Docker and Docker Compose

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/PedroM2626/PriceWatcher.git
   cd PriceWatcher
   ```

2. Copy the example configuration file:
   ```bash
   cp config.example.yaml config.yaml
   ```

3. Edit the `config.yaml` file with your settings:
   ```yaml
   database:
     host: localhost
     port: 5432
     user: your_username
     password: your_password
     dbname: pricewatcher
     sslmode: disable
   
   # ... other settings
   ```

4. Set up the database:
   ```bash
   # Create the database
   createdb pricewatcher
   
   # Run migrations
   psql -U your_username -d pricewatcher -f migrations/001_initial_schema.up.sql
   ```

5. Install dependencies:
   ```bash
   go mod download
   ```

## Running the Application

```bash
# Build the application
go build -o pricewatcher cmd/pricewatcher/main.go

# Run the application
./pricewatcher
```

### Using Docker

1. Build the Docker image:
   ```bash
   docker-compose build
   ```

2. Start the application:
   ```bash
   docker-compose up -d
   ```

## Configuration

The application is configured using a YAML file. See `config.example.yaml` for all available options.

## Adding Products

To add a product to monitor, you can use the provided API or add it directly to the database:

```sql
INSERT INTO products (id, name, url, current_price, currency, is_available, website, created_at, updated_at)
VALUES (
  '550e8400-e29b-41d4-a716-446655440000',
  'Example Product',
  'https://example.com/product',
  99.99,
  'BRL',
  true,
  'example.com',
  NOW(),
  NOW()
);
```

## Setting Up Alerts

Create a price alert using the API or directly in the database:

```sql
INSERT INTO alerts (id, product_id, target_price, is_active, notification_type, created_at)
VALUES (
  '550e8400-e29b-41d4-a716-446655440001',
  '550e8400-e29b-41d4-a716-446655440000',
  80.00,
  true,
  'email',
  NOW()
);
```

## Development

### Project Structure

```
PriceWatcher/
├── cmd/                  # Main application entry points
├── internal/             # Private application code
│   ├── config/           # Configuration management
│   ├── models/           # Data models
│   ├── notifier/         # Notification services
│   ├── scraper/          # Web scrapers
│   └── storage/          # Database layer
├── migrations/           # Database migrations
├── pkg/                  # Public library code
├── web/                  # Web assets and templates
├── config.example.yaml   # Example configuration
└── Dockerfile            # Docker configuration
```

### Adding a New Website Scraper

1. Create a new file in `internal/scraper/` (e.g., `amazon_scraper.go`)
2. Implement the `Scraper` interface
3. Add the scraper to the `getScraperForDomain` function in `scraper.go`

### Running Tests

```bash
go test ./...
```

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Support

For support, please open an issue on GitHub.
