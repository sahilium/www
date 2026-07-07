package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"sahil-api/internal/model"
)

func LastSong(apiKey, username string) (*model.Song, error) {
	if apiKey == "" || username == "" {
		return nil, nil
	}

	u := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.getRecentTracks&user=%s&api_key=%s&format=json&limit=1",
		url.QueryEscape(username),
		url.QueryEscape(apiKey),
	)

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("lastfm request: %w", err)
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
		return nil, fmt.Errorf("lastfm decode: %w", err)
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

	var ts int64
	var playedAt string
	var ago string
	if track.Date.Uts != "" {
		ts, _ = strconv.ParseInt(track.Date.Uts, 10, 64)
		t := time.Unix(ts, 0)
		playedAt = t.Format(time.RFC3339)
		ago = timeAgo(t)
	} else {
		ago = "listening now"
	}

	return &model.Song{
		Title:    track.Name,
		Artist:   track.Artist.Content,
		Album:    track.Album.Content,
		Image:    imageURL,
		Url:      track.URL,
		PlayedAt: playedAt,
		TimeAgo:  ago,
	}, nil
}
