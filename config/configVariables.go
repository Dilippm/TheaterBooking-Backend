package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// Exported environment variables
var (
    ENV            string
    PORT           string
    MONGO_URI      string
    JWT_SECRET_KEY string
    STRIPE_KEY     string
)

func init() {
    // Set environment variables from the environment (Render will provide these in production)
    ENV = os.Getenv("ENV")
    
    // Load .env only if in development mode
    if ENV == "development" {
        err := godotenv.Load() // Load the .env file only in development
        if err != nil {
            log.Println("Error loading .env file; using environment variables")
        }
    }

    // Set environment variables
    PORT = os.Getenv("PORT")
    if PORT == "" {
        PORT = "8023" // Default port (for local development)
    }

    MONGO_URI = os.Getenv("MONGODB_URI")
    JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
    STRIPE_KEY = os.Getenv("STRIPE_KEY")

    // Log the loaded variables for debugging purposes (optional)
    log.Println("Loaded environment variables:")
    log.Println("PORT:", PORT)
    log.Println("MONGO_URI:", MONGO_URI)
}
