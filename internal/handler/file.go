package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/response"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/AH134/nanoshare/internal/storage"
	"github.com/AH134/nanoshare/internal/token"
	"github.com/alexedwards/scs/v2"
)

const DefaultFileKey = "file"

type FileHandler struct {
	files          *database.FileRepository
	links          *database.LinkRepository
	sessionManager *scs.SessionManager
	storage        storage.Storage
	logger         *slog.Logger
}

func NewFileHandler(files *database.FileRepository, links *database.LinkRepository, sessionManager *scs.SessionManager, storage storage.Storage, logger *slog.Logger) *FileHandler {
	return &FileHandler{
		files:          files,
		links:          links,
		sessionManager: sessionManager,
		storage:        storage,
		logger:         logger.With("handler", "file"),
	}
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt64(r.Context(), session.DefaultUserIDKey)

	// limit memory usage to 32mb
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Warn("failed to parse multipart form", "error", err)
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "INVALID_MULTIPART_FORM",
			Message: "Failed to process the uploaded form data.",
		})
		return
	}

	file, fileHeader, err := r.FormFile(DefaultFileKey)
	if errors.Is(err, http.ErrMissingFile) {
		h.logger.Warn("missing file in upload request", "error", err)
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "MISSING_FILE",
			Message: "The required file parameter is missing from the request.",
		})
		return
	}
	if err != nil {
		h.logger.Warn("invalid multipart form", "error", err)
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "INVALID_MULTIPART_FORM",
			Message: "Failed to process the uploaded form data.",
		})
		return
	}
	defer file.Close()

	storageKey, err := token.Generate(token.DefaultLength)
	if err != nil {
		response.InternalError(w, h.logger, "failed to generate token for file", err)
		return
	}

	if err := h.storage.Save(r.Context(), storageKey, file); err != nil {
		h.logger.Error("failed to save file to storage", "storage_key", storageKey, "error", err)
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An error occurred while saving your uploaded file(s).",
		})
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	uploadedFile, err := h.files.Create(database.UploadedFile{
		OwnerID:          userID,
		OriginalFilename: fileHeader.Filename,
		StorageKey:       storageKey,
		SizeBytes:        fileHeader.Size,
		MimeType:         mimeType,
	})
	if err != nil {
		// remove from storage if failed to save to db
		if delErr := h.storage.Delete(r.Context(), storageKey); delErr != nil {
			h.logger.Error("failed to clean up orphaned file", "storage_key", storageKey, "file", file, "error", err)
		}

		response.InternalError(w, h.logger, "failed to save file to database", err)
		return
	}

	response.Success(w, http.StatusCreated, uploadedFile)
}

func (h *FileHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt64(r.Context(), session.DefaultUserIDKey)

	files, err := h.files.GetAllByOwnerID(userID)
	if err != nil {
		h.logger.Error("failed to get files by owner ID", "user_id", userID, "error", err)
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to fetch files.",
		})
		return
	}

	fileIDs := make([]int64, len(files))
	fileByID := make(map[int64]*database.UploadedFile, len(files))
	for i := range files {
		fileIDs[i] = files[i].ID
		fileByID[files[i].ID] = &files[i]
	}

	links, err := h.links.GetAllByFileIDs(fileIDs)
	if err != nil {
		h.logger.Error("failed to get links by fileID", "file_id", userID, "error", err)
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to fetch links.",
		})
		return
	}

	for _, link := range links {
		if file, ok := fileByID[link.FileID]; ok {
			file.Links = append(file.Links, link)
		}
	}

	response.Success(w, http.StatusOK, files)
}
