package main

import (
	"log"

	"github.com/AH134/nanoshare/internal/config"
	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/server"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/AH134/nanoshare/internal/storage"
)

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

	if err := database.RunMigrations(db.Conn()); err != nil {
		log.Fatal(err)
	}

	userRepo := database.NewUserRepository(db)
	if err := userRepo.SeedAdmin(); err != nil {
		log.Fatal(err)
	}

	sessionManager := session.New(db.Conn(), cfg.Prod)

	app := server.NewApplication(db, userRepo, sessionManager, storage, cfg)

	if err := app.Run(app.Mount()); err != nil {
		log.Fatalf("nanoshare has failed to start: %s", err)
	}
}
