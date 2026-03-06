package controllers

import (
	"net/http"

	"github.com/cuappdev/chimes-backend/models"
	"github.com/gin-gonic/gin"
)

// GET /health
// Get healthcheck
func HealthCheck(c *gin.Context) {
	sqlDB, err := models.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "database": "disconnected"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "database": "disconnected"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "database": "connected"})
}
