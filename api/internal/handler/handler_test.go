package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sahil-api/internal/cache"
	"sahil-api/internal/config"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	cfg := config.FromEnv()
	cfg.LastfmAPIKey = ""
	cfg.LastfmUser = ""
	cfg.AnilistUser = ""
	cfg.GoodreadsUserID = ""
	cfg.LetterboxdUser = ""
	return New(cache.New(time.Minute), cfg)
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	newTestHandler(t).Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", rec.Code)
	}
	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Fatalf("body = %v", body)
	}
}

func TestNowEmptyConfig(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/now", nil)
	rec := httptest.NewRecorder()
	h.Now(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", rec.Code)
	}
	// with no credentials all services return nil -> empty response
	var body map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&body)
	if len(body) != 0 {
		t.Fatalf("expected empty now response, got %v", body)
	}
}

func TestNowCached(t *testing.T) {
	h := newTestHandler(t)
	h.cache.Set("now_response", map[string]interface{}{"cached": true})

	req := httptest.NewRequest(http.MethodGet, "/now", nil)
	rec := httptest.NewRecorder()
	h.Now(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}
	var body map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&body)
	if body["cached"] != true {
		t.Fatalf("expected cached response, got %v", body)
	}
}

func TestKooImageMissingQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/koo", nil)
	rec := httptest.NewRecorder()

	newTestHandler(t).KooImage(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestServiceHandlers(t *testing.T) {
	h := newTestHandler(t)
	cases := []struct {
		name string
		fn   func(http.ResponseWriter, *http.Request)
	}{
		{"LastFM", h.LastFM},
		{"AniList", h.AniList},
		{"Letterboxd", h.Letterboxd},
		{"Goodreads", h.Goodreads},
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c.fn(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: code = %d, want 200 (empty config -> nil data)", c.name, rec.Code)
		}
	}
}

func TestRespondJSONAndError(t *testing.T) {
	rec := httptest.NewRecorder()
	respondJSON(rec, http.StatusTeapot, map[string]string{"a": "b"})
	if rec.Code != http.StatusTeapot {
		t.Fatalf("code = %d", rec.Code)
	}

	rec2 := httptest.NewRecorder()
	respondError(rec2, http.StatusBadRequest, "oops")
	if rec2.Code != http.StatusBadRequest {
		t.Fatalf("code = %d", rec2.Code)
	}
}
