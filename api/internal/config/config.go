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

	CMSAPIToken         string
	CloudflareAccountID string
	CloudflareAPIToken  string
	D1DatabaseID        string
	R2Endpoint          string
	R2AccessKeyID       string
	R2SecretAccessKey   string
	R2Bucket            string
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

		CMSAPIToken:         os.Getenv("CMS_API_TOKEN"),
		CloudflareAccountID: os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		CloudflareAPIToken:  os.Getenv("CLOUDFLARE_API_TOKEN"),
		D1DatabaseID:        os.Getenv("D1_DATABASE_ID"),
		R2Endpoint:          os.Getenv("R2_ENDPOINT"),
		R2AccessKeyID:       os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey:   os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2Bucket:            os.Getenv("R2_BUCKET"),
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
