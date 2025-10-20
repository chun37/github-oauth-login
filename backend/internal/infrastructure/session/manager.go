package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewSessionManager creates a new session manager with PostgreSQL store
// Session lifetime is set to 1 year (365 days) as per requirements
func NewSessionManager(pool *pgxpool.Pool, cookieDomain string) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.NewWithCleanupInterval(pool, 30*time.Minute)

	// Set session lifetime to 1 year
	sessionManager.Lifetime = 365 * 24 * time.Hour

	// Cookie settings
	sessionManager.Cookie.Name = "session_id"
	sessionManager.Cookie.Domain = cookieDomain // Set domain from environment variable
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = false                      // Set to false for HTTP environments
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode // SameSite=Lax (appropriate for HTTP environments)
	sessionManager.Cookie.Persist = true

	return sessionManager
}

// SessionKeys defines the keys used in session storage
const (
	AccessTokenKey = "access_token"
	StateKey       = "oauth_state"
)
