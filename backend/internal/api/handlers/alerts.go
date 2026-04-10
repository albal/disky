package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetAlerts returns a stubbed list of alerts for the user
func (h *Handlers) GetAlerts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"alerts": []map[string]interface{}{
		{
			"id": "alert-1",
			"type": "category_price_per_gb",
			"category": "m2-nvme",
			"threshold_price_per_gb": 0.05,
			"is_active": true,
		},
	}})
}

// CreateAlert stubs creating a new price alert
func (h *Handlers) CreateAlert(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Alert created"})
}
