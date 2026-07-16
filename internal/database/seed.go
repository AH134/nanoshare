package database

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func (d *DB) CreateAdmin(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}
	_, err = d.conn.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
	return err
}

func (d *DB) SeedAdmin() error {
	var count int
	if err := d.conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("no admin exists and ADMIN_USERNAME/ADMIN_PASSWORD not set")
	}

	return d.CreateAdmin(username, password)
}
