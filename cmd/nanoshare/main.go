package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("nanoshare started on addr :8080")
	srv.ListenAndServe()
}
