package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /healthcheck
// Get healthcheck
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Chimes is healthy!"})
}
