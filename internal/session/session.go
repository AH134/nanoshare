package session

import (
	"database/sql"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

func New(db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()

	sessionManager.Store = sqlite3store.New(db)
	// sessionManager.Lifetime = 30 * time.Second
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute

	return sessionManager
}
