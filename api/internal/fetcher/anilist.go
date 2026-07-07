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
		current: MediaListCollection(userName: $username, type: ANIME, status: CURRENT, sort: UPDATED_TIME_DESC, limit: 1) {
			lists { entries { media { title { romaji english } coverImage { medium } siteUrl } updatedAt } }
		}
		completed: MediaListCollection(userName: $username, type: ANIME, status: COMPLETED, sort: UPDATED_TIME_DESC, limit: 1) {
			lists { entries { media { title { romaji english } coverImage { medium } siteUrl } updatedAt } }
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
			Current   anilistCollection `json:"current"`
			Completed anilistCollection `json:"completed"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("anilist decode: %w", err)
	}

	entry := pickAnimeEntry(result.Data.Current, result.Data.Completed)
	if entry == nil {
		return nil, nil
	}

	status := "watching"
	if entry == firstAnimeEntry(result.Data.Completed) {
		status = "completed"
	}

	return &model.Anime{
		Title:     pickTitle(entry.Media.Title.Romaji, entry.Media.Title.English),
		Image:     entry.Media.CoverImage.Medium,
		Url:       entry.Media.SiteUrl,
		Status:    status,
		UpdatedAt: time.Unix(int64(entry.UpdatedAt), 0).Format(time.RFC3339),
	}, nil
}

type anilistCollection struct {
	Lists []struct {
		Entries []anilistEntry `json:"entries"`
	} `json:"lists"`
}

type anilistEntry struct {
	Media struct {
		Title struct {
			Romaji  string `json:"romaji"`
			English string `json:"english"`
		} `json:"title"`
		CoverImage struct {
			Medium string `json:"medium"`
		} `json:"coverImage"`
		SiteUrl string `json:"siteUrl"`
	} `json:"media"`
	UpdatedAt int `json:"updatedAt"`
}

func pickAnimeEntry(current, completed anilistCollection) *anilistEntry {
	if e := firstAnimeEntry(current); e != nil {
		return e
	}
	return firstAnimeEntry(completed)
}

func firstAnimeEntry(c anilistCollection) *anilistEntry {
	for _, list := range c.Lists {
		if len(list.Entries) > 0 {
			return &list.Entries[0]
		}
	}
	return nil
}

func pickTitle(romaji, english string) string {
	if english != "" {
		return english
	}
	return romaji
}
