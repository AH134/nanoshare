package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/AH134/nanoshare/internal/storage"
	"github.com/AH134/nanoshare/internal/token"
	"github.com/alexedwards/scs/v2"
)

const DefaultFileKey = "file"

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

func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), session.DefaultUserIDKey)

	// limit memory usage to 32mb
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile(DefaultFileKey)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	storageKey, err := token.Generate(token.DefaultLength)
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

func (h *StorageHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), session.DefaultUserIDKey)

	files, err := h.files.GetAllByOwnerID(userID)
	if err != nil {
		http.Error(w, "failed to fetch all files", http.StatusInternalServerError)
		return
	}

	if files == nil {
		files = make([]*database.UploadedFile, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "failed to fetch all files", http.StatusInternalServerError)
		return
	}
}
