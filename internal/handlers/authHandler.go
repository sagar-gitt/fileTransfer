package handlers

import (
	"context"
	"fileTransfer/internal/config"
	"fileTransfer/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
)

var googleOauthConfig *oauth2.Config

func InitGoogleAuth() {
	log.Printf("Initializing Google OAuth configuration...")
	log.Printf("Redirect URL: %s", config.GoogleConfig.RedirectURL)
	
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  config.GoogleConfig.RedirectURL,
		ClientID:     config.GoogleConfig.ClientID,
		ClientSecret: config.GoogleConfig.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	
	log.Printf("Google OAuth configuration initialized successfully")
}

func (h *Handlers) GoogleLogin(c *gin.Context) {
	log.Printf("Starting Google login process...")
	url := googleOauthConfig.AuthCodeURL("random") // In production, generate a secure state value
	log.Printf("Redirecting to Google auth URL: %s", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handlers) GoogleCallback(c *gin.Context) {
	log.Printf("Received Google callback request")
	
	state := c.Query("state")
	log.Printf("State parameter: %s", state)
	
	if state != "random" {
		log.Printf("Invalid state parameter: %s", state)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	log.Printf("Received auth code: %s", code)
	
	if code == "" {
		log.Printf("No code provided in callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}
	log.Printf("Successfully exchanged code for token")

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()
	log.Printf("Successfully retrieved user info")

	user, err := models.ParseGoogleUser(resp.Body)
	if err != nil {
		log.Printf("Failed to parse user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user"})
		return
	}
	log.Printf("Successfully parsed user info for: %s", user.Email)

	dbUser, err := h.UserDbRepo.FindOrCreateUser(user)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	log.Printf("Successfully found/created user in database")

	jwtToken, err := h.JWT.GenerateToken(dbUser.Email)
	if err != nil {
		log.Printf("Failed to generate JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}
	log.Printf("Successfully generated JWT token")

	// Redirect to frontend with token
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, jwtToken))
}
