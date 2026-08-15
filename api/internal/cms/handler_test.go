package cms

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sahil-api/internal/config"
)

func newTestHandler(t *testing.T) (*Handler, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req d1Request
		json.NewDecoder(r.Body).Decode(&req)
		if strings.HasPrefix(req.SQL, "SELECT") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{}}},
				"success": true,
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{{"id": 1}}}},
			"success": true,
		})
	}))
	t.Cleanup(srv.Close)

	cfg := config.FromEnv()
	cfg.CMSAPIToken = "secret"
	d1 := newTestD1(srv)
	return NewHandler(d1, cfg), srv
}

func TestPostFeedUnauthorized(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	h.PostFeed(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("code = %d, want 401", rec.Code)
	}
}

func TestPostFeedWrongToken(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	rec := httptest.NewRecorder()
	h.PostFeed(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("code = %d, want 401", rec.Code)
	}
}

func TestPostFeedBadBody(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not json"))
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	h.PostFeed(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestPostFeedMissingContent(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"slug":"x","content":""}`))
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	h.PostFeed(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestPostFeedSuccess(t *testing.T) {
	h, _ := newTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"slug":"my-post","content":"# Hello\nbody"}`))
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	h.PostFeed(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200 (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestGetFeedNotFound(t *testing.T) {
	h, _ := newTestHandler(t)
	// point at a fresh server that returns no rows
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{}}},
			"success": true,
		})
	}))
	defer srv2.Close()
	h.d1 = newTestD1(srv2)

	req := httptest.NewRequest(http.MethodGet, "/?slug=x", nil)
	rec := httptest.NewRecorder()
	h.GetFeed(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404", rec.Code)
	}
}

func TestGetFeedFoundHandler(t *testing.T) {
	h, _ := newTestHandler(t)
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{{"slug": "x", "title": "T", "content": "C", "updated_at": "2026"}}}},
			"success": true,
		})
	}))
	defer srv2.Close()
	h.d1 = newTestD1(srv2)

	req := httptest.NewRequest(http.MethodGet, "/?slug=x", nil)
	rec := httptest.NewRecorder()
	h.GetFeed(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"title":"T"`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestExtractTitle(t *testing.T) {
	if got := extractTitle("# Hello World\nsome body"); got != "Hello World" {
		t.Fatalf("extractTitle = %q, want Hello World", got)
	}
	if got := extractTitle("no heading here"); got != "" {
		t.Fatalf("extractTitle = %q, want empty", got)
	}
	if got := extractTitle("## Subheading\n# Real"); got != "Real" {
		t.Fatalf("extractTitle = %q, want Real", got)
	}
}
