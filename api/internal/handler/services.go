package handler

import (
	"net/http"

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

func (h *Handler) Hardcover(w http.ResponseWriter, r *http.Request) {
	handleService(w, h.cache, "hardcover", func() (interface{}, error) {
		return fetcher.LastBook(h.cfg.HardcoverToken)
	})
}

func (h *Handler) GitHubCommits(w http.ResponseWriter, r *http.Request) {
	handleService(w, h.cache, "github:commits", func() (interface{}, error) {
		return fetcher.LastCommit(h.cfg.GitHubUser)
	})
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
