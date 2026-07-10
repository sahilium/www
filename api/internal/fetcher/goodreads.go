package fetcher

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"sahil-api/internal/model"
)

func LastBook(userID string) (*model.Book, error) {
	if userID == "" {
		return nil, nil
	}

	rssURL := fmt.Sprintf("https://www.goodreads.com/review/list_rss/%s?print=true&shelf=currently-reading&per_page=1", userID)
	req, _ := http.NewRequest(http.MethodGet, rssURL, nil)
	req.Header.Set("User-Agent", "sahil-api/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("goodreads request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("goodreads: %s", resp.Status)
	}

	var feed struct {
		Channel struct {
			Items []struct {
				Title        string `xml:"title"`
				Link         string `xml:"link"`
				PubDate      string `xml:"pubDate"`
				AuthorName   string `xml:"author_name"`
				LargeImage   string `xml:"book_large_image_url"`
				MediumImage  string `xml:"book_medium_image_url"`
				UserRating   int    `xml:"user_rating"`
				BookPublished string `xml:"book_published"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("goodreads decode: %w", err)
	}

	if len(feed.Channel.Items) == 0 {
		return nil, nil
	}

	item := feed.Channel.Items[0]

	image := item.MediumImage
	if image == "" {
		image = item.LargeImage
	}

	var ago string
	if item.PubDate != "" {
		for _, f := range []string{time.RFC1123Z, time.RFC1123, "Mon, 2 Jan 2006 15:04:05 -0700"} {
			if t, err := time.Parse(f, item.PubDate); err == nil {
				ago = timeAgo(t)
				break
			}
		}
	}

	return &model.Book{
		Title:  item.Title,
		Author: item.AuthorName,
		Cover:  image,
		Url:    item.Link,
		Rating: item.UserRating,
		TimeAgo: ago,
	}, nil
}
