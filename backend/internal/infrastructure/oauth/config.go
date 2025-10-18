package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// Config holds OAuth configuration
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// NewGitHubOAuthConfig creates a new GitHub OAuth2 config
func NewGitHubOAuthConfig(cfg Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}
