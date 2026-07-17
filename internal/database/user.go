package database

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int
	Username     string
	PasswordHash string
	CreateAt     time.Time
}

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, password_hash, created_at FROM users WHERE username = ?"
	err := r.db.Conn().QueryRow(query, username).Scan(
		&user.Id,
		&user.Username,
		&user.PasswordHash,
		&user.CreateAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateAdmin(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}
	_, err = r.db.Conn().Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
	return err
}

func (r *UserRepository) SeedAdmin() error {
	var count int
	db := r.db.Conn()
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
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

	return r.CreateAdmin(username, password)
}
