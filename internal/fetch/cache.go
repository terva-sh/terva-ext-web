package fetch

import (
	"sync"
	"time"
)

// Image is one image found on a fetched page. ID is the page-local handle the
// model sees as `[image:ID]`; URL is resolved absolute against the page URL.
type Image struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
	Alt string `json:"alt,omitempty"`
}

// page is the fully rendered result of a fetch, cached so a follow-up
// web_images call needs no network. Markdown is untruncated (carrying
// `[image:N]` placeholders unless inline images are configured); per-request
// character limits are applied at response time.
type page struct {
	URL           string // the requested URL (cache key)
	FinalURL      string // URL after redirects
	Title         string
	ContentType   string
	Status        int
	Markdown      string
	Images        []Image
	BodyTruncated bool // raw body hit the byte cap before rendering
}

// cache is a small, concurrency-safe, TTL + LRU page cache. Tool handlers run
// in their own goroutines (see proto.Run), so every access takes the lock.
type cache struct {
	mu      sync.Mutex
	ttl     time.Duration
	max     int
	entries map[string]*entry
}

type entry struct {
	page     page
	stored   time.Time
	accessed time.Time
}

// newCache builds a cache. A non-positive max disables caching entirely (get
// always misses, put is a no-op); a non-positive ttl means entries never
// expire by age (still bounded by max).
func newCache(ttl time.Duration, max int) *cache {
	return &cache{ttl: ttl, max: max, entries: map[string]*entry{}}
}

// get returns the cached page for url and whether it was a live hit. Expired
// entries are dropped and reported as a miss.
func (c *cache) get(url string) (page, bool) {
	if c.max <= 0 {
		return page{}, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[url]
	if !ok {
		return page{}, false
	}
	if c.ttl > 0 && time.Since(e.stored) > c.ttl {
		delete(c.entries, url)
		return page{}, false
	}
	e.accessed = time.Now()
	return e.page, true
}

// put stores p, evicting the least-recently-accessed entry if over capacity.
func (c *cache) put(p page) {
	if c.max <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.entries[p.URL] = &entry{page: p, stored: now, accessed: now}
	for len(c.entries) > c.max {
		var oldestKey string
		var oldest time.Time
		for k, e := range c.entries {
			if oldestKey == "" || e.accessed.Before(oldest) {
				oldestKey, oldest = k, e.accessed
			}
		}
		delete(c.entries, oldestKey)
	}
}
