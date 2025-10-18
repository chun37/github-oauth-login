package handler

import (
	"net/http"

	"github-oauth-backend/internal/application/usecase"
	"github-oauth-backend/internal/infrastructure/session"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

// UserHandler handles user-related requests
type UserHandler struct {
	githubUseCase  *usecase.GitHubUseCase
	sessionManager *scs.SessionManager
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(githubUseCase *usecase.GitHubUseCase, sessionManager *scs.SessionManager) *UserHandler {
	return &UserHandler{
		githubUseCase:  githubUseCase,
		sessionManager: sessionManager,
	}
}

// GetProfile retrieves the user's GitHub profile
func (h *UserHandler) GetProfile(c echo.Context) error {
	ctx := c.Request().Context()

	// Get access token from session
	accessToken := h.sessionManager.GetString(ctx, session.AccessTokenKey)
	if accessToken == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Not authenticated",
		})
	}

	// Fetch user profile from GitHub
	profile, err := h.githubUseCase.GetUserProfile(ctx, accessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch profile",
		})
	}

	return c.JSON(http.StatusOK, profile)
}

// CheckAuth checks if the user is authenticated
func (h *UserHandler) CheckAuth(c echo.Context) error {
	ctx := c.Request().Context()

	accessToken := h.sessionManager.GetString(ctx, session.AccessTokenKey)
	if accessToken == "" {
		return c.JSON(http.StatusOK, map[string]bool{
			"authenticated": false,
		})
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"authenticated": true,
	})
}
