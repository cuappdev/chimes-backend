package controllers

import (
"net/http"
"strings"

"github.com/gin-gonic/gin"
"github.com/cuappdev/chimes-backend/models"
"github.com/cuappdev/chimes-backend/middleware"
	"github.com/cuappdev/chimes-backend/auth"
	firebaseauth "firebase.google.com/go/v4/auth"
)

// GET /users
// Get all users
func FindUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// POST /users
// Create new user
func CreateUser(c *gin.Context) {
  // Validate input
  var input models.CreateUserInput
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  
  uid := middleware.UIDFrom(c)
  if uid == "" {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "midding firebase uid"})
    return
  }

  // Create user
  user := models.User{
    FirstName: input.FirstName, 
    LastName: input.LastName, 
    Email: input.Email,
    Firebase_UID: uid,
  }
  models.DB.Create(&user)

  c.JSON(http.StatusOK, gin.H{"data": user})
}

// VerifyTokenRequest represents the request body for token verification
type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// POST /api/verify-token
// Verify Firebase token and return custom JWT tokens
func VerifyToken(firebaseAuthClient *firebaseauth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid token"})
			return
		}

		// Verify Firebase token
		firebaseToken, err := firebaseAuthClient.VerifyIDToken(c.Request.Context(), req.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Firebase token"})
			return
		}

		// Extract user data from Firebase token
		claims := firebaseToken.Claims
		firebaseUID := firebaseToken.UID
		
		// Get user info from Firebase token claims
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)
		
		// Parse name into first and last name
		nameParts := strings.Fields(name)
		firstName := ""
		lastName := ""
		if len(nameParts) > 0 {
			firstName = nameParts[0]
		}
		if len(nameParts) > 1 {
			lastName = strings.Join(nameParts[1:], " ")
		}

		// Find or create user in database
		user, err := models.FindOrCreateUser(firebaseUID, email, firstName, lastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/find user"})
			return
		}

		// Generate JWT tokens
		jwtService := auth.NewJWTService()
		tokenPair, err := jwtService.GenerateTokenPair(firebaseUID, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// Update user's refresh token in database
		if err := user.UpdateRefreshToken(tokenPair.RefreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
			return
		}

		// Return tokens and user info
		c.JSON(http.StatusOK, gin.H{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"expires_in":    tokenPair.ExpiresIn,
			"user": gin.H{
				"id":         user.ID,
				"firebase_uid": user.Firebase_UID,
				"email":      user.Email,
				"firstname":  user.FirstName,
				"lastname":   user.LastName,
			},
		})
	}
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// POST /api/refresh-token
// Refresh access token using refresh token
func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
			return
		}

		// Validate refresh token
		jwtService := auth.NewJWTService()
		userID, err := jwtService.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		// Find user in database
		var user models.User
		if err := models.DB.Where("firebase_uid = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Verify refresh token matches stored token
		if user.Refresh_Token != req.RefreshToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		// Generate new token pair
		tokenPair, err := jwtService.GenerateTokenPair(user.Firebase_UID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// Update user's refresh token
		if err := user.UpdateRefreshToken(tokenPair.RefreshToken); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update refresh token"})
			return
		}

		// Return new tokens
		c.JSON(http.StatusOK, gin.H{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"expires_in":    tokenPair.ExpiresIn,
		})
	}
}

