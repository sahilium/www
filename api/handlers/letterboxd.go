package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sahil-api/cache"

	"sahil-api/models"
)

func FetchLastMovie() (*models.Movie, error) {
	username := os.Getenv("LETTERBOXD_USERNAME")
	if username == "" {
		return nil, nil
	}

	rssURL := fmt.Sprintf("https://letterboxd.com/%s/rss/", username)
	req, _ := http.NewRequest(http.MethodGet, rssURL, nil)
	req.Header.Set("User-Agent", "sahil-api/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feed struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			Items []struct {
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				Description string `xml:"description"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	if len(feed.Channel.Items) == 0 {
		return nil, nil
	}

	item := feed.Channel.Items[0]
	title := strings.Trim(item.Title, "\" ")

	yearRe := regexp.MustCompile(`\((\d{4})\)`)
	yearMatch := yearRe.FindStringSubmatch(item.Title)
	year := ""
	if len(yearMatch) > 1 {
		year = yearMatch[1]
	}

	imgRe := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
	imgMatch := imgRe.FindStringSubmatch(item.Description)
	image := ""
	if len(imgMatch) > 1 {
		image = imgMatch[1]
	}

	return &models.Movie{
		Title: title,
		Year:  year,
		Image: image,
		Url:   item.Link,
	}, nil
}

func LetterboxdHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cached, ok := c.Get("letterboxd"); ok {
			respondJSON(w, http.StatusOK, cached)
			return
		}
		movie, err := FetchLastMovie()
		if err != nil {
			respondError(w, http.StatusBadGateway, err.Error())
			return
		}
		c.Set("letterboxd", movie)
		respondJSON(w, http.StatusOK, movie)
	}
}
