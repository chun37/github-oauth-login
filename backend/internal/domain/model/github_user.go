package model

import "time"

// GitHubUser represents a GitHub user profile (Value Object)
// This is immutable and serves as a value object in DDD
type GitHubUser struct {
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

// NewGitHubUser creates a new GitHubUser value object
func NewGitHubUser(
	id int,
	login, name, email, avatarURL, bio, company, location, blog string,
	createdAt, updatedAt time.Time,
) GitHubUser {
	return GitHubUser{
		ID:        id,
		Login:     login,
		Name:      name,
		Email:     email,
		AvatarURL: avatarURL,
		Bio:       bio,
		Company:   company,
		Location:  location,
		Blog:      blog,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
