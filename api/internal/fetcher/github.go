package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sahil-api/internal/model"
)

func LastCommit(username string) (*model.Commit, error) {
	if username == "" {
		return nil, nil
	}

	repo, sha, err := latestPushRepo(username)
	if err != nil || repo == "" {
		return nil, err
	}

	return fetchCommit(username, repo, sha)
}

func latestPushRepo(username string) (repo, sha string, err error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public?per_page=5", username)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "sahil-api/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("github events request: %w", err)
	}
	defer resp.Body.Close()

	var events []struct {
		Type string `json:"type"`
		Repo struct {
			Name string `json:"name"`
		} `json:"repo"`
		Payload struct {
			Head string `json:"head"`
			Ref  string `json:"ref"`
		} `json:"payload"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return "", "", fmt.Errorf("github events decode: %w", err)
	}

	for _, e := range events {
		if e.Type == "PushEvent" && e.Payload.Head != "" {
			return e.Repo.Name, e.Payload.Head, nil
		}
	}

	return "", "", nil
}

func fetchCommit(username, repo, sha string) (*model.Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s", repo, sha)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "sahil-api/1.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github commit request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Sha    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Date string `json:"date"`
			} `json:"author"`
		} `json:"commit"`
		HtmlUrl string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("github commit decode: %w", err)
	}

	t, _ := time.Parse(time.RFC3339, result.Commit.Author.Date)

	return &model.Commit{
		Message:   result.Commit.Message,
		Repo:      repo,
		Url:       result.HtmlUrl,
		Timestamp: result.Commit.Author.Date,
		TimeAgo:   timeAgo(t),
	}, nil
}
