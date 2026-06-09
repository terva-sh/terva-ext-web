package fetch

import (
	"testing"
	"time"
)

func TestCacheHitMissAndTTL(t *testing.T) {
	c := newCache(50*time.Millisecond, 4)
	c.put(page{URL: "a", Markdown: "A"})

	if p, ok := c.get("a"); !ok || p.Markdown != "A" {
		t.Fatalf("expected hit for a, got ok=%v p=%+v", ok, p)
	}
	if _, ok := c.get("missing"); ok {
		t.Fatal("expected miss for unknown key")
	}
	time.Sleep(70 * time.Millisecond)
	if _, ok := c.get("a"); ok {
		t.Fatal("expected expiry after TTL")
	}
}

func TestCacheLRUEviction(t *testing.T) {
	c := newCache(0, 2) // no TTL, cap 2
	c.put(page{URL: "a"})
	c.put(page{URL: "b"})
	// Touch a so b becomes least-recently-accessed.
	if _, ok := c.get("a"); !ok {
		t.Fatal("a should be present")
	}
	c.put(page{URL: "c"}) // over cap -> evict LRU (b)
	if _, ok := c.get("b"); ok {
		t.Error("b should have been evicted as least-recently-used")
	}
	if _, ok := c.get("a"); !ok {
		t.Error("a should survive (recently accessed)")
	}
	if _, ok := c.get("c"); !ok {
		t.Error("c should be present (just inserted)")
	}
}

func TestCacheDisabled(t *testing.T) {
	c := newCache(time.Minute, 0) // max 0 disables
	c.put(page{URL: "a", Markdown: "A"})
	if _, ok := c.get("a"); ok {
		t.Fatal("disabled cache should never hit")
	}
}
