package main

import (
	"context"
	"log"
	"os"

	"github.com/cuappdev/chimes-backend/auth"
	"github.com/cuappdev/chimes-backend/controllers"
	"github.com/cuappdev/chimes-backend/middleware"
	"github.com/cuappdev/chimes-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file.")
	}

	log.Println("Starting chimes-backend...")
	r := gin.Default()
	r.SetTrustedProxies(nil) // Trust no proxies for now, will configure later
	log.Println("Connecting to database...")
	// Connect to DB
	if err := models.ConnectDatabase(); err != nil {
		log.Printf("[FATAL] Database connection failed: %v", err)
	}

	// Migrate fcm token table
	models.DB.AutoMigrate(&models.FCMToken{})

	// Initialize Firebase Auth SAFELY
	serviceAccountPath := "service-account-key.json"
	// Log working dir and check file exists
	if _, err := os.Stat(serviceAccountPath); err != nil {
		log.Printf("[FATAL] Missing service account file: %s (cwd: %s): %v", serviceAccountPath, getwdSafe(), err)
	}
	ac, err := auth.NewAuthClient(context.Background(), serviceAccountPath)
	if err != nil {
		log.Printf("[FATAL] Firebase init failed: %v", err)
	}

	// Initialize Firebase Messaging
	if err := auth.InitFirebase(serviceAccountPath); err != nil {
		log.Printf("[FATAL] Firebase Messaging init failed: %v", err)
	}

	log.Println("Setting up routes...")
	// Public routes
	r.GET("/health", controllers.HealthCheck)

	// Auth routes (public)
	api := r.Group("/api")
	{
		api.POST("/verify-token", controllers.VerifyToken(ac))
		api.POST("/refresh-token", controllers.RefreshToken())
	}

	// Protected routes
	authd := api.Group("")
	authd.Use(middleware.RequireAuth(ac))
	{
		// User routes
		authd.GET("/users", controllers.FindUsers)
		authd.POST("/users", controllers.CreateUser)
		// Notification routes
		authd.POST("/fcm/register", controllers.RegisterFCMToken)
		authd.DELETE("/fcm/delete", controllers.DeleteFCMToken)
		authd.POST("/fcm/test", controllers.SendTestNotification)

	}
	log.Println("Server starting on :8080")

	r.Run()
}

func getwdSafe() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}
