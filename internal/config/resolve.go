package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Resolve applies explicit host values without manifest defaults. Errors name
// fields, never their values: even malformed legacy JSON may contain secrets.
func Resolve(dataDir, extensionDir string, host map[string]json.RawMessage) (Config, error) {
	c := Config{SearchBackend: "tavily", FetchMaxBytes: DefaultFetchMaxBytes, FetchImageMaxBytes: DefaultFetchImageMaxBytes, FetchTimeoutSec: DefaultFetchTimeoutSec, FetchCacheTTLSec: DefaultFetchCacheTTLSec, FetchCacheMaxEntries: DefaultFetchCacheMaxEntries, FetchCacheMaxBytes: DefaultFetchCacheMaxBytes, AllowLocalHosts: append([]string(nil), DefaultAllowLocalHosts...)}
	source := "legacy"
	if raw, ok := host["configuration_source"]; ok {
		if err := decodeField("configuration_source", raw, &source); err != nil {
			return c, err
		}
	}
	if value, ok := webEnv("CONFIGURATION_SOURCE"); ok {
		source = value
	}
	if source != "legacy" && source != "host" {
		return c, invalid("configuration_source")
	}
	fields := map[string]any{
		"search_backend": &c.SearchBackend, "searxng_url": &c.SearxngURL, "tavily_api_key": &c.TavilyAPIKey,
		"user_agent": &c.UserAgent, "allow_local_hosts": &c.AllowLocalHosts, "fetch_max_bytes": &c.FetchMaxBytes,
		"fetch_image_max_bytes": &c.FetchImageMaxBytes, "fetch_timeout_sec": &c.FetchTimeoutSec, "fetch_inline_images": &c.FetchInlineImages,
		"fetch_cache_ttl_sec": &c.FetchCacheTTLSec, "fetch_cache_max_entries": &c.FetchCacheMaxEntries, "fetch_cache_max_bytes": &c.FetchCacheMaxBytes,
	}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if source == "legacy" {
		for _, dir := range []string{dataDir, extensionDir} {
			if dir == "" {
				continue
			}
			b, err := os.ReadFile(filepath.Join(dir, "config.json"))
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return c, fmt.Errorf("legacy config.json cannot be read")
			}
			var legacy map[string]json.RawMessage
			if err := json.Unmarshal(b, &legacy); err != nil || legacy == nil {
				return c, fmt.Errorf("legacy config.json must be a JSON object")
			}
			for _, key := range keys {
				if raw, ok := legacy[key]; ok {
					if err := decodeField(key, raw, fields[key]); err != nil {
						return c, err
					}
				}
			}
			break
		}
	}
	for _, key := range keys {
		// Host secret migration is a separate batch; legacy/env owns it for now.
		if key == "tavily_api_key" {
			continue
		}
		if raw, ok := host[key]; ok {
			if key == "allow_local_hosts" {
				var list string
				if err := decodeField(key, raw, &list); err != nil {
					return c, err
				}
				raw = json.RawMessage(list)
			}
			if err := decodeField(key, raw, fields[key]); err != nil {
				return c, err
			}
		}
	}
	for _, key := range keys {
		if key == "tavily_api_key" {
			continue
		}
		value, ok := webEnv(strings.ToUpper(key))
		if !ok {
			continue
		}
		if key == "allow_local_hosts" {
			for _, entry := range strings.Split(value, ",") {
				if entry = strings.TrimSpace(entry); entry != "" {
					c.AllowLocalHosts = append(c.AllowLocalHosts, entry)
				}
			}
			continue
		}
		raw := json.RawMessage(value)
		if _, ok := fields[key].(*string); ok {
			raw, _ = json.Marshal(value)
		}
		if err := decodeField(key, raw, fields[key]); err != nil {
			return c, err
		}
	}
	if value, ok := os.LookupEnv("TAVILY_API_KEY"); ok {
		c.TavilyAPIKey = value
	}
	return c, validate(c)
}

func webEnv(suffix string) (string, bool) {
	if value, ok := os.LookupEnv("TERVA_EXT_WEB_" + suffix); ok {
		return value, true
	}
	return os.LookupEnv("ZOT_WEB_" + suffix)
}
func invalid(key string) error { return fmt.Errorf("invalid setting %s", key) }
func decodeField(key string, raw json.RawMessage, dest any) error {
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, dest) != nil {
		return invalid(key)
	}
	return nil
}
func validate(c Config) error {
	if c.SearchBackend != "tavily" && c.SearchBackend != "searxng" {
		return invalid("search_backend")
	}
	if c.SearxngURL != "" {
		u, err := url.Parse(c.SearxngURL)
		if err != nil || u.Hostname() == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil {
			return invalid("searxng_url")
		}
	}
	if strings.ContainsAny(c.UserAgent, "\r\n\x00") {
		return invalid("user_agent")
	}
	for _, limit := range []struct {
		key             string
		value, min, max int64
	}{
		{"fetch_max_bytes", c.FetchMaxBytes, 1, MaxFetchMaxBytes},
		{"fetch_image_max_bytes", c.FetchImageMaxBytes, 1, MaxFetchImageMaxBytes},
		{"fetch_timeout_sec", int64(c.FetchTimeoutSec), 1, MaxFetchTimeoutSec},
		{"fetch_cache_ttl_sec", int64(c.FetchCacheTTLSec), 0, MaxFetchCacheTTLSec},
		{"fetch_cache_max_entries", int64(c.FetchCacheMaxEntries), 0, MaxFetchCacheMaxEntries},
		{"fetch_cache_max_bytes", c.FetchCacheMaxBytes, 0, MaxFetchCacheMaxBytes},
	} {
		if limit.value < limit.min || limit.value > limit.max {
			return invalid(limit.key)
		}
	}
	for _, entry := range c.AllowLocalHosts {
		if net.ParseIP(entry) != nil {
			continue
		}
		if _, _, err := net.ParseCIDR(entry); err == nil {
			continue
		}
		hostname := strings.TrimSuffix(entry, ".")
		if hostname == "" || len(hostname) > 253 {
			return invalid("allow_local_hosts")
		}
		for _, label := range strings.Split(hostname, ".") {
			if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
				return invalid("allow_local_hosts")
			}
			for _, ch := range label {
				if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || ch == '-') {
					return invalid("allow_local_hosts")
				}
			}
		}
	}
	return nil
}
