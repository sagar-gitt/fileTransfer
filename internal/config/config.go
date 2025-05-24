package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var GoogleConfig GoogleOAuthConfig

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Continue execution even if .env file is not found
	}

	// Log the environment variables (without sensitive data)
	log.Printf("Loading Google OAuth configuration...")
	log.Printf("GOOGLE_REDIRECT_URL: %s", os.Getenv("GOOGLE_REDIRECT_URL"))
	log.Printf("GOOGLE_CLIENT_ID exists: %v", os.Getenv("GOOGLE_CLIENT_ID") != "")
	log.Printf("GOOGLE_CLIENT_SECRET exists: %v", os.Getenv("GOOGLE_CLIENT_SECRET") != "")

	GoogleConfig = GoogleOAuthConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	}

	// Validate the configuration
	if GoogleConfig.ClientID == "" {
		log.Fatal("GOOGLE_CLIENT_ID is not set")
	}
	if GoogleConfig.ClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_SECRET is not set")
	}
	if GoogleConfig.RedirectURL == "" {
		log.Fatal("GOOGLE_REDIRECT_URL is not set")
	}

	log.Printf("Google OAuth configuration loaded successfully")
}
