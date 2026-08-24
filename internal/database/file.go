package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UploadedFile struct {
	ID               int64     `json:"id"`
	OwnerID          int64     `json:"-"`
	OriginalFilename string    `json:"originalFilename"`
	StorageKey       string    `json:"-"`
	SizeBytes        int64     `json:"sizeBytes"`
	MimeType         string    `json:"mimeType"`
	UploadedAt       time.Time `json:"uploadedAt"`
	Links            []Link    `json:"links"`
}

type FileRepository struct {
	db *DB
}

func NewFileRepository(db *DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) Create(file UploadedFile) (*UploadedFile, error) {
	file.UploadedAt = time.Now().UTC()

	query := `
		INSERT INTO files (owner_id, original_filename, storage_key, size_bytes, mime_type, uploaded_at)
		VALUES (?, ?, ?, ?, ?, ?)
		`

	result, err := r.db.Conn().Exec(query,
		file.OwnerID,
		file.OriginalFilename,
		file.StorageKey,
		file.SizeBytes,
		file.MimeType,
		file.UploadedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting file: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting inserted file: %w", err)
	}

	file.ID = id
	file.Links = make([]Link, 0)
	return &file, nil
}

func (r *FileRepository) GetByID(id int64) (*UploadedFile, error) {
	query := `
		SELECT id, owner_id, original_filename, storage_key, size_bytes, mime_type, uploaded_at
		FROM files
		WHERE id = ?
	`

	var file UploadedFile
	if err := r.db.Conn().QueryRow(query, id).Scan(
		&file.ID,
		&file.OwnerID,
		&file.OriginalFilename,
		&file.StorageKey,
		&file.SizeBytes,
		&file.MimeType,
		&file.UploadedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting file: %w", err)
	}

	file.Links = make([]Link, 0)
	return &file, nil
}

func (r *FileRepository) DeleteByID(id int64) error {
	query := `
	DELTED FROM files
	WHERE id = ?
	`

	result, err := r.db.Conn().Exec(query, id)
	if err != nil {
		return fmt.Errorf("deleting file: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gettings rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("finding file: %w", err)
	}

	return nil
}

func (r *FileRepository) GetAllByOwnerID(ownerID int64) ([]UploadedFile, error) {
	query := `
		SELECT id, owner_id, original_filename, storage_key, size_bytes, mime_type, uploaded_at
		FROM files
		WHERE owner_id = ?
	`

	rows, err := r.db.Conn().Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("querying files: %w", err)
	}
	defer rows.Close()

	files := make([]UploadedFile, 0)
	for rows.Next() {
		var f UploadedFile
		if err := rows.Scan(
			&f.ID,
			&f.OwnerID,
			&f.OriginalFilename,
			&f.StorageKey,
			&f.SizeBytes,
			&f.MimeType,
			&f.UploadedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning file row: %w", err)
		}
		f.Links = make([]Link, 0)
		files = append(files, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating file rows: %w", err)
	}

	return files, nil
}
