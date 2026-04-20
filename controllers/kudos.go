package controllers

import (
	"net/http"

	"github.com/cuappdev/chimes-backend/middleware"
	"github.com/cuappdev/chimes-backend/models"
	"github.com/gin-gonic/gin"
)

// POST/api/sessions/:id/kudos
// Student sends a kudo for a session

// POST /api/kudos
// Student sends a kudo for a session
func CreateKudo(c *gin.Context) {
	var input models.KudoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := middleware.UIDFrom(c)
	var user models.User
	if err := models.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	kudo := models.Kudo{
		SessionID:  input.SessionID,
		KudoTypeID: input.KudoTypeID,
		UserID:     user.ID,
	}

	if err := models.DB.Create(&kudo).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "kudo already sent"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": kudo})
}

// GET/api/kudos
// Admin-only route to see kudo counts grouped by type for a session

func GetSessionKudos(c *gin.Context) {
	var req models.KudoSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	counts, err := models.GetKudoCountsForSession(req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch kudos"})
	}
	c.JSON(http.StatusOK, gin.H{"data": counts})
}
