package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/PedroM2626/PriceWatcher/internal/models"
	"github.com/PedroM2626/PriceWatcher/internal/storage"
)

// Handler handles HTTP requests
type Handler struct {
	storage storage.Storage
}

// NewHandler creates a new handler instance
func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Public routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.register)
			auth.POST("/login", h.login)
			auth.POST("/refresh", h.refreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(h.authMiddleware())
		{
			// Products
			products := protected.Group("/products")
			{
				products.GET("", h.listProducts)
				products.POST("", h.createProduct)
				products.GET("/:id", h.getProduct)
				products.PUT("/:id", h.updateProduct)
				products.DELETE("/:id", h.deleteProduct)
			}

			// Alerts
			alerts := protected.Group("/alerts")
			{
				alerts.GET("", h.listAlerts)
				alerts.POST("", h.createAlert)
				alerts.GET("/:id", h.getAlert)
				alerts.PUT("/:id", h.updateAlert)
				alerts.DELETE("/:id", h.deleteAlert)
			}
		}
	}
}

// authMiddleware handles JWT authentication
func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// TODO: Validate JWT token
		// claims, err := validateToken(tokenString)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	c.Abort()
		// 	return
		// }

		// c.Set("userID", claims.UserID)
		c.Next()
	}
}

// Product handlers
func (h *Handler) listProducts(c *gin.Context) {
	// TODO: Implement pagination
	products, err := h.storage.ListProducts(c.Request.Context(), 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) createProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate new ID
	product.ID = uuid.New()

	if err := h.storage.CreateProduct(c.Request.Context(), &product); err != nil {
		log.Error().Err(err).Msg("Failed to create product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// ... (other product handler methods)

// Alert handlers
func (h *Handler) listAlerts(c *gin.Context) {
	// TODO: Get user ID from context
	// userID := c.MustGet("userID").(string)

	alerts, err := h.storage.ListAlerts(c.Request.Context(), "userID")
	if err != nil {
		log.Error().Err(err).Msg("Failed to list alerts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list alerts"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// ... (other alert handler methods)

// Auth handlers
func (h *Handler) register(c *gin.Context) {
	// TODO: Implement user registration
}

func (h *Handler) login(c *gin.Context) {
	// TODO: Implement user login
}

func (h *Handler) refreshToken(c *gin.Context) {
	// TODO: Implement token refresh
}
