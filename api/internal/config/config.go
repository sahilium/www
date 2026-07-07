package config

import (
	"os"
	"time"
)

type Config struct {
	Port          string
	APIBaseURL    string
	CacheTTL      time.Duration
	RequestTimeout time.Duration

	LastfmAPIKey  string
	LastfmUser    string
	AnilistUser   string
	LetterboxdUser string
	HardcoverToken string
}

func FromEnv() *Config {
	return &Config{
		Port:           getenv("PORT", "8080"),
		APIBaseURL:     getenv("API_BASE_URL", "http://localhost:8080"),
		CacheTTL:       5 * time.Minute,
		RequestTimeout: 10 * time.Second,
		LastfmAPIKey:   os.Getenv("LASTFM_API_KEY"),
		LastfmUser:     os.Getenv("LASTFM_USERNAME"),
		AnilistUser:    os.Getenv("ANILIST_USERNAME"),
		LetterboxdUser: os.Getenv("LETTERBOXD_USERNAME"),
		HardcoverToken: os.Getenv("HARDCOVER_TOKEN"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
