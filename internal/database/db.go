package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn   *sql.DB
	dbPath string
}

func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	conn.SetMaxOpenConns(1)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("pinging db: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (d *DB) Ping() error {
	return d.conn.Ping()
}

func (d *DB) Close() error {
	return d.conn.Close()
}
