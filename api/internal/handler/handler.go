package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"sahil-api/internal/cache"
	"sahil-api/internal/config"
	"sahil-api/internal/fetcher"
	"sahil-api/internal/model"
)

type Handler struct {
	cache *cache.Cache
	cfg   *config.Config
}

func New(c *cache.Cache, cfg *config.Config) *Handler {
	return &Handler{cache: c, cfg: cfg}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Now(w http.ResponseWriter, r *http.Request) {
	const cacheKey = "now_response"
	if cached, ok := h.cache.Get(cacheKey); ok {
		respondJSON(w, http.StatusOK, cached)
		return
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	result := model.NowResponse{}

	fns := []struct {
		name string
		fn   func() (interface{}, error)
	}{
		{"song", func() (interface{}, error) { return fetcher.LastSong(h.cfg.LastfmAPIKey, h.cfg.LastfmUser) }},
		{"book", func() (interface{}, error) { return fetcher.LastBook(h.cfg.HardcoverToken) }},
		{"anime", func() (interface{}, error) { return fetcher.LastAnime(h.cfg.AnilistUser) }},
		{"movie", func() (interface{}, error) { return fetcher.LastMovie(h.cfg.LetterboxdUser) }},
	}

	for _, f := range fns {
		wg.Add(1)
		go func(name string, fn func() (interface{}, error)) {
			defer wg.Done()
			val, err := fn()
			if err != nil || val == nil {
				return
			}
			mu.Lock()
			switch v := val.(type) {
			case *model.Song:
				result.LastSong = v
			case *model.Book:
				result.LastBook = v
			case *model.Anime:
				result.LastAnime = v
			case *model.Movie:
				result.LastMovie = v
			}
			mu.Unlock()
		}(f.name, f.fn)
	}
	wg.Wait()

	h.cache.Set(cacheKey, &result)
	respondJSON(w, http.StatusOK, &result)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{Error: message})
}
