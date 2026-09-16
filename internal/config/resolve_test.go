package config

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func isolateResolveEnv(t *testing.T) {
	t.Helper()
	for _, item := range os.Environ() {
		key, _, _ := strings.Cut(item, "=")
		if strings.HasPrefix(key, "ZOT_WEB_") || strings.HasPrefix(key, "TERVA_EXT_WEB_") || key == "TAVILY_API_KEY" {
			t.Setenv(key, "")
			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestResolvePrecedence(t *testing.T) {
	isolateResolveEnv(t)
	host := map[string]json.RawMessage{"search_backend": json.RawMessage(`"tavily"`), "fetch_inline_images": json.RawMessage(`false`), "fetch_cache_ttl_sec": json.RawMessage(`0`), "user_agent": json.RawMessage(`""`), "allow_local_hosts": json.RawMessage(`"[]"`)}
	c, err := Resolve(host)
	if err != nil || c.SearchBackend != "tavily" || c.FetchInlineImages || c.FetchCacheTTLSec != 0 || c.UserAgent != "" || len(c.AllowLocalHosts) != 0 {
		t.Fatalf("explicit host values lost: %v", err)
	}
	t.Setenv("TERVA_EXT_WEB_USER_AGENT", "env-client")
	t.Setenv("TERVA_EXT_WEB_ALLOW_LOCAL_HOSTS", " new.invalid, ,10.0.0.0/8 ")
	c, err = Resolve(host)
	if err != nil || c.UserAgent != "env-client" || !reflect.DeepEqual(c.AllowLocalHosts, []string{"new.invalid", "10.0.0.0/8"}) {
		t.Fatalf("env precedence failed: %v", err)
	}
	t.Setenv("TERVA_EXT_WEB_USER_AGENT", "")
	c, err = Resolve(host)
	if err != nil || c.UserAgent != "" {
		t.Fatal("empty explicit env lost")
	}
}

func TestResolveIgnoresRetiredOverrides(t *testing.T) {
	isolateResolveEnv(t)
	want, err := Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, suffix := range []string{"SEARCH_BACKEND", "TAVILY_API_KEY", "SEARXNG_URL", "USER_AGENT", "ALLOW_LOCAL_HOSTS", "FETCH_MAX_BYTES", "FETCH_IMAGE_MAX_BYTES", "FETCH_TIMEOUT_SEC", "FETCH_INLINE_IMAGES", "FETCH_CACHE_TTL_SEC", "FETCH_CACHE_MAX_ENTRIES", "FETCH_CACHE_MAX_BYTES", "CONFIGURATION_SOURCE"} {
		t.Setenv("ZOT_WEB_"+suffix, "ignored-value")
	}
	t.Setenv("TERVA_EXT_WEB_CONFIGURATION_SOURCE", "legacy")
	got, err := Resolve(map[string]json.RawMessage{"configuration_source": json.RawMessage(`"legacy"`)})
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("retired inputs changed configuration: %v", err)
	}
}

func TestResolveRejectsInvalidWithoutEcho(t *testing.T) {
	isolateResolveEnv(t)
	for _, tc := range []struct{ key, value string }{
		{"fetch_timeout_sec", "0"}, {"fetch_max_bytes", "33554433"}, {"fetch_inline_images", `"private-marker"`}, {"fetch_cache_ttl_sec", "null"}, {"allow_local_hosts", `"[\\\"bad/private-marker\\\"]"`}, {"searxng_url", `"https://user:private-marker@example.invalid"`},
	} {
		t.Run(tc.key, func(t *testing.T) {
			_, err := Resolve(map[string]json.RawMessage{tc.key: json.RawMessage(tc.value)})
			if err == nil || strings.Contains(err.Error(), "private-marker") {
				t.Fatalf("unsafe/missing error: %v", err)
			}
		})
	}
	t.Setenv("TERVA_EXT_WEB_FETCH_TIMEOUT_SEC", "")
	if _, err := Resolve(nil); err == nil {
		t.Fatal("invalid environment override fell through")
	}
}

func TestSecretHostAndEnvironmentPrecedence(t *testing.T) {
	isolateResolveEnv(t)
	raw, _ := json.Marshal("synthetic-host-secret")
	values := map[string]json.RawMessage{"tavily_api_key": raw}
	check := func(want string) {
		t.Helper()
		c, err := Resolve(values)
		if err != nil || c.TavilyAPIKey != want {
			t.Fatalf("credential source mismatch: %v", err)
		}
	}
	check("synthetic-host-secret")
	t.Setenv("TAVILY_API_KEY", "synthetic-env-secret")
	check("synthetic-env-secret")
	t.Setenv("TAVILY_API_KEY", "")
	check("")
	if err := os.Unsetenv("TAVILY_API_KEY"); err != nil {
		t.Fatal(err)
	}
	delete(values, "tavily_api_key")
	check("") // Missing or undecryptable host keys remain unconfigured.
}
