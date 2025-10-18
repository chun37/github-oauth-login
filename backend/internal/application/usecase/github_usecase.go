package usecase

import (
	"context"
	"github-oauth-backend/internal/application/dto"
	"github-oauth-backend/internal/domain/service"
)

// GitHubUseCase handles GitHub-related use cases
// This follows the Single Responsibility Principle (SOLID)
type GitHubUseCase struct {
	githubService service.GitHubService
}

// NewGitHubUseCase creates a new GitHubUseCase
func NewGitHubUseCase(githubService service.GitHubService) *GitHubUseCase {
	return &GitHubUseCase{
		githubService: githubService,
	}
}

// GetUserProfile retrieves the GitHub user profile and converts it to a DTO
func (uc *GitHubUseCase) GetUserProfile(ctx context.Context, accessToken string) (*dto.GitHubUserResponse, error) {
	user, err := uc.githubService.GetUserProfile(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	// Convert domain model to DTO
	return &dto.GitHubUserResponse{
		ID:        user.ID,
		Login:     user.Login,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Company:   user.Company,
		Location:  user.Location,
		Blog:      user.Blog,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
