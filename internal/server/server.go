package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AH134/nanoshare/internal/config"
	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/handler"
	"github.com/AH134/nanoshare/internal/middleware"
	"github.com/alexedwards/scs/v2"
)

type Application struct {
	db             *database.DB
	userRepo       *database.UserRepository
	sessionManager *scs.SessionManager
	config         *config.EnvConfig
}

func NewApplication(db *database.DB, userRepo *database.UserRepository, sm *scs.SessionManager, cfg *config.EnvConfig) *Application {
	return &Application{
		db:             db,
		userRepo:       userRepo,
		sessionManager: sm,
		config:         cfg,
	}
}

func (a *Application) Mount() http.Handler {
	r := http.NewServeMux()

	// handlers
	authHandler := handler.NewAuthHandler(a.userRepo, a.sessionManager)

	// middlewares
	authMiddleware := middleware.NewAuthMiddleware(a.sessionManager)

	// routes
	r.HandleFunc("GET /api/health", handler.HealthCheck)

	r.HandleFunc("POST /api/auth/login", authHandler.Login)
	r.Handle("POST /api/auth/logout", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Logout)))

	return r
}

func (a *Application) Run(h http.Handler) error {
	addr := fmt.Sprintf(":%s", a.config.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      a.sessionManager.LoadAndSave(h),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("nanoshare has started at addr %s", addr)
	return srv.ListenAndServe()
}
