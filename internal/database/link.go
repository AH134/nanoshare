package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	ID            int                 `json:"id"`
	FileID        int                 `json:"fileID"`
	Token         string              `json:"token"`
	MaxDownloads  Nullable[int64]     `json:"maxDownloads"`
	DownloadCount int                 `json:"downloadCount"`
	CreatedAt     time.Time           `json:"createdAt"`
	ExpiresAt     Nullable[time.Time] `json:"expiresAt"`
	RevokedAt     Nullable[time.Time] `json:"revokedAt"`
}

type LinkRepository struct {
	db *DB
}

func NewLinkRepository(db *DB) *LinkRepository {
	return &LinkRepository{
		db: db,
	}
}

func (r *LinkRepository) Create(link *Link) error {
	query := "INSERT INTO links (file_id, token, max_downloads, expires_at) VALUES (?, ?, ?, ?)"

	_, err := r.db.Conn().Exec(query, link.FileID, link.Token, link.MaxDownloads, link.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to insert link: %w", err)
	}

	return nil
}

func (r *LinkRepository) GetByToken(token string) (*Link, error) {
	var link Link

	query := "SELECT id, file_id, token, max_downloads, download_count, created_at, expires_at, revoked_at FROM links WHERE token = ?"
	err := r.db.Conn().QueryRow(query, token).Scan(
		&link.ID,
		&link.FileID,
		&link.Token,
		&link.MaxDownloads,
		&link.DownloadCount,
		&link.CreatedAt,
		&link.ExpiresAt,
		&link.RevokedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed ot get link with token %s: %w", token, err)
	}
	return &link, nil
}

func (r *LinkRepository) Revoke(id int) error {
	query := "UPDATE links SET revoked_at = ? WHERE id = ?"
	_, err := r.db.Conn().Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update revoked_at attribute for link %d: %w", id, err)
	}
	return nil
}
