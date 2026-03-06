package middleware

import (
	"context"
	"net/http"
	"strings"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/cuappdev/chimes-backend/auth"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const UIDKey ctxKey = "uid"

// RequireAuth validates either Firebase tokens or custom JWT tokens
func RequireAuth(firebaseAuthClient *firebaseauth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		const pref = "Bearer "
		authz := c.GetHeader("Authorization")
		if !strings.HasPrefix(authz, pref) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(authz, pref)

		// Try to validate as custom JWT token first
		jwtService := auth.NewJWTService()
		if claims, err := jwtService.ValidateToken(token); err == nil {
			// Custom JWT token is valid
			ctx := context.WithValue(c.Request.Context(), UIDKey, claims.UserID)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}

		// If custom JWT fails, try Firebase token
		firebaseToken, err := firebaseAuthClient.VerifyIDToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Firebase token is valid
		ctx := context.WithValue(c.Request.Context(), UIDKey, firebaseToken.UID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequireFirebaseUser validates only Firebase tokens (for backward compatibility)
func RequireFirebaseUser(ac *firebaseauth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		const pref = "Bearer "
		authz := c.GetHeader("Authorization")
		if !strings.HasPrefix(authz, pref) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		idToken := strings.TrimPrefix(authz, pref)

		tok, err := ac.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// put uid into context for handlers
		ctx := context.WithValue(c.Request.Context(), UIDKey, tok.UID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Helper to read the uid in handlers:
func UIDFrom(c *gin.Context) string {
	v := c.Request.Context().Value(UIDKey)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
