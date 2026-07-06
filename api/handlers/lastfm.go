package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"sahil-api/cache"

	"sahil-api/models"
)

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

func FetchLastSong() (*models.Song, error) {
	apiKey := os.Getenv("LASTFM_API_KEY")
	username := os.Getenv("LASTFM_USERNAME")
	if apiKey == "" || username == "" {
		return nil, nil
	}

	u := fmt.Sprintf(
		"http://ws.audioscrobbler.com/2.0/?method=user.getRecentTracks&user=%s&api_key=%s&format=json&limit=1",
		url.QueryEscape(username),
		url.QueryEscape(apiKey),
	)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		RecentTracks struct {
			Track []struct {
				Name   string `json:"name"`
				Artist struct {
					Content string `json:"#text"`
				} `json:"artist"`
				Album struct {
					Content string `json:"#text"`
				} `json:"album"`
				Image []struct {
					Size string `json:"size"`
					URL  string `json:"#text"`
				} `json:"image"`
				URL  string `json:"url"`
				Date struct {
					Uts string `json:"uts"`
				} `json:"date"`
			} `json:"track"`
		} `json:"recenttracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	if len(raw.RecentTracks.Track) == 0 {
		return nil, nil
	}

	track := raw.RecentTracks.Track[0]
	var imageURL string
	for _, img := range track.Image {
		if img.Size == "large" {
			imageURL = img.URL
			break
		}
	}

	var playedAt string
	if track.Date.Uts != "" {
		playedAt = time.Unix(parseInt64(track.Date.Uts), 0).Format(time.RFC3339)
	}

	return &models.Song{
		Title:    track.Name,
		Artist:   track.Artist.Content,
		Album:    track.Album.Content,
		Image:    imageURL,
		Url:      track.URL,
		PlayedAt: playedAt,
	}, nil
}

func LastFMHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cached, ok := c.Get("lastfm"); ok {
			respondJSON(w, http.StatusOK, cached)
			return
		}
		song, err := FetchLastSong()
		if err != nil {
			respondError(w, http.StatusBadGateway, err.Error())
			return
		}
		c.Set("lastfm", song)
		respondJSON(w, http.StatusOK, song)
	}
}
