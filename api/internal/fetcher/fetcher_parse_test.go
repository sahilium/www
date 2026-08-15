package fetcher

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func mustParse(t *testing.T, layout string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006-01-02", layout)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

type mockRoundTripper struct {
	srv  *httptest.Server
	orig http.RoundTripper
}

func (m mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	r.URL.Host = m.srv.Listener.Addr().String()
	return m.orig.RoundTrip(r)
}

// routeAllHTTP points the default transport at the test server for the
// duration of the test, so fetchers that call out over http.DefaultClient or
// http.Get resolve to the mock without code changes.
func routeAllHTTP(t *testing.T, srv *httptest.Server) {
	t.Helper()
	orig := http.DefaultTransport
	http.DefaultTransport = mockRoundTripper{srv: srv, orig: orig}
	t.Cleanup(func() { http.DefaultTransport = orig })
}

func TestLastSongParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"recenttracks":{"track":[
				{"name":"Song X","artist":{"#text":"Artist Y"},
				 "album":{"#text":"Album Z"},
				 "image":[{"size":"small","#text":"s.jpg"},{"size":"large","#text":"l.jpg"}],
				 "url":"http://song","date":{"uts":"1700000000"}}
			]}
		}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	song, err := LastSong("key", "user")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if song.Title != "Song X" || song.Artist != "Artist Y" || song.Album != "Album Z" {
		t.Fatalf("song = %+v", song)
	}
	if song.Image != "l.jpg" {
		t.Fatalf("Image = %q, want l.jpg", song.Image)
	}
	if song.PlayedAt == "" {
		t.Fatal("expected PlayedAt")
	}
}

func TestLastSongListeningNow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"recenttracks":{"track":[
			{"name":"Now","artist":{"#text":"A"},"album":{"#text":"B"},
			 "image":[],"url":"http://s","date":{"uts":""}}
		]}}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	song, err := LastSong("key", "user")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if song.TimeAgo != "listening now" {
		t.Fatalf("TimeAgo = %q, want listening now", song.TimeAgo)
	}
}

func TestLastSongEmptyTracks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"recenttracks":{"track":[]}}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	if s, err := LastSong("key", "user"); err != nil || s != nil {
		t.Fatalf("got (%v, %v), want nil", s, err)
	}
}

func TestLastMovieParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<?xml version="1.0"?>
			<rss version="2.0"><channel><item>
				<title>Movie Title, 2020 - ★★★★</title>
				<link>http://movie</link>
				<pubDate>Mon, 2 Jan 2026 15:04:05 +0000</pubDate>
				<description><![CDATA[<img src="http://img.jpg"/>]]></description>
			</item></channel></rss>`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	movie, err := LastMovie("user")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if movie.Title != "Movie Title" {
		t.Fatalf("Title = %q, want Movie Title", movie.Title)
	}
	if movie.Year != "2020" {
		t.Fatalf("Year = %q, want 2020", movie.Year)
	}
	if movie.Image != "http://img.jpg" {
		t.Fatalf("Image = %q", movie.Image)
	}
	if movie.Url != "http://movie" {
		t.Fatalf("Url = %q", movie.Url)
	}
}

func TestLastMovieEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<rss><channel></channel></rss>`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	if m, err := LastMovie("user"); err != nil || m != nil {
		t.Fatalf("got (%v, %v), want nil", m, err)
	}
}

func TestLastBookParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<?xml version="1.0"?>
			<rss version="2.0"><channel><item>
				<title>Book Title</title>
				<link>http://book</link>
				<pubDate>Mon, 2 Jan 2026 15:04:05 +0000</pubDate>
				<author_name>Author</author_name>
				<book_large_image_url>http://large</book_large_image_url>
				<user_rating>4</user_rating>
			</item></channel></rss>`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	book, err := LastBook("123")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if book.Title != "Book Title" || book.Author != "Author" {
		t.Fatalf("book = %+v", book)
	}
	if book.Cover != "http://large" {
		t.Fatalf("Cover = %q, want http://large (fallback to large)", book.Cover)
	}
	if book.Rating != 4 {
		t.Fatalf("Rating = %d, want 4", book.Rating)
	}
}

func TestLastAnimeParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"Page":{"mediaList":[
			{"status":"COMPLETED","progress":12,"updatedAt":1700000000,
			 "media":{"title":{"romaji":"R","english":"English Title"},
			  "coverImage":{"medium":"http://cover"},
			  "siteUrl":"http://anime","episodes":24}}
		]}}}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	anime, err := LastAnime("user")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if anime.Title != "English Title" {
		t.Fatalf("Title = %q, want English Title", anime.Title)
	}
	if anime.Status != "completed" {
		t.Fatalf("Status = %q, want completed", anime.Status)
	}
	if anime.Episode != 12 || anime.TotalEpisodes != 24 {
		t.Fatalf("anime = %+v", anime)
	}
}

func TestLastAnimeFallbackTitle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"Page":{"mediaList":[
			{"status":"CURRENT","progress":1,"updatedAt":0,
			 "media":{"title":{"romaji":"R","english":""},"coverImage":{"medium":""},"siteUrl":"","episodes":0}}
		]}}}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	anime, err := LastAnime("user")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if anime.Title != "R" {
		t.Fatalf("Title = %q, want R (romaji fallback)", anime.Title)
	}
	if anime.Status != "watching" {
		t.Fatalf("Status = %q, want watching", anime.Status)
	}
}

func TestOpenverseImageParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"results":[{
			"title":"A","url":"http://a","direct_url":"http://d","thumbnail":"http://t",
			"foreign_landing_url":"http://f","creator":"C","license":"BY","license_url":"http://l",
			"provider":"P","width":100,"height":200
		}]}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	img, err := OpenverseImage("query")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if img.Title != "A" || img.Provider != "P" || img.Width != 100 || img.Height != 200 {
		t.Fatalf("img = %+v", img)
	}
}

func TestMoonPhaseParse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"date":"2026-01-01","moon_phase":"Full Moon","moon_illumination":100,"moonrise":"6:00 PM","moonset":"6:00 AM","status":"OK"}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	moon, err := MoonPhase("1", "2", mustParse(t, "2026-01-01"))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if moon.PhaseName != "Full Moon" || moon.Emoji != "🌕" {
		t.Fatalf("moon = %+v", moon)
	}
	if moon.Illumination != 100 {
		t.Fatalf("Illumination = %v", moon.Illumination)
	}
}

func TestOpenverseQueryEscaping(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "q=hello+world") {
			t.Errorf("RawQuery = %q, want escaped query", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"results":[]}`))
	}))
	defer srv.Close()
	routeAllHTTP(t, srv)

	if _, err := OpenverseImage("hello world"); err != nil {
		t.Fatalf("err = %v", err)
	}
}
