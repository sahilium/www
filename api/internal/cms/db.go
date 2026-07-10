package cms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const cloudflareAPI = "https://api.cloudflare.com/client/v4"

type D1Client struct {
	accountID  string
	databaseID string
	apiToken   string
	httpClient *http.Client
}

func NewD1Client(accountID, databaseID, apiToken string) *D1Client {
	return &D1Client{
		accountID:  accountID,
		databaseID: databaseID,
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type d1Request struct {
	SQL    string        `json:"sql"`
	Params []interface{} `json:"params,omitempty"`
}

type d1Response struct {
	Result  []d1QueryResult  `json:"result"`
	Success bool             `json:"success"`
	Errors  []d1Error        `json:"errors"`
}

type d1QueryResult struct {
	Results []map[string]interface{} `json:"results"`
	Success bool                     `json:"success"`
}

type d1Error struct {
	Message string `json:"message"`
}

func (d *D1Client) Query(sql string, params ...interface{}) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/accounts/%s/d1/database/%s/query", cloudflareAPI, d.accountID, d.databaseID)

	body, _ := json.Marshal(d1Request{SQL: sql, Params: params})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+d.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("d1 request: %w", err)
	}
	defer resp.Body.Close()

	var d1resp d1Response
	if err := json.NewDecoder(resp.Body).Decode(&d1resp); err != nil {
		return nil, fmt.Errorf("d1 decode: %w", err)
	}

	if !d1resp.Success {
		msg := "unknown error"
		if len(d1resp.Errors) > 0 {
			msg = d1resp.Errors[0].Message
		}
		return nil, fmt.Errorf("d1: %s", msg)
	}

	if len(d1resp.Result) == 0 || !d1resp.Result[0].Success {
		return nil, fmt.Errorf("d1 query failed")
	}

	return d1resp.Result[0].Results, nil
}

func (d *D1Client) UpsertFeed(slug, title, content, metadata string) error {
	existing, err := d.Query("SELECT id FROM feed_entries WHERE slug = ?", slug)
	if err != nil {
		return fmt.Errorf("d1 upsert check: %w", err)
	}

	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	if len(existing) > 0 {
		_, err = d.Query(
			"UPDATE feed_entries SET title = ?, content = ?, metadata = ?, updated_at = ? WHERE slug = ?",
			title, content, metadata, now, slug,
		)
	} else {
		_, err = d.Query(
			"INSERT INTO feed_entries (slug, title, content, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			slug, title, content, metadata, now, now,
		)
	}
	return err
}

func (d *D1Client) GetFeed(slug string) (map[string]interface{}, error) {
	rows, err := d.Query(
		"SELECT slug, title, content, metadata, created_at, updated_at FROM feed_entries WHERE slug = ? ORDER BY updated_at DESC LIMIT 1",
		slug,
	)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return rows[0], nil
}
