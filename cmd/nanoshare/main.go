package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AH134/nanoshare/internal/config"
	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/server"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/AH134/nanoshare/internal/storage"
	"github.com/alexedwards/scs/sqlite3store"
)

func temp(db *database.DB) {
	linkRepo := database.NewLinkRepository(db)
	// link := &database.Link{
	// 	FileID:    1,
	// 	Token:     "test-token-123",
	// 	ExpiresAt: time.Now().Add(24 * time.Hour),
	// }

	// if err := linkRepo.Create(link); err != nil {
	// 	log.Fatal(err)
	// }

	l, err := linkRepo.GetByToken("test-token-123")
	if err != nil {
		log.Fatal(err)
	}
	j, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", l)
	fmt.Printf("%s\n", j)

}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewLocalStorage(cfg.StoragePath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg.DbPath)
	if err != nil {
		log.Fatal(err)
	}

	// temp(db)

	if err := database.RunMigrations(db.Conn()); err != nil {
		log.Fatal(err)
	}

	userRepo := database.NewUserRepository(db)
	if err := userRepo.SeedAdmin(); err != nil {
		log.Fatal(err)
	}

	sessionManager := session.New(db.Conn())
	sqlite3store.NewWithCleanupInterval(db.Conn(), 10*time.Minute)

	app := server.NewApplication(db, userRepo, sessionManager, storage, cfg)

	if err := app.Run(app.Mount()); err != nil {
		log.Fatalf("nanoshare has failed to start: %s", err)
	}
}
