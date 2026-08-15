package config

import (
	"os"
	"testing"
	"time"
)

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("API_BASE_URL", "")
	t.Setenv("MOON_LAT", "")
	t.Setenv("MOON_LNG", "")

	cfg := FromEnv()

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.APIBaseURL != "http://localhost:8080" {
		t.Errorf("APIBaseURL = %q, want http://localhost:8080", cfg.APIBaseURL)
	}
	if cfg.MoonLat != "35.6895" {
		t.Errorf("MoonLat = %q, want 35.6895", cfg.MoonLat)
	}
	if cfg.MoonLng != "139.6917" {
		t.Errorf("MoonLng = %q, want 139.6917", cfg.MoonLng)
	}
	if cfg.CacheTTL != 5*time.Minute {
		t.Errorf("CacheTTL = %v, want 5m", cfg.CacheTTL)
	}
	if cfg.RequestTimeout != 10*time.Second {
		t.Errorf("RequestTimeout = %v, want 10s", cfg.RequestTimeout)
	}
}

func TestFromEnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("API_BASE_URL", "https://api.example.com")
	t.Setenv("LASTFM_USERNAME", "bob")
	t.Setenv("LASTFM_API_KEY", "k1")
	t.Setenv("ANILIST_USERNAME", "alice")
	t.Setenv("GOODREADS_USER_ID", "123")
	t.Setenv("LETTERBOXD_USERNAME", "carol")
	t.Setenv("CMS_API_TOKEN", "tok")

	cfg := FromEnv()

	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want 9090", cfg.Port)
	}
	if cfg.APIBaseURL != "https://api.example.com" {
		t.Errorf("APIBaseURL = %q", cfg.APIBaseURL)
	}
	if cfg.LastfmUser != "bob" || cfg.LastfmAPIKey != "k1" {
		t.Errorf("lastfm user/key not read: %+v", cfg)
	}
	if cfg.AnilistUser != "alice" {
		t.Errorf("AnilistUser = %q", cfg.AnilistUser)
	}
	if cfg.GoodreadsUserID != "123" {
		t.Errorf("GoodreadsUserID = %q", cfg.GoodreadsUserID)
	}
	if cfg.LetterboxdUser != "carol" {
		t.Errorf("LetterboxdUser = %q", cfg.LetterboxdUser)
	}
	if cfg.CMSAPIToken != "tok" {
		t.Errorf("CMSAPIToken = %q", cfg.CMSAPIToken)
	}
}

func TestLoadEnvFile(t *testing.T) {
	path := ".env"
	if err := os.WriteFile(path, []byte("# comment\nPORT=9999\nEMPTY=\nBROKEN_LINE\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	t.Setenv("PORT", "")
	loadEnvFile(path)

	if v := os.Getenv("PORT"); v != "9999" {
		t.Errorf("PORT = %q, want 9999 (loaded from .env)", v)
	}
}

func TestGetenvFallback(t *testing.T) {
	t.Setenv("FOO_GETENV", "")
	if got := getenv("FOO_GETENV", "fb"); got != "fb" {
		t.Errorf("getenv = %q, want fb", got)
	}
	t.Setenv("FOO_GETENV", "set")
	if got := getenv("FOO_GETENV", "fb"); got != "set" {
		t.Errorf("getenv = %q, want set", got)
	}
}
