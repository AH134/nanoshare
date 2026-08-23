package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/response"
	"github.com/AH134/nanoshare/internal/storage"
	"github.com/AH134/nanoshare/internal/token"
	"github.com/alexedwards/scs/v2"
)

type LinkRequest struct {
	MaxDownloads *int64     `json:"maxDownloads"`
	ExpiresAt    *time.Time `json:"expiresAt"`
}

type LinkHandler struct {
	links          *database.LinkRepository
	files          *database.FileRepository
	sessionManager *scs.SessionManager
	storage        storage.Storage
	logger         *slog.Logger
}

func NewLinkHandler(links *database.LinkRepository, files *database.FileRepository, sessionManager *scs.SessionManager, storage storage.Storage, logger *slog.Logger) *LinkHandler {
	return &LinkHandler{
		links:          links,
		files:          files,
		sessionManager: sessionManager,
		storage:        storage,
		logger:         logger.With("handler", "link"),
	}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("id")
	fileID, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		h.logger.Warn("failed to parse file id", "file_id", fileID, "error", err)
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "BAD_REQUEST",
			Message: "Failed to parse file id.",
		})
		return
	}

	file, err := h.files.GetByID(fileID)
	if err != nil {
		h.logger.Warn("file not found", "file_id", fileID, "error", err)
		response.Error(w, http.StatusNotFound, response.APIError{
			Code:    "NOT_FOUND",
			Message: "Failed to fetch file",
		})
		return
	}

	var req LinkRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode link request body", "error", err)
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "BAD_REQUEST",
			Message: "Invalid JSON payload",
		})
		return
	}

	var maxDownloads database.Nullable[int64]
	if req.MaxDownloads != nil {
		maxDownloads.Valid = true
		maxDownloads.V = *req.MaxDownloads
	}

	var expiresAt database.Nullable[time.Time]
	if req.ExpiresAt != nil {
		expiresAt.Valid = true
		expiresAt.V = *req.ExpiresAt

	}

	linkToken, err := token.Generate(token.DefaultLength)
	if err != nil {
		response.InternalError(w, h.logger, "failed to generate token for file", err)
		return
	}

	createdLink, err := h.links.Create(database.Link{
		FileID:       file.ID,
		Token:        linkToken,
		MaxDownloads: maxDownloads,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		response.InternalError(w, h.logger, "failed to save link to database", err)
		return
	}

	response.Success(w, http.StatusCreated, createdLink)
}

func (h *LinkHandler) Download(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	link, err := h.links.GetByToken(token)
	if errors.Is(err, database.ErrNotFound) {
		h.logger.Warn("link not found", "token", token, "error", err)
		response.Error(w, http.StatusNotFound, response.APIError{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("No link with token %s.", token),
		})
		return
	}
	if err != nil {
		response.InternalError(w, h.logger, "failed to get link", err)
		return
	}

	validLink := link.IsValid()
	if !validLink {
		response.Error(w, http.StatusGone, response.APIError{
			Code:    "GONE",
			Message: "This link is no longer available.",
		})
		return
	}

	file, err := h.files.GetByID(link.FileID)
	if errors.Is(err, database.ErrNotFound) {
		h.logger.Warn("file not found", "error", err)
		response.Error(w, http.StatusNotFound, response.APIError{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("No file with associated link with token %s.", token),
		})
		return
	}
	if err != nil {
		response.InternalError(w, h.logger, "failed to get file", err)
		return
	}

	storageFile, err := h.storage.Open(r.Context(), file.StorageKey)
	if err != nil {
		h.logger.Error("failed to open file from storage", "storage_key", file.StorageKey, "error", err)
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An error occurred while trying to fetch uploaded file.",
		})
		return
	}
	defer storageFile.Close()

	escapedFilename := url.PathEscape(file.OriginalFilename)
	dispositionValue := fmt.Sprintf("attachment; filename=\"fallback.bin\"; filename*=UTF-8''%s", escapedFilename)

	w.Header().Set("Content-Disposition", dispositionValue)
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	seeker, ok := storageFile.(io.ReadSeeker)
	if !ok {
		response.InternalError(w, h.logger, "stream does not support seeking", err)
		return
	}

	if err := h.links.IncrementDownloadCount(link.ID); err != nil {
		h.logger.Warn("failed to increment download count", "link_id", link.ID, "error", err)
	}

	http.ServeContent(w, r, "", time.Now(), seeker)
}
