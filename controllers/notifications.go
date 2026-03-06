package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/cuappdev/chimes-backend/models"
	"github.com/cuappdev/chimes-backend/services"
)

// Struct for register token
type RegisterTokenInput struct {
    Token    string `json:"token" binding:"required"`
    Platform string `json:"platform" binding:"required,oneof=android ios"`
}

// POST /fcm/register
// Register an fcm token
func RegisterFCMToken(c *gin.Context) {
    var input RegisterTokenInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Get user ID from auth middleware
    userID := c.GetUint("userID")
    
    err := models.SaveOrUpdateToken(userID, input.Token, input.Platform)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Token registered successfully"})
}


// DELETE /fcm/delete
// Delete an fcm token
func DeleteFCMToken(c *gin.Context) {
    var input struct {
        Token string `json:"token" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    err := models.DeleteToken(input.Token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Token deleted successfully"})
}

// POST /fcm/test
// Send a test notification to the user
func SendTestNotification(c *gin.Context) {
    userID := c.GetUint("userID")
    
    payload := services.NotificationPayload{
        Title: "Test Notification",
        Body:  "This is a test notification",
        Data:  map[string]string{"type": "test"},
    }
    
    err := services.SendToUser(userID, payload)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}

// Backend only (not exposed to clients)
// Send a notification to a specific token
func SendNotificationToToken(c *gin.Context) {
	var input struct {
		Token   string `json:"token" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := services.NotificationPayload{
		Title: input.Title,
		Body:  input.Body,
	}

	err := services.SendToToken(input.Token, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}