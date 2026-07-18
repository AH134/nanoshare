package database

import (
	"fmt"
	"time"
)

type UploadedFile struct {
	ID               int       `json:"id"`
	OwnerID          int       `json:"-"`
	OriginalFilename string    `json:"originalFilename"`
	StorageKey       string    `json:"-"`
	SizeBytes        int64     `json:"sizeBytes"`
	MimeType         string    `json:"mimeType"`
	UploadedAt       time.Time `json:"uploadedAt"`
}

type FileRepository struct {
	db *DB
}

func NewFileRepository(db *DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) Create(file *UploadedFile) error {
	query := "INSERT INTO files (owner_id, original_filename, storage_key, size_bytes, mime_type, uploaded_at) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := r.db.Conn().Exec(query, file.OwnerID, file.OriginalFilename, file.StorageKey, file.SizeBytes, file.MimeType, file.UploadedAt)
	if err != nil {
		return fmt.Errorf("failed to insert uploaded file: %w", err)
	}

	return nil
}

func (r *FileRepository) GetByID(id int) (*UploadedFile, error) {
	var file UploadedFile
	query := "SELECT id, owner_id, original_filename, storage_key, size_bytes, mime_type, uploaded_at FROM files WHERE id = ?"

	err := r.db.Conn().QueryRow(query, id).Scan(
		&file.ID,
		&file.OwnerID,
		&file.OriginalFilename,
		&file.StorageKey,
		&file.SizeBytes,
		&file.MimeType,
		&file.UploadedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get file with id %d: %w", id, err)
	}

	return &file, nil
}

func (r *FileRepository) GetAllByOwnerID(ownerID int) ([]*UploadedFile, error) {
	query := "SELECT id, owner_id, original_filename, storage_key, size_bytes, mime_type, uploaded_at FROM files WHERE owner_id = ?"

	rows, err := r.db.Conn().Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch files for owner %d: %w", ownerID, err)
	}
	defer rows.Close()

	var files []*UploadedFile
	for rows.Next() {
		var file UploadedFile
		err := rows.Scan(
			&file.ID,
			&file.OwnerID,
			&file.OriginalFilename,
			&file.StorageKey,
			&file.SizeBytes,
			&file.MimeType,
			&file.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file row: %w", err)
		}

		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate through file rows: %w", err)
	}

	return files, nil
}
