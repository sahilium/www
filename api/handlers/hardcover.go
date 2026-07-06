package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sahil-api/cache"

	"sahil-api/models"
)

const hardcoverEndpoint = "https://api.hardcover.app/v1/graphql"

func FetchLastBook() (*models.Book, error) {
	apiKey := os.Getenv("HARDCOVER_API_KEY")
	userID := os.Getenv("HARDCOVER_USER_ID")
	if apiKey == "" || userID == "" {
		return nil, nil
	}

	query := fmt.Sprintf(`{
		"query": "query { readingLogs(where: {user_id: {_eq: %s}, status: {_eq: \\\"read\\\"}}, order_by: {updated_at: desc}, limit: 1) { book { title image { url(transform: {quality: medium}) } slug } rating } }"
	}`, userID)

	body := bytes.NewReader([]byte(query))
	req, _ := http.NewRequest(http.MethodPost, hardcoverEndpoint, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ReadingLogs []struct {
				Book struct {
					Title string `json:"title"`
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
		return nil, err
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

	return &models.Book{
		Title:  log.Book.Title,
		Author: "",
		Cover:  coverURL,
		Url:    fmt.Sprintf("https://hardcover.app/books/%s", log.Book.Slug),
		Rating: rating,
	}, nil
}

func HardcoverHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cached, ok := c.Get("hardcover"); ok {
			respondJSON(w, http.StatusOK, cached)
			return
		}
		book, err := FetchLastBook()
		if err != nil {
			respondError(w, http.StatusBadGateway, err.Error())
			return
		}
		c.Set("hardcover", book)
		respondJSON(w, http.StatusOK, book)
	}
}
