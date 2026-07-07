package fetcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"sahil-api/internal/model"
)

const hardcoverEndpoint = "https://api.hardcover.app/v1/graphql"

func LastBook(token string) (*model.Book, error) {
	if token == "" {
		return nil, nil
	}

	query := `{
		"query": "query { readingLogs(where: {status: {_eq: \\\"read\\\"}}, order_by: {updated_at: desc}, limit: 1) { book { title image { url(transform: {quality: medium}) } slug } rating } }"
	}`

	body := bytes.NewReader([]byte(query))
	req, _ := http.NewRequest(http.MethodPost, hardcoverEndpoint, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hardcover request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ReadingLogs []struct {
				Book struct {
					Title string  `json:"title"`
					Image *struct {
						URL string `json:"url"`
					} `json:"image"`
					Slug string `json:"slug"`
				} `json:"book"`
				Rating *int `json:"rating"`
			} `json:"readingLogs"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("hardcover decode: %w", err)
	}

	if len(result.Data.ReadingLogs) == 0 {
		return nil, nil
	}

	log := result.Data.ReadingLogs[0]
	var coverURL string
	if log.Book.Image != nil {
		coverURL = log.Book.Image.URL
	}

	rating := 0
	if log.Rating != nil {
		rating = *log.Rating
	}

	return &model.Book{
		Title:  log.Book.Title,
		Author: "",
		Cover:  coverURL,
		Url:    fmt.Sprintf("https://hardcover.app/books/%s", log.Book.Slug),
		Rating: rating,
	}, nil
}
