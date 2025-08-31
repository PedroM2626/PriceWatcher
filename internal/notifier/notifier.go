package notifier

import (
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/PedroM2626/PriceWatcher/internal/models"
)

// Notifier defines the interface for sending notifications
type Notifier interface {
	// Send sends a notification with the given message
	Send(ctx context.Context, recipient string, subject, message string) error
	// SendPriceAlert sends a price alert notification
	SendPriceAlert(ctx context.Context, alert *models.Alert, product *models.Product, oldPrice float64) error
}

// NotificationConfig holds configuration for notifications
type NotificationConfig struct {
	Email    EmailConfig
	Telegram TelegramConfig
}

// EmailConfig holds email notification configuration
type EmailConfig struct {
	Enabled  bool
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

// TelegramConfig holds Telegram notification configuration
type TelegramConfig struct {
	Enabled bool
	Token   string
	ChatID  int64
}

// NotificationService handles sending notifications through multiple channels
type NotificationService struct {
	emailNotifier    *EmailNotifier
	telegramNotifier *TelegramNotifier
}

// NewNotificationService creates a new notification service
func NewNotificationService(cfg NotificationConfig) (*NotificationService, error) {
	var emailNotifier *EmailNotifier
	var telegramNotifier *TelegramNotifier
	var err error

	if cfg.Email.Enabled {
		emailNotifier = NewEmailNotifier(cfg.Email)
	}

	if cfg.Telegram.Enabled {
		telegramNotifier, err = NewTelegramNotifier(cfg.Telegram)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Telegram notifier: %w", err)
		}
	}

	return &NotificationService{
		emailNotifier:    emailNotifier,
		telegramNotifier: telegramNotifier,
	}, nil
}

// Send sends a notification to all enabled channels
func (s *NotificationService) Send(ctx context.Context, recipient, subject, message string) error {
	// Send email notification if enabled
	if s.emailNotifier != nil {
		if err := s.emailNotifier.Send(ctx, recipient, subject, message); err != nil {
			return fmt.Errorf("failed to send email notification: %w", err)
		}
	}

	// Send Telegram notification if enabled
	if s.telegramNotifier != nil {
		if err := s.telegramNotifier.Send(ctx, "", subject, message); err != nil {
			return fmt.Errorf("failed to send Telegram notification: %w", err)
		}
	}

	return nil
}

// SendPriceAlert sends a price alert notification
func (s *NotificationService) SendPriceAlert(ctx context.Context, alert *models.Alert, product *models.Product, oldPrice float64) error {
	// Prepare the message
	subject := fmt.Sprintf("🚨 Price Alert: %s", product.Name)
	
	tmpl := `
🔔 *Price Alert!*

*{{.ProductName}}*

🏷 *Old Price:* {{.OldPrice}} {{.Currency}}
💰 *New Price:* {{.NewPrice}} {{.Currency}}
📉 *Price Drop:* {{.PriceDrop}}% ({{.PriceDifference}} {{.Currency}})

🛒 [View on Website]({{.ProductURL}})
`

	priceDrop := ((oldPrice - product.CurrentPrice) / oldPrice) * 100
	priceDifference := oldPrice - product.CurrentPrice

	data := struct {
		ProductName   string
		OldPrice      float64
		NewPrice      float64
		Currency      string
		PriceDrop     float64
		PriceDifference float64
		ProductURL    string
	}{
		ProductName:   product.Name,
		OldPrice:      oldPrice,
		NewPrice:      product.CurrentPrice,
		Currency:      product.Currency,
		PriceDrop:     priceDrop,
		PriceDifference: priceDifference,
		ProductURL:    product.URL,
	}

	var messageBuf strings.Builder
	t := template.Must(template.New("alert").Parse(tmpl))
	if err := t.Execute(&messageBuf, data); err != nil {
		return fmt.Errorf("failed to generate alert message: %w", err)
	}

	// Send the notification
	switch alert.NotificationType {
	case "email":
		if s.emailNotifier == nil {
			return fmt.Errorf("email notifications are not enabled")
		}
		// In a real implementation, you would get the user's email from the database
		return s.emailNotifier.Send(ctx, "user@example.com", subject, messageBuf.String())
		
	case "telegram":
		if s.telegramNotifier == nil {
			return fmt.Errorf("Telegram notifications are not enabled")
		}
		return s.telegramNotifier.Send(ctx, "", subject, messageBuf.String())
		
	default:
		return fmt.Errorf("unsupported notification type: %s", alert.NotificationType)
	}
}

// EmailNotifier handles sending email notifications
type EmailNotifier struct {
	config EmailConfig
}

// NewEmailNotifier creates a new email notifier
func NewEmailNotifier(config EmailConfig) *EmailNotifier {
	return &EmailNotifier{config: config}
}

// Send sends an email notification
func (n *EmailNotifier) Send(ctx context.Context, to, subject, message string) error {
	// TODO: Implement email sending using net/smtp or a third-party package
	// This is a placeholder implementation
	return nil
}

// TelegramNotifier handles sending Telegram notifications
type TelegramNotifier struct {
	config TelegramConfig
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(config TelegramConfig) (*TelegramNotifier, error) {
	// TODO: Initialize Telegram bot
	return &TelegramNotifier{config: config}, nil
}

// Send sends a Telegram notification
func (n *TelegramNotifier) Send(ctx context.Context, chatID string, subject, message string) error {
	// TODO: Implement Telegram message sending using the Telegram Bot API
	// This is a placeholder implementation
	return nil
}
