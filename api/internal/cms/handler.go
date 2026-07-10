package cms

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"sahil-api/internal/config"
)

type Handler struct {
	d1  *D1Client
	cfg *config.Config
}

func NewHandler(d1 *D1Client, cfg *config.Config) *Handler {
	return &Handler{d1: d1, cfg: cfg}
}

type feedRequest struct {
	Slug     string `json:"slug"`
	Content  string `json:"content"`
}

type feedResponse struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Updated string `json:"updated"`
}

func (h *Handler) PostFeed(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	if strings.TrimPrefix(token, "Bearer ") != h.cfg.CMSAPIToken {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req feedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err), http.StatusBadRequest)
		return
	}

	if req.Slug == "" {
		req.Slug = "feed"
	}
	if req.Content == "" {
		http.Error(w, `{"error":"content required"}`, http.StatusBadRequest)
		return
	}

	title := extractTitle(req.Content)

	if err := h.d1.UpsertFeed(req.Slug, title, req.Content, "{}"); err != nil {
		slog.Error("d1 upsert", "slug", req.Slug, "error", err)
		http.Error(w, `{"error":"storage failed"}`, http.StatusInternalServerError)
		return
	}

	slog.Info("feed updated", "slug", req.Slug, "title", title)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "slug": req.Slug})
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		slug = "feed"
	}

	row, err := h.d1.GetFeed(slug)
	if err != nil {
		slog.Error("d1 get feed", "slug", slug, "error", err)
		http.Error(w, `{"error":"query failed"}`, http.StatusInternalServerError)
		return
	}

	if row == nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	content, _ := row["content"].(string)
	title, _ := row["title"].(string)
	updated, _ := row["updated_at"].(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedResponse{
		Slug:    slug,
		Title:   title,
		Content: content,
		Updated: updated,
	})
}

func extractTitle(content string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimPrefix(trimmed, "# ")
		}
	}
	return ""
}
