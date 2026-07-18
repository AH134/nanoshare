package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/token"
	"github.com/alexedwards/scs/v2"
)

type LinkRequest struct {
	MaxDownloads *int       `json:"maxDownloads"`
	ExpiresAt    *time.Time `json:"expiresAt"`
}

type LinkHandler struct {
	links          *database.LinkRepository
	files          *database.FileRepository
	sessionManager *scs.SessionManager
}

func NewLinkHandler(links *database.LinkRepository, files *database.FileRepository, sessionManager *scs.SessionManager) *LinkHandler {
	return &LinkHandler{
		links:          links,
		files:          files,
		sessionManager: sessionManager,
	}
}

func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("id")
	intValue, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, "invalid file id", http.StatusBadRequest)
		return
	}

	file, err := h.files.GetByID(intValue)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	var req LinkRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to process request body", http.StatusBadRequest)
		return
	}

	var maxDownloads database.Nullable[int64]
	if req.MaxDownloads != nil {
		maxDownloads.Valid = true
		maxDownloads.V = int64(*req.MaxDownloads)
	}

	var expiresAt database.Nullable[time.Time]
	if req.ExpiresAt != nil {
		expiresAt.Valid = true
		expiresAt.V = *req.ExpiresAt

	}

	linkToken, err := token.Generate(token.DefaultLength)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	link := &database.Link{
		FileID:       file.ID,
		Token:        linkToken,
		MaxDownloads: maxDownloads,
		ExpiresAt:    expiresAt,
	}

	if err := h.links.Create(link); err != nil {
		http.Error(w, "failed to upload link", http.StatusInternalServerError)
		return
	}

	createdLink, err := h.links.GetByToken(linkToken)
	if err != nil {
		http.Error(w, "failed to fetch link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdLink); err != nil {
		log.Printf("create link: failed to encode response: %v", err)
	}
}
