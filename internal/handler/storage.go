package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/response"
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
	userID := h.sessionManager.GetInt64(r.Context(), session.DefaultUserIDKey)

	// limit memory usage to 32mb
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "INVALID_MULTIPART_FORM",
			Message: "Failed to process the uploaded form data.",
		})
	}

	file, fileHeader, err := r.FormFile(DefaultFileKey)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			response.Error(w, http.StatusBadRequest, response.APIError{
				Code:    "MISSING_FILE",
				Message: "The required file parameter is missing from the request.",
			})
			return
		}

		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "INVALID_MULTIPART_FORM",
			Message: "Failed to process the uploaded form data.",
		})
	}
	defer file.Close()

	storageKey, err := token.Generate(token.DefaultLength)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An error occured while processing your uploaded file(s).",
		})
		return
	}

	if err := h.storage.Save(r.Context(), storageKey, file); err != nil {
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An error occured while saving your uploaded file(s).",
		})
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	createdFile, err := h.files.Create(
		database.UploadedFile{
			OwnerID:          userID,
			OriginalFilename: fileHeader.Filename,
			StorageKey:       storageKey,
			SizeBytes:        fileHeader.Size,
			MimeType:         mimeType,
		})

	if err != nil {
		// remove from storage if failed to save to db
		if delErr := h.storage.Delete(r.Context(), storageKey); delErr != nil {
			log.Printf("upload: failed to clean up orphaned file %s: %v", storageKey, delErr)
		}

		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to upload file.",
		})
		return
	}

	response.Success(w, http.StatusCreated, createdFile)
}

func (h *StorageHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt(r.Context(), session.DefaultUserIDKey)

	files, err := h.files.GetAllByOwnerID(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to fetch files.",
		})
		return
	}

	if files == nil {
		files = make([]*database.UploadedFile, 0)
	}

	response.Success(w, http.StatusOK, files)
}
