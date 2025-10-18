package dto

import "time"

// GitHubUserResponse represents the response DTO for GitHub user data
// This separates domain models from API responses (Clean Architecture principle)
type GitHubUserResponse struct {
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
