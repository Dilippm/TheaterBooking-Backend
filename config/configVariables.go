package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Exported environment variables
var (
	ENV       string
	PORT      string
	MONGO_URI string
	JWT_SECRET_KEY string
	STRIPE_KEY string
)

func init() {
	// Load environment variables from .env file
	if os.Getenv("ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file; using environment variables")
		}

	// Set environment variables
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8023" // Default port (for local development, for example)
	}
	MONGO_URI = os.Getenv("MONGODB_URI")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	STRIPE_KEY = os.Getenv("STRIPE_KEY")
	}
}