// Package config resolves Terva host settings and environment
// overrides. Resolve validates settings before the runtime accepts them.
package config

const (
	DefaultFetchMaxBytes        int64 = 2 << 20 // 2 MiB
	DefaultFetchImageMaxBytes   int64 = 5 << 20 // 5 MiB
	DefaultFetchTimeoutSec            = 25
	DefaultFetchCacheTTLSec           = 600
	DefaultFetchCacheMaxEntries       = 32
	DefaultFetchCacheMaxBytes   int64 = 64 << 20 // 64 MiB

	MaxFetchMaxBytes        int64 = 32 << 20 // 32 MiB
	MaxFetchImageMaxBytes   int64 = 20 << 20 // 20 MiB
	MaxFetchTimeoutSec            = 60
	MaxFetchCacheTTLSec           = 3600 // 1 hour
	MaxFetchCacheMaxEntries       = 128
	MaxFetchCacheMaxBytes   int64 = 256 << 20 // 256 MiB
)

// DefaultAllowLocalHosts is the out-of-the-box SSRF allowlist: loopback only,
// by name and by literal address (the hostname entry already covers whatever
// "localhost" resolves to; the IPs cover URLs that dial 127.0.0.1/[::1]
// directly). Everything else private/reserved stays blocked until the user
// opts in.
var DefaultAllowLocalHosts = []string{"localhost", "127.0.0.1", "::1"}

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

	// FetchMaxBytes caps a fetched response body. Default 2 MiB; max 32 MiB.
	FetchMaxBytes int64 `json:"fetch_max_bytes"`
	// FetchImageMaxBytes caps the encoded size of an image returned by
	// web_fetch_image for multimodal injection. Default 5 MiB; max 20 MiB.
	// Images larger than this (after any requested resize) are rejected with a
	// hint to resubmit with a smaller max_dimension. The raw download is allowed
	// to exceed this so an oversized original can be decoded and resized down.
	FetchImageMaxBytes int64 `json:"fetch_image_max_bytes"`
	// FetchTimeoutSec is the overall per-fetch timeout. Default 25s; max 60s.
	FetchTimeoutSec int `json:"fetch_timeout_sec"`

	// FetchInlineImages keeps image URLs inline in web_fetch output. Default
	// false: images are replaced with `[image:N]` placeholders and the URLs
	// are retrieved separately via the web_images tool.
	FetchInlineImages bool `json:"fetch_inline_images"`
	// FetchCacheTTLSec is how long a fetched+rendered page stays cached so a
	// follow-up web_images call needs no network. Default 600s. 0 means no expiry.
	FetchCacheTTLSec int `json:"fetch_cache_ttl_sec"`
	// FetchCacheMaxEntries bounds the in-memory page cache (LRU). Default 32; max 128.
	FetchCacheMaxEntries int `json:"fetch_cache_max_entries"`
	// FetchCacheMaxBytes bounds the page cache by total retained bytes (rendered
	// Markdown + compressed raw body + harvested links/images), evicting LRU
	// entries until under budget. A single page larger than the budget is retained
	// as the only entry. Default 64 MiB;
	// max 256 MiB. 0 disables the byte bound (entry count still applies).
	FetchCacheMaxBytes int64 `json:"fetch_cache_max_bytes"`

	// UserAgent overrides the User-Agent sent on every fetch. Empty means the
	// default "terva-ext-web/<version>". The special value "browser" expands to a
	// common desktop-browser UA, for sites that block non-browser clients.
	// A per-call user_agent tool parameter takes precedence over this.
	UserAgent string `json:"user_agent"`

	// AllowLocalHosts is the SSRF escape hatch: targets that resolve to
	// private/reserved addresses are refused UNLESS they match an entry here.
	// Each entry is a hostname (matched against the request host), an IP, or
	// a CIDR (matched against the resolved IP). Defaults to loopback
	// (DefaultAllowLocalHosts); the host field replaces the default. Supply
	// the full list to extend it, or [] to lock loopback back down. The env
	// override appends instead.
	AllowLocalHosts []string `json:"allow_local_hosts"`
}
