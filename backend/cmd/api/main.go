package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github-oauth-backend/internal/application/usecase"
	"github-oauth-backend/internal/infrastructure/database"
	"github-oauth-backend/internal/infrastructure/oauth"
	"github-oauth-backend/internal/infrastructure/session"
	"github-oauth-backend/internal/interfaces/handler"
	customMiddleware "github-oauth-backend/internal/interfaces/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	// Load environment variables
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "github_oauth_app"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	oauthCfg := oauth.Config{
		ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		RedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/api/auth/callback"),
	}

	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	port := getEnv("BACKEND_PORT", "8080")

	// Initialize database connection
	pool, err := database.NewPostgresPool(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Database connection established")

	// Initialize session manager
	sessionManager := session.NewSessionManager(pool)

	// Initialize OAuth config
	oauthConfig := oauth.NewGitHubOAuthConfig(oauthCfg)

	// Initialize services
	githubService := oauth.NewGitHubService()

	// Initialize use cases
	githubUseCase := usecase.NewGitHubUseCase(githubService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(oauthConfig, sessionManager)
	userHandler := handler.NewUserHandler(githubUseCase, sessionManager)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(customMiddleware.CORSConfig(frontendURL))
	e.Use(echo.WrapMiddleware(sessionManager.LoadAndSave))

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Auth routes
	e.GET("/api/auth/login", authHandler.Login)
	e.GET("/api/auth/callback", authHandler.Callback)
	e.GET("/api/auth/check", userHandler.CheckAuth)

	// User routes
	e.GET("/api/user/profile", userHandler.GetProfile)

	// Start server
	address := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", address)
	if err := e.Start(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv retrieves environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
