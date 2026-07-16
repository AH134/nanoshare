package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AH134/nanoshare/internal/database"
)

func main() {
	db, err := database.New("./data/nanoshare.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SeedAdmin()

	router := http.NewServeMux()

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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("nanoshare started on addr :8080")
	srv.ListenAndServe()
}
