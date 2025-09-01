package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/google/uuid"
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
	api := router.Group("/api/v1")
	{
		// Auth routes (basic placeholders returning 501)
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.register)
			auth.POST("/login", h.login)
			auth.POST("/refresh", h.refreshToken)
		}

		// Protected routes (middleware can be enhanced later)
		protected := api.Group("")
		protected.Use(h.authMiddleware())
		{
			products := protected.Group("/products")
			{
				products.GET("", h.listProducts)
				products.POST("", h.createProduct)
				products.GET(":id", h.getProduct)
				products.PUT(":id", h.updateProduct)
				products.DELETE(":id", h.deleteProduct)
			}

			alerts := protected.Group("/alerts")
			{
				alerts.GET("", h.listAlerts)
				alerts.POST("", h.createAlert)
				alerts.GET(":id", h.getAlert)
				alerts.PUT(":id", h.updateAlert)
				alerts.DELETE(":id", h.deleteAlert)
			}
		}
	}
}

// authMiddleware handles JWT authentication
func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, allow all requests. Add JWT validation here if needed.
		c.Next()
	}
}

// Product handlers
func (h *Handler) listProducts(c *gin.Context) {
	products, err := h.storage.ListProducts(c.Request.Context(), 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": products})
}

func (h *Handler) createProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if product.ID == uuid.Nil { product.ID = uuid.New() }
	if err := h.storage.CreateProduct(c.Request.Context(), &product); err != nil {
		log.Error().Err(err).Msg("Failed to create product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *Handler) getProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	product, err := h.storage.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product"})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) updateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID = id
	if err := h.storage.UpdateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) deleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.storage.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Alert handlers
func (h *Handler) listAlerts(c *gin.Context) {
	alerts, err := h.storage.ListAlerts(c.Request.Context(), 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list alerts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list alerts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": alerts})
}

func (h *Handler) createAlert(c *gin.Context) {
	var alert models.Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if alert.ID == uuid.Nil { alert.ID = uuid.New() }
	if err := h.storage.CreateAlert(c.Request.Context(), &alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert"})
		return
	}
	c.JSON(http.StatusCreated, alert)
}

func (h *Handler) getAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	alert, err := h.storage.GetAlertByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alert"})
		return
	}
	if alert == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}
	c.JSON(http.StatusOK, alert)
}

func (h *Handler) updateAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var alert models.Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	alert.ID = id
	if err := h.storage.UpdateAlert(c.Request.Context(), &alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alert"})
		return
	}
	c.JSON(http.StatusOK, alert)
}

func (h *Handler) deleteAlert(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.storage.DeleteAlert(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete alert"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Auth handlers (basic 501 placeholders)
func (h *Handler) register(c *gin.Context)    { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (h *Handler) login(c *gin.Context)       { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (h *Handler) refreshToken(c *gin.Context){ c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
