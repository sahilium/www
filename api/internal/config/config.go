package config

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port           string
	APIBaseURL     string
	CacheTTL       time.Duration
	RequestTimeout time.Duration

	LastfmAPIKey    string
	LastfmUser      string
	AnilistUser     string
	LetterboxdUser  string
	GoodreadsUserID string
}

func Load() *Config {
	loadEnvFile(".env")
	return FromEnv()
}

func FromEnv() *Config {
	return &Config{
		Port:           getenv("PORT", "8080"),
		APIBaseURL:     getenv("API_BASE_URL", "http://localhost:8080"),
		CacheTTL:       5 * time.Minute,
		RequestTimeout: 10 * time.Second,
		LastfmAPIKey:    os.Getenv("LASTFM_API_KEY"),
		LastfmUser:      os.Getenv("LASTFM_USERNAME"),
		AnilistUser:     os.Getenv("ANILIST_USERNAME"),
		LetterboxdUser:  os.Getenv("LETTERBOXD_USERNAME"),
		GoodreadsUserID: os.Getenv("GOODREADS_USER_ID"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok || k == "" {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}
