package fetcher

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"sahil-api/internal/model"
)

var letterboxdTimeFormats = []string{
	time.RFC1123Z,
	time.RFC1123,
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
}

func LastMovie(username string) (*model.Movie, error) {
	if username == "" {
		return nil, nil
	}

	rssURL := fmt.Sprintf("https://letterboxd.com/%s/rss/", username)
	req, _ := http.NewRequest(http.MethodGet, rssURL, nil)
	req.Header.Set("User-Agent", "sahil-api/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("letterboxd request: %w", err)
	}
	defer resp.Body.Close()

	var feed struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			Items []struct {
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				PubDate     string `xml:"pubDate"`
				Description string `xml:"description"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("letterboxd decode: %w", err)
	}

	if len(feed.Channel.Items) == 0 {
		return nil, nil
	}

	item := feed.Channel.Items[0]

	ratingRe := regexp.MustCompile(`\s+-\s+[★☆½¼¾]+$`)
	title := ratingRe.ReplaceAllString(strings.Trim(item.Title, "\" "), "")

	yearRe := regexp.MustCompile(`, (\d{4})`)
	yearMatch := yearRe.FindStringSubmatch(item.Title)
	year := ""
	if len(yearMatch) > 1 {
		year = yearMatch[1]
	}

	title = strings.TrimSuffix(title, ", "+year)

	imgRe := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
	imgMatch := imgRe.FindStringSubmatch(item.Description)
	image := ""
	if len(imgMatch) > 1 {
		image = imgMatch[1]
	}

	var ago string
	if item.PubDate != "" {
		for _, f := range letterboxdTimeFormats {
			if t, err := time.Parse(f, item.PubDate); err == nil {
				ago = timeAgo(t)
				break
			}
		}
	}

	return &model.Movie{
		Title:  title,
		Year:   year,
		Image:  image,
		Url:    item.Link,
		TimeAgo: ago,
	}, nil
}
