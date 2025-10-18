package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github-oauth-backend/internal/infrastructure/session"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// AuthHandler handles OAuth authentication
type AuthHandler struct {
	oauthConfig    *oauth2.Config
	sessionManager *scs.SessionManager
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(oauthConfig *oauth2.Config, sessionManager *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		oauthConfig:    oauthConfig,
		sessionManager: sessionManager,
	}
}

// Login initiates the GitHub OAuth flow
func (h *AuthHandler) Login(c echo.Context) error {
	// Generate a random state token
	state, err := generateStateToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate state token",
		})
	}

	// Store state in session
	h.sessionManager.Put(c.Request().Context(), session.StateKey, state)

	// Redirect to GitHub authorization URL
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles the OAuth callback from GitHub
func (h *AuthHandler) Callback(c echo.Context) error {
	ctx := c.Request().Context()

	// Verify state token
	storedState := h.sessionManager.GetString(ctx, session.StateKey)
	receivedState := c.QueryParam("state")

	if storedState == "" || storedState != receivedState {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid state token",
		})
	}

	// Remove state from session
	h.sessionManager.Remove(ctx, session.StateKey)

	// Exchange authorization code for access token
	code := c.QueryParam("code")
	token, err := h.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to exchange token",
		})
	}

	// Store access token in session
	h.sessionManager.Put(ctx, session.AccessTokenKey, token.AccessToken)

	// Redirect to frontend
	frontendURL := c.Request().Header.Get("X-Frontend-URL")
	if frontendURL == "" {
		frontendURL = "http://127.0.0.1:3000"
	}

	return c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/profile")
}

// generateStateToken generates a random state token for OAuth (alphanumeric only)
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Use hex encoding to ensure only alphanumeric characters (0-9, a-f)
	return hex.EncodeToString(b), nil
}
