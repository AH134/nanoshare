package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Nullable[T any] struct {
	sql.Null[T]
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(n.V)
}

type Link struct {
	ID            int64               `json:"id"`
	FileID        int64               `json:"fileID"`
	Token         string              `json:"token"`
	MaxDownloads  Nullable[int64]     `json:"maxDownloads"`
	DownloadCount int64               `json:"downloadCount"`
	CreatedAt     time.Time           `json:"createdAt"`
	ExpiresAt     Nullable[time.Time] `json:"expiresAt"`
	RevokedAt     Nullable[time.Time] `json:"revokedAt"`
}

func (l *Link) IsValid() bool {
	revokedAt := l.RevokedAt
	expiresAt := l.ExpiresAt
	maxDownloads := l.MaxDownloads
	downloadCount := l.DownloadCount

	if revokedAt.Valid {
		return false
	}

	if expiresAt.Valid && expiresAt.V.Before(time.Now()) {
		return false
	}

	if maxDownloads.Valid && maxDownloads.V <= downloadCount {
		return false
	}

	return true
}

type LinkRepository struct {
	db *DB
}

func NewLinkRepository(db *DB) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) Create(link Link) (*Link, error) {
	link.CreatedAt = time.Now().UTC()

	query := `
		INSERT INTO links (file_id, token, max_downloads, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Conn().Exec(query,
		link.FileID,
		link.Token,
		link.MaxDownloads,
		link.CreatedAt,
		link.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert link: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting inserted link: %w", err)
	}

	link.ID = id
	return &link, nil
}

func (r *LinkRepository) GetAllByFileIDs(fileIDs []int64) ([]Link, error) {
	if len(fileIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(fileIDs))
	args := make([]any, len(fileIDs))
	for i, id := range fileIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	inClause := strings.Join(placeholders, ", ")
	query := `
		SELECT id, file_id, token, max_downloads, download_count, created_at, expires_at, revoked_at
		FROM links
		WHERE file_id IN (` + inClause + `)
		ORDER BY created_at ASC
	`

	rows, err := r.db.Conn().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying links: %w", err)
	}
	defer rows.Close()

	links := make([]Link, 0)
	for rows.Next() {
		var l Link
		if err := rows.Scan(
			&l.ID,
			&l.FileID,
			&l.Token,
			&l.MaxDownloads,
			&l.DownloadCount,
			&l.CreatedAt,
			&l.ExpiresAt,
			&l.RevokedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning link row %w", err)
		}
		links = append(links, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating link rows: %w", err)
	}

	return links, nil
}

func (r *LinkRepository) GetByToken(token string) (*Link, error) {
	query := `
		SELECT id, file_id, token, max_downloads, download_count, created_at, expires_at, revoked_at
		FROM links
		WHERE token = ?
	`

	var link Link
	if err := r.db.Conn().QueryRow(query, token).Scan(
		&link.ID,
		&link.FileID,
		&link.Token,
		&link.MaxDownloads,
		&link.DownloadCount,
		&link.CreatedAt,
		&link.ExpiresAt,
		&link.RevokedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("gettin link with token %s: %w", token, err)
	}

	return &link, nil
}

func (r *LinkRepository) Revoke(id int64) error {
	query := `
		UPDATE links
		SET revoked_at = ?
		WHERE id = ?
	`

	_, err := r.db.Conn().Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("setting revoked_at attribute for link %d: %w", id, err)
	}

	return nil
}

func (r *LinkRepository) IncrementDownloadCount(id int64) error {
	query := `
		UPDATE links
		SET download_count = download_count + 1
		WHERE id = ?
	`

	_, err := r.db.Conn().Exec(query, id)
	if err != nil {
		return fmt.Errorf("setting download_count attribute for link %d: %w", id, err)
	}

	return nil
}
