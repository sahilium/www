package cache

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	c := New(time.Minute)
	c.Set("key", "value")

	v, ok := c.Get("key")
	if !ok {
		t.Fatal("expected hit")
	}
	if v != "value" {
		t.Fatalf("value = %v, want value", v)
	}
}

func TestGetMissing(t *testing.T) {
	c := New(time.Minute)
	if _, ok := c.Get("nope"); ok {
		t.Fatal("expected miss for missing key")
	}
}

func TestGetExpired(t *testing.T) {
	c := New(time.Minute)
	c.Set("key", "value")

	c.mu.Lock()
	c.items["key"].ExpiresAt = time.Now().Add(-time.Second)
	c.mu.Unlock()

	if _, ok := c.Get("key"); ok {
		t.Fatal("expected miss for expired item")
	}
}

func TestEvictLoopRemovesExpired(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set("key", "value")

	// give the evict loop a tick
	deadline := time.Now().Add(50 * time.Millisecond)
	for time.Now().Before(deadline) {
		if _, ok := c.Get("key"); !ok {
			return // evicted
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatal("expected expired item to be evicted by the loop")
}
