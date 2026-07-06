package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"sahil-api/cache"

	"sahil-api/models"
)

func NowHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cached, ok := c.Get("now_response"); ok {
			respondJSON(w, http.StatusOK, cached)
			return
		}

		var mu sync.Mutex
		var wg sync.WaitGroup
		result := models.NowResponse{}

		fetchers := []struct {
			key string
			fn  func() (interface{}, error)
		}{
			{"song", func() (interface{}, error) { return FetchLastSong() }},
			{"book", func() (interface{}, error) { return FetchLastBook() }},
			{"anime", func() (interface{}, error) { return FetchLastAnime() }},
			{"movie", func() (interface{}, error) { return FetchLastMovie() }},
		}

		for _, f := range fetchers {
			wg.Add(1)
			go func(key string, fn func() (interface{}, error)) {
				defer wg.Done()
				val, err := fn()
				if err != nil || val == nil {
					return
				}
				mu.Lock()
				switch v := val.(type) {
				case *models.Song:
					result.LastSong = v
				case *models.Book:
					result.LastBook = v
				case *models.Anime:
					result.LastAnime = v
				case *models.Movie:
					result.LastMovie = v
				}
				mu.Unlock()
			}(f.key, f.fn)
		}
		wg.Wait()

		c.Set("now_response", &result)
		respondJSON(w, http.StatusOK, &result)
	}
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
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
