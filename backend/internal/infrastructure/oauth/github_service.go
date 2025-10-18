package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github-oauth-backend/internal/domain/model"
	"github-oauth-backend/internal/domain/service"
)

// GitHubServiceImpl implements the GitHubService interface
// This follows the Dependency Inversion Principle (SOLID)
type GitHubServiceImpl struct {
	httpClient *http.Client
}

// NewGitHubService creates a new GitHubServiceImpl
func NewGitHubService() service.GitHubService {
	return &GitHubServiceImpl{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetUserProfile fetches the GitHub user profile using an access token
func (s *GitHubServiceImpl) GetUserProfile(ctx context.Context, accessToken string) (*model.GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API returned status %d: %s", resp.StatusCode, string(body))
	}

	var githubUser struct {
		ID        int       `json:"id"`
		Login     string    `json:"login"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		AvatarURL string    `json:"avatar_url"`
		Bio       string    `json:"bio"`
		Company   string    `json:"company"`
		Location  string    `json:"location"`
		Blog      string    `json:"blog"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	user := model.NewGitHubUser(
		githubUser.ID,
		githubUser.Login,
		githubUser.Name,
		githubUser.Email,
		githubUser.AvatarURL,
		githubUser.Bio,
		githubUser.Company,
		githubUser.Location,
		githubUser.Blog,
		githubUser.CreatedAt,
		githubUser.UpdatedAt,
	)

	return &user, nil
}
