// Package config loads the extension's configuration from its data_dir
// config.json (persisted beside the binary) with environment-variable
// overrides taking precedence — secrets in particular are best passed via env.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config is the effective settings for the web extension.
type Config struct {
	// SearchBackend selects the web_search provider: "tavily" (default) or
	// "searxng".
	SearchBackend string `json:"search_backend"`
	// TavilyAPIKey authenticates the Tavily backend.
	TavilyAPIKey string `json:"tavily_api_key"`
	// SearxngURL is the base URL of a self-hosted SearXNG instance (JSON
	// format must be enabled in its settings.yml).
	SearxngURL string `json:"searxng_url"`

	// FetchMaxBytes caps a fetched response body. Default 2 MiB.
	FetchMaxBytes int64 `json:"fetch_max_bytes"`
	// FetchImageMaxBytes caps the encoded size of an image returned by
	// fetch_image for multimodal injection. Default 5 MiB (≈ provider limits).
	// Images larger than this (after any requested resize) are rejected with a
	// hint to resubmit with a smaller max_dimension. The raw download is allowed
	// to exceed this so an oversized original can be decoded and resized down.
	FetchImageMaxBytes int64 `json:"fetch_image_max_bytes"`
	// FetchTimeoutSec is the overall per-fetch timeout. Default 25s (well
	// under zot's 60s tool budget).
	FetchTimeoutSec int `json:"fetch_timeout_sec"`

	// FetchInlineImages keeps image URLs inline in web_fetch output. Default
	// false: images are replaced with `[image:N]` placeholders and the URLs
	// are retrieved separately via the web_images tool.
	FetchInlineImages bool `json:"fetch_inline_images"`
	// FetchCacheTTLSec is how long a fetched+rendered page stays cached so a
	// follow-up web_images call needs no network. Default 600s. 0 disables.
	FetchCacheTTLSec int `json:"fetch_cache_ttl_sec"`
	// FetchCacheMaxEntries bounds the in-memory page cache (LRU). Default 32.
	FetchCacheMaxEntries int `json:"fetch_cache_max_entries"`

	// AllowLocalHosts is the SSRF escape hatch: targets that resolve to
	// private/reserved addresses are refused UNLESS they match an entry here.
	// Each entry is a hostname (matched against the request host), an IP, or
	// a CIDR (matched against the resolved IP).
	AllowLocalHosts []string `json:"allow_local_hosts"`
}

// Load reads dataDir/config.json (if present), then applies env overrides.
func Load(dataDir string) Config {
	c := Config{
		SearchBackend:        "tavily",
		FetchMaxBytes:        2 << 20, // 2 MiB
		FetchImageMaxBytes:   5 << 20, // 5 MiB
		FetchTimeoutSec:      25,
		FetchCacheTTLSec:     600,
		FetchCacheMaxEntries: 32,
	}
	if dataDir != "" {
		if b, err := os.ReadFile(filepath.Join(dataDir, "config.json")); err == nil {
			_ = json.Unmarshal(b, &c)
		}
	}

	if v := os.Getenv("ZOT_WEB_SEARCH_BACKEND"); v != "" {
		c.SearchBackend = v
	}
	if v := os.Getenv("TAVILY_API_KEY"); v != "" {
		c.TavilyAPIKey = v
	}
	if v := os.Getenv("ZOT_WEB_SEARXNG_URL"); v != "" {
		c.SearxngURL = v
	}
	if v := os.Getenv("ZOT_WEB_ALLOW_LOCAL_HOSTS"); v != "" {
		for _, h := range strings.Split(v, ",") {
			if h = strings.TrimSpace(h); h != "" {
				c.AllowLocalHosts = append(c.AllowLocalHosts, h)
			}
		}
	}
	if v := os.Getenv("ZOT_WEB_FETCH_MAX_BYTES"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			c.FetchMaxBytes = n
		}
	}
	if v := os.Getenv("ZOT_WEB_FETCH_IMAGE_MAX_BYTES"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			c.FetchImageMaxBytes = n
		}
	}
	if v := os.Getenv("ZOT_WEB_FETCH_TIMEOUT_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.FetchTimeoutSec = n
		}
	}
	if v := os.Getenv("ZOT_WEB_FETCH_INLINE_IMAGES"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			c.FetchInlineImages = b
		}
	}
	if v := os.Getenv("ZOT_WEB_FETCH_CACHE_TTL_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			c.FetchCacheTTLSec = n
		}
	}
	if v := os.Getenv("ZOT_WEB_FETCH_CACHE_MAX_ENTRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			c.FetchCacheMaxEntries = n
		}
	}

	c.SearchBackend = strings.ToLower(strings.TrimSpace(c.SearchBackend))
	if c.SearchBackend == "" {
		c.SearchBackend = "tavily"
	}
	if c.FetchMaxBytes <= 0 {
		c.FetchMaxBytes = 2 << 20
	}
	if c.FetchImageMaxBytes <= 0 {
		c.FetchImageMaxBytes = 5 << 20
	}
	if c.FetchTimeoutSec <= 0 {
		c.FetchTimeoutSec = 25
	}
	if c.FetchCacheMaxEntries < 0 {
		c.FetchCacheMaxEntries = 0
	}
	return c
}
