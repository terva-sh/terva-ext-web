package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestResolveDefaults(t *testing.T) {
	isolateResolveEnv(t)
	got, err := Resolve(nil)
	want := Config{
		SearchBackend: "tavily", FetchMaxBytes: 2 << 20, FetchImageMaxBytes: 5 << 20,
		FetchTimeoutSec: 25, FetchCacheTTLSec: 600, FetchCacheMaxEntries: 32,
		FetchCacheMaxBytes: 64 << 20, AllowLocalHosts: []string{"localhost", "127.0.0.1", "::1"},
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected application defaults: %v", err)
	}
	got.AllowLocalHosts[0] = "changed.invalid"
	next, err := Resolve(nil)
	if err != nil || next.AllowLocalHosts[0] != "localhost" {
		t.Fatal("one configuration changed the default allowlist")
	}
}

func TestResolveAllSettingsFromHostAndEnvironment(t *testing.T) {
	isolateResolveEnv(t)
	values := map[string]any{
		"search_backend": " SEARXNG ", "searxng_url": "http://search.invalid",
		"tavily_api_key": "synthetic-config-fixture", "user_agent": "test-client",
		"fetch_max_bytes": 1048576, "fetch_image_max_bytes": 2097152,
		"fetch_timeout_sec": 10, "fetch_inline_images": true,
		"fetch_cache_ttl_sec": 300, "fetch_cache_max_entries": 16,
		"fetch_cache_max_bytes": 33554432,
	}
	want := Config{
		SearchBackend: "searxng", SearxngURL: "http://search.invalid",
		TavilyAPIKey: "synthetic-config-fixture", UserAgent: "test-client",
		FetchMaxBytes: 1048576, FetchImageMaxBytes: 2097152, FetchTimeoutSec: 10,
		FetchInlineImages: true, FetchCacheTTLSec: 300, FetchCacheMaxEntries: 16,
		FetchCacheMaxBytes: 33554432, AllowLocalHosts: []string{"localhost", "127.0.0.1", "::1"},
	}
	host := make(map[string]json.RawMessage)
	for key, value := range values {
		host[key], _ = json.Marshal(value)
	}
	got, err := Resolve(host)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("host settings not applied: %v", err)
	}
	for key, value := range values {
		name := "TERVA_EXT_WEB_" + strings.ToUpper(key)
		if key == "tavily_api_key" {
			name = "TAVILY_API_KEY"
		}
		t.Setenv(name, fmt.Sprint(value))
	}
	got, err = Resolve(nil)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("environment settings not applied: %v", err)
	}
}

// The production resolver rejects invalid limits. The retired loader silently
// clamped them, so its tests could pass while real configuration failed.
func TestResolveNumericLimits(t *testing.T) {
	isolateResolveEnv(t)
	for _, tc := range []struct {
		key      string
		min, max int64
	}{
		{"fetch_max_bytes", 1, 32 << 20}, {"fetch_image_max_bytes", 1, 20 << 20},
		{"fetch_timeout_sec", 1, 60}, {"fetch_cache_ttl_sec", 0, 3600},
		{"fetch_cache_max_entries", 0, 128}, {"fetch_cache_max_bytes", 0, 256 << 20},
	} {
		for _, value := range []int64{tc.min - 1, tc.min, tc.max, tc.max + 1} {
			for _, source := range []string{"host", "env"} {
				t.Run(fmt.Sprintf("%s/%s/%d", tc.key, source, value), func(t *testing.T) {
					valid := value >= tc.min && value <= tc.max
					var host map[string]json.RawMessage
					switch source {
					case "host":
						host = map[string]json.RawMessage{tc.key: json.RawMessage(fmt.Sprint(value))}
					case "env":
						t.Setenv("TERVA_EXT_WEB_"+strings.ToUpper(tc.key), fmt.Sprint(value))
					}
					_, err := Resolve(host)
					if (err == nil) != valid {
						t.Fatalf("valid=%v, error=%v", valid, err)
					}
				})
			}
		}
	}
}
