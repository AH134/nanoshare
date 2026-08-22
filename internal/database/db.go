package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const dbPragmaOptions = "?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_time_format=sqlite&_timezone=UTC&_texttotime=1"

type DB struct {
	conn *sql.DB
}

func New(path string) (*DB, error) {
	dsn := path + dbPragmaOptions

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("pinging db: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (d *DB) Conn() *sql.DB {
	return d.conn
}

func (d *DB) Ping() error {
	return d.conn.Ping()
}

func (d *DB) Close() error {
	return d.conn.Close()
}
