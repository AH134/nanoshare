package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/AH134/nanoshare/internal/storage"
	"github.com/alexedwards/scs/v2"
)

const (
	STORAGE_KEY_LENGTH = 32
	REQ_FILE_KEY       = "file"
)

type StorageHandler struct {
	files          *database.FileRepository
	sessionManager *scs.SessionManager
	storage        storage.Storage
}

func NewStorageHandler(files *database.FileRepository, sessionManager *scs.SessionManager, storage storage.Storage) *StorageHandler {
	return &StorageHandler{
		files:          files,
		sessionManager: sessionManager,
		storage:        storage,
	}
}

func generateKey(length int) (string, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}

	return base64.URLEncoding.EncodeToString(key), nil
}

func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), session.UserIDKey)

	// limit memory usage to 32mb
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile(REQ_FILE_KEY)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	storageKey, err := generateKey(STORAGE_KEY_LENGTH)
	if err != nil {
		http.Error(w, "failed to upload file", http.StatusInternalServerError)
		return
	}

	if err := h.storage.Save(r.Context(), storageKey, file); err != nil {
		http.Error(w, "failed to upload file", http.StatusInternalServerError)
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	uploadedFile := &database.UploadedFile{
		OwnerID:          userID,
		OriginalFilename: fileHeader.Filename,
		StorageKey:       storageKey,
		SizeBytes:        fileHeader.Size,
		MimeType:         mimeType,
		UploadedAt:       time.Now(),
	}
	if err := h.files.Create(uploadedFile); err != nil {
		// remove from storage if failed to save to db
		if delErr := h.storage.Delete(r.Context(), storageKey); delErr != nil {
			log.Printf("upload: failed to clean up orphaned file %s: %v", storageKey, delErr)
		}
		http.Error(w, "failed to upload file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
