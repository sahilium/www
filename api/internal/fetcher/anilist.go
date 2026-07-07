package fetcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sahil-api/internal/model"
)

const anilistEndpoint = "https://graphql.anilist.co"

func LastAnime(username string) (*model.Anime, error) {
	if username == "" {
		return nil, nil
	}

	query := `query ($username: String) {
		Page(perPage: 1) {
			mediaList(userName: $username, type: ANIME, sort: UPDATED_TIME_DESC) {
				status
				progress
				media {
					title { romaji english }
					coverImage { medium }
					siteUrl
					episodes
				}
				updatedAt
			}
		}
	}`

	body := map[string]interface{}{
		"query": query,
		"variables": map[string]string{"username": username},
	}
	b, _ := json.Marshal(body)

	resp, err := http.Post(anilistEndpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("anilist request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Page struct {
				MediaList []struct {
					Status    string `json:"status"`
					Progress  int    `json:"progress"`
					UpdatedAt int    `json:"updatedAt"`
					Media     struct {
						Title struct {
							Romaji  string `json:"romaji"`
							English string `json:"english"`
						} `json:"title"`
						CoverImage struct {
							Medium string `json:"medium"`
						} `json:"coverImage"`
						SiteUrl  string `json:"siteUrl"`
						Episodes int    `json:"episodes"`
					} `json:"media"`
				} `json:"mediaList"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("anilist decode: %w", err)
	}

	if len(result.Data.Page.MediaList) == 0 {
		return nil, nil
	}

	entry := result.Data.Page.MediaList[0]

	status := "watching"
	if entry.Status == "COMPLETED" {
		status = "completed"
	}

	title := entry.Media.Title.English
	if title == "" {
		title = entry.Media.Title.Romaji
	}

	return &model.Anime{
		Title:         title,
		Image:         entry.Media.CoverImage.Medium,
		Url:           entry.Media.SiteUrl,
		Status:        status,
		Episode:       entry.Progress,
		TotalEpisodes: entry.Media.Episodes,
		UpdatedAt:     time.Unix(int64(entry.UpdatedAt), 0).Format(time.RFC3339),
	}, nil
}
