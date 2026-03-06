package models

import (
    "fmt"
    "log"
    "os"
    "time"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

var DB *gorm.DB

func ConnectDatabase() error {
    // Build connection string from env variables
    host := getEnv("DB_HOST", "localhost")
    user := getEnv("DB_USER", "postgres")
    dbname := getEnv("DB_NAME", "chimes")
    port := getEnv("DB_PORT", "5432")
    password := getEnv("DB_PASSWORD", "")
    sslmode := getEnv("DB_SSLMODE", "require")
    
    dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    sqlDB, err := database.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %w", err)
    }
    
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    // Make sure to include all models to migrate here
    err = database.AutoMigrate(&User{}, &Seller{}, &FCMToken{}, &ServiceListing{}, &Service{})
    if err != nil {
        return fmt.Errorf("failed to migrate database: %w", err)
    }

    DB = database
    log.Println("Database connected successfully")
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}