package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

// Server represents the API server
type Server struct {
	httpServer *http.Server
	handler    *Handler
}

// NewServer creates a new API server
func NewServer(cfg Config, storage storage.Storage) *Server {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create handler and register routes
	handler := NewHandler(storage)
	handler.RegisterRoutes(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		handler:    handler,
	}
}

// Start starts the API server
func (s *Server) Start() error {
	log.Info().Str("address", s.httpServer.Addr).Msg("Starting API server")
	return s.httpServer.ListenAndServe()
}

// Stop gracefully shuts down the API server
func (s *Server) Stop(ctx context.Context) error {
	log.Info().Msg("Shutting down API server")
	return s.httpServer.Shutdown(ctx)
}

// Config holds the API server configuration
type Config struct {
	Address         string
	Environment     string
	AllowedOrigins  []string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}
