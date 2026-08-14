package handler

import (
	"net/http"
	"strings"
	"time"

	"sahil-api/internal/fetcher"
)

func (h *Handler) LastFM(w http.ResponseWriter, r *http.Request) {
	handleService(w, h.cache, "lastfm", func() (interface{}, error) {
		return fetcher.LastSong(h.cfg.LastfmAPIKey, h.cfg.LastfmUser)
	})
}

func (h *Handler) AniList(w http.ResponseWriter, r *http.Request) {
	handleService(w, h.cache, "anilist", func() (interface{}, error) {
		return fetcher.LastAnime(h.cfg.AnilistUser)
	})
}

func (h *Handler) Letterboxd(w http.ResponseWriter, r *http.Request) {
	handleService(w, h.cache, "letterboxd", func() (interface{}, error) {
		return fetcher.LastMovie(h.cfg.LetterboxdUser)
	})
}

func (h *Handler) Goodreads(w http.ResponseWriter, r *http.Request) {
	handleService(w, h.cache, "goodreads", func() (interface{}, error) {
		return fetcher.LastBook(h.cfg.GoodreadsUserID)
	})
}

func (h *Handler) Moon(w http.ResponseWriter, r *http.Request) {
	handleService(w, h.cache, "moon", func() (interface{}, error) {
		return fetcher.MoonPhase(h.cfg.MoonLat, h.cfg.MoonLng, time.Now())
	})
}

func (h *Handler) KooImage(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		respondError(w, http.StatusBadRequest, "missing query parameter")
		return
	}
	key := "kooimage:" + q
	if cached, ok := h.cache.Get(key); ok {
		respondJSON(w, http.StatusOK, cached)
		return
	}
	data, err := fetcher.OpenverseImage(q)
	if err != nil {
		respondError(w, http.StatusBadGateway, err.Error())
		return
	}
	h.cache.Set(key, data)
	respondJSON(w, http.StatusOK, data)
}

func handleService(w http.ResponseWriter, c cacheInterface, key string, fn func() (interface{}, error)) {
	if cached, ok := c.Get(key); ok {
		respondJSON(w, http.StatusOK, cached)
		return
	}
	data, err := fn()
	if err != nil {
		respondError(w, http.StatusBadGateway, err.Error())
		return
	}
	c.Set(key, data)
	respondJSON(w, http.StatusOK, data)
}

type cacheInterface interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
}
