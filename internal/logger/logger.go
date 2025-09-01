package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds the logger configuration
type Config struct {
	Level      string `yaml:"level"`
	JSONOutput bool   `yaml:"json_output"`
	Caller     bool   `yaml:"caller"`
}

// Init initializes the global logger with the given configuration
func Init(cfg Config) error {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// Configure logger output
	var output io.Writer = os.Stdout
	if !cfg.JSONOutput {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	// Initialize global logger
	logger := zerolog.New(output).
		With().
		Timestamp().
		Logger()

	// Add caller info if enabled
	if cfg.Caller {
		logger = logger.With().Caller().Logger()
	}

	// Set as global logger
	log.Logger = logger

	return nil
}

// Logger returns a new logger with the specified context
func Logger() *zerolog.Logger {
	return &log.Logger
}

// WithContext adds context fields to the logger
func WithContext(fields map[string]interface{}) zerolog.Logger {
	return log.Logger.With().Fields(fields).Logger()
}

// Debug logs a debug message
func Debug(msg string, fields ...interface{}) {
	log.Debug().Fields(parseFields(fields...)).Msg(msg)
}

// Info logs an info message
func Info(msg string, fields ...interface{}) {
	log.Info().Fields(parseFields(fields...)).Msg(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...interface{}) {
	log.Warn().Fields(parseFields(fields...)).Msg(msg)
}

// Error logs an error message
func Error(msg string, err error, fields ...interface{}) {
	log.Error().Err(err).Fields(parseFields(fields...)).Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, err error, fields ...interface{}) {
	log.Fatal().Err(err).Fields(parseFields(fields...)).Msg(msg)
}

// parseFields converts variadic key-value pairs to a map
func parseFields(fields ...interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if !ok {
				continue
			}
			result[key] = fields[i+1]
		}
	}
	return result
}
