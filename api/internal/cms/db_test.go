package cms

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// newTestD1 returns a D1Client whose requests are routed to the local test
// server regardless of the hardcoded Cloudflare URL.
func newTestD1(srv *httptest.Server) *D1Client {
	return &D1Client{
		accountID:  "acc",
		databaseID: "db",
		apiToken:   "tok",
		httpClient: &http.Client{
			Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				r.URL.Scheme = "http"
				r.URL.Host = srv.Listener.Addr().String()
				return http.DefaultTransport.RoundTrip(r)
			}),
		},
	}
}

func TestQuerySuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tok" {
			t.Errorf("Authorization = %q, want Bearer tok", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		var req d1Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if req.SQL != "SELECT 1" {
			t.Errorf("SQL = %q", req.SQL)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{{"slug": "x"}}}},
			"success": true,
		})
	}))
	defer srv.Close()

	d := newTestD1(srv)

	rows, err := d.Query("SELECT 1")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if len(rows) != 1 || rows[0]["slug"] != "x" {
		t.Fatalf("rows = %v", rows)
	}
}

func TestQueryFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"errors":  []map[string]interface{}{{"message": "auth failed"}},
		})
	}))
	defer srv.Close()

	d := newTestD1(srv)
	if _, err := d.Query("SELECT 1"); err == nil {
		t.Fatal("expected error for unsuccessful response")
	}
}

func TestQueryResultUnsuccessful(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": false, "results": []map[string]interface{}{}}},
			"success": true,
		})
	}))
	defer srv.Close()

	d := newTestD1(srv)
	if _, err := d.Query("SELECT 1"); err == nil {
		t.Fatal("expected error for unsuccessful query result")
	}
}

func TestUpsertFeedInserts(t *testing.T) {
	var gotSQL string
	var gotParams []interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req d1Request
		json.NewDecoder(r.Body).Decode(&req)
		gotSQL = req.SQL
		gotParams = req.Params
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{}}},
			"success": true,
		})
	}))
	defer srv.Close()

	d := newTestD1(srv)
	if err := d.UpsertFeed("my-slug", "Title", "body", "{}"); err != nil {
		t.Fatalf("err = %v", err)
	}
	if gotSQL == "" {
		t.Fatal("expected a query to be issued")
	}
	if len(gotParams) == 0 {
		t.Fatal("expected query params")
	}
}

func TestGetFeedFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{{"slug": "s", "title": "T", "content": "C"}}}},
			"success": true,
		})
	}))
	defer srv.Close()

	d := newTestD1(srv)
	row, err := d.GetFeed("s")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if row == nil || row["slug"] != "s" {
		t.Fatalf("row = %v", row)
	}
}

func TestGetFeedEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  []map[string]interface{}{{"success": true, "results": []map[string]interface{}{}}},
			"success": true,
		})
	}))
	defer srv.Close()

	d := newTestD1(srv)
	row, err := d.GetFeed("s")
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if row != nil {
		t.Fatalf("row = %v, want nil", row)
	}
}

func TestQueryHTTPError(t *testing.T) {
	// server that errors on the wire
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("connection reset")
	}))
	defer srv.Close()

	d := newTestD1(srv)
	if _, err := d.Query("SELECT 1"); err == nil {
		t.Fatal("expected error on HTTP failure")
	}
}
