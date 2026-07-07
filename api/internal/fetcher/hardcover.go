package fetcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sahil-api/internal/model"
)

const hardcoverEndpoint = "https://api.hardcover.app/v1/graphql"

func LastBook(token string) (*model.Book, error) {
	if token == "" {
		return nil, nil
	}

	query := `query {
		user_book_reads(
			where: {finished_at: {_is_null: false}},
			order_by: {finished_at: desc},
			limit: 1
		) {
			finished_at
			progress
			user_book {
				book {
					title
					slug
					image { url }
				}
				rating
			}
		}
	}`

	body := map[string]string{"query": query}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, hardcoverEndpoint, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hardcover request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			UserBookReads []struct {
				FinishedAt string   `json:"finished_at"`
				Progress   *float64 `json:"progress"`
				UserBook   struct {
					Book struct {
						Title string `json:"title"`
						Slug  string `json:"slug"`
						Image *struct {
							URL string `json:"url"`
						} `json:"image"`
					} `json:"book"`
					Rating *int `json:"rating"`
				} `json:"user_book"`
			} `json:"user_book_reads"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("hardcover decode: %w", err)
	}

	if len(result.Data.UserBookReads) == 0 {
		return nil, nil
	}

	read := result.Data.UserBookReads[0]

	var coverURL string
	if read.UserBook.Book.Image != nil {
		coverURL = read.UserBook.Book.Image.URL
	}

	rating := 0
	if read.UserBook.Rating != nil {
		rating = *read.UserBook.Rating
	}

	var ago string
	if read.FinishedAt != "" {
		if t, err := time.Parse("2006-01-02", read.FinishedAt); err == nil {
			ago = timeAgo(t)
		}
	}

	return &model.Book{
		Title:  read.UserBook.Book.Title,
		Author: "",
		Cover:  coverURL,
		Url:    fmt.Sprintf("https://hardcover.app/books/%s", read.UserBook.Book.Slug),
		Rating: rating,
		TimeAgo: ago,
	}, nil
}
