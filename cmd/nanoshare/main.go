package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/handlers"
	"github.com/AH134/nanoshare/internal/middleware"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/alexedwards/scs/sqlite3store"
)

func main() {
	db, err := database.New("./data/nanoshare.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sessionManager := session.New(db.Conn())
	sqlite3store.NewWithCleanupInterval(db.Conn(), 5*time.Minute)

	router := http.NewServeMux()

	userRepo := database.NewUserRepository(db)
	userRepo.SeedAdmin()
	authHandler := handlers.NewAuthHandler(userRepo, sessionManager)
	authMW := middleware.NewAuthMiddleware(sessionManager)

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`{"status":"ok"}`))
	})

	router.HandleFunc("GET /health/db", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"status": "ok"}`))
	})

	router.HandleFunc("POST /api/login", authHandler.Login)
	router.Handle("POST /api/logout", authMW.RequireAuth(http.HandlerFunc(authHandler.Logout)))
	router.Handle("GET /me", authMW.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "authorized"}`))
	})))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: sessionManager.LoadAndSave(router),
	}

	fmt.Println("nanoshare started on addr :8080")
	srv.ListenAndServe()
}
