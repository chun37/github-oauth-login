package service

import (
	"context"
	"github-oauth-backend/internal/domain/model"
)

// GitHubService defines the domain service interface for GitHub operations
// This follows the Dependency Inversion Principle (SOLID)
type GitHubService interface {
	// GetUserProfile fetches the GitHub user profile using an access token
	GetUserProfile(ctx context.Context, accessToken string) (*model.GitHubUser, error)
}
