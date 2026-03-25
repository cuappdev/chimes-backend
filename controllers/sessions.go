package controllers

import (
	"net/http"

	"github.com/cuappdev/chimes-backend/models"
	"github.com/gin-gonic/gin"
)

// GET /sessions
// Get all sessions
func FindSessions(c *gin.Context) {
	var sessions []models.Session
	models.DB.Find(&sessions)

	c.JSON(http.StatusOK, gin.H{"data": sessions})
}

// POST /sessions
// Create a new session
func CreateSession(c *gin.Context) {
	var input models.SessionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := models.Session{
		Day:        input.Day,
		TimeOfDay:  input.TimeOfDay,
		KudosCount: 0,
	}

	if err := models.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": session})
}
