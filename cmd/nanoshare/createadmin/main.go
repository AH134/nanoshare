package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/AH134/nanoshare/internal/database"
)

func main() {
	username := flag.String("username", "admin", "admin username")
	password := flag.String("password", "admin", "admin password")
	flag.Parse()

	db, err := database.New("./data/nanoshare.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.CreateAdmin(*username, *password); err != nil {
		log.Fatal(err)
	}

	fmt.Println("admin user created")
}
