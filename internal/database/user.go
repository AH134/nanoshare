package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
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
	query := `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = ?
	`

	var user User
	err := r.db.Conn().QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int64) (*User, error) {
	query := `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := r.db.Conn().QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdatePassword(id int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = ?
		WHERE id = ?
	`
	result, err := r.db.Conn().Exec(query, passwordHash, id)
	if err != nil {
		return fmt.Errorf("changing password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("finding user: %w", err)
	}

	return nil
}

func (r *UserRepository) CreateAdmin(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	query := `
		INSERT INTO users (username, password_hash, created_at)
		VALUES (?, ?, ?)
	`

	_, err = r.db.Conn().Exec(query, username, string(hash), time.Now().UTC())
	if err != nil {
		return fmt.Errorf("inserting admin user: %w", err)
	}

	return nil
}

func (r *UserRepository) SeedAdmin() error {
	query := `
		SELECT COUNT(*)
		FROM users
	`

	var count int
	if err := r.db.Conn().QueryRow(query).Scan(&count); err != nil {
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
