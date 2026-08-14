package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"sahil-api/internal/model"
)

const openverseSearchURL = "https://api.openverse.org/v1/images/"

// OpenverseImage searches the Openverse image API (anonymous access, no key
// required) for the top image matching query.
func OpenverseImage(query string) (*model.OpenverseImage, error) {
	u := fmt.Sprintf("%s?q=%s&page_size=1&mature=false", openverseSearchURL, url.QueryEscape(query))
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Set("User-Agent", "sahil-api/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openverse search: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openverse search: %s", resp.Status)
	}

	var res struct {
		Results []struct {
			Title             string `json:"title"`
			URL               string `json:"url"`
			DirectURL         string `json:"direct_url"`
			Thumbnail         string `json:"thumbnail"`
			ForeignLandingURL string `json:"foreign_landing_url"`
			Creator           string `json:"creator"`
			License           string `json:"license"`
			LicenseURL        string `json:"license_url"`
			Provider          string `json:"provider"`
			Width             int    `json:"width"`
			Height            int    `json:"height"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("openverse search decode: %w", err)
	}

	if len(res.Results) == 0 {
		return nil, nil
	}

	r := res.Results[0]
	return &model.OpenverseImage{
		Title:             r.Title,
		URL:               r.URL,
		DirectURL:         r.DirectURL,
		Thumbnail:         r.Thumbnail,
		ForeignLandingURL: r.ForeignLandingURL,
		Creator:           r.Creator,
		License:           r.License,
		LicenseURL:        r.LicenseURL,
		Provider:          r.Provider,
		Width:             r.Width,
		Height:            r.Height,
	}, nil
}
