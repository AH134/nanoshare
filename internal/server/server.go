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
	"github.com/AH134/nanoshare/internal/storage"
	"github.com/alexedwards/scs/v2"
)

type Application struct {
	db             *database.DB
	userRepo       *database.UserRepository
	sessionManager *scs.SessionManager
	storage        storage.Storage
	config         *config.EnvConfig
}

func NewApplication(db *database.DB, userRepo *database.UserRepository, sm *scs.SessionManager, storage storage.Storage, cfg *config.EnvConfig) *Application {
	return &Application{
		db:             db,
		userRepo:       userRepo,
		sessionManager: sm,
		storage:        storage,
		config:         cfg,
	}
}

func (a *Application) Mount() http.Handler {
	r := http.NewServeMux()

	// db repositories
	fileRepo := database.NewFileRepository(a.db)
	linkRepo := database.NewLinkRepository(a.db)

	// handlers
	authHandler := handler.NewAuthHandler(a.userRepo, a.sessionManager)
	storageHandler := handler.NewStorageHandler(fileRepo, a.sessionManager, a.storage)
	linkHandler := handler.NewLinkHandler(linkRepo, fileRepo, a.sessionManager, a.storage)

	// middlewares
	mwChain := middleware.Chain(
		middleware.Logging,
		a.sessionManager.LoadAndSave,
	)
	requireAuth := middleware.RequireAuth(a.sessionManager)

	// routes
	r.HandleFunc("GET /api/health", handler.HealthCheck)
	r.Handle("GET /api/me", requireAuth(http.HandlerFunc(authHandler.Me)))

	r.HandleFunc("POST /api/auth/login", authHandler.Login)
	r.Handle("POST /api/auth/logout", requireAuth(http.HandlerFunc(authHandler.Logout)))

	r.Handle("GET /api/files", requireAuth(http.HandlerFunc(storageHandler.ListFiles)))
	r.Handle("POST /api/files", requireAuth(http.HandlerFunc(storageHandler.Upload)))

	r.Handle("POST /api/files/{id}/links", requireAuth(http.HandlerFunc(linkHandler.Create)))

	r.HandleFunc("GET /d/{token}", linkHandler.Download)

	return mwChain(r)
}

func (a *Application) Run(h http.Handler) error {
	addr := fmt.Sprintf(":%s", a.config.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("nanoshare has started at addr %s", addr)
	return srv.ListenAndServe()
}
