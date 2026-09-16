package config

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"search_backend":"searxng","searxng_url":"https://example.invalid","fetch_inline_images":true,"fetch_cache_ttl_sec":900,"allow_local_hosts":["legacy.invalid"],"user_agent":"legacy"}`), 0600); err != nil {
		t.Fatal(err)
	}
	c, err := Resolve(dir, "", nil)
	if err != nil || c.SearchBackend != "searxng" {
		t.Fatalf("legacy not preserved: %v", err)
	}
	host := map[string]json.RawMessage{"search_backend": json.RawMessage(`"tavily"`), "fetch_inline_images": json.RawMessage(`false`), "fetch_cache_ttl_sec": json.RawMessage(`0`), "user_agent": json.RawMessage(`""`), "allow_local_hosts": json.RawMessage(`"[]"`)}
	c, err = Resolve(dir, "", host)
	if err != nil || c.SearchBackend != "tavily" || c.FetchInlineImages || c.FetchCacheTTLSec != 0 || c.UserAgent != "" || len(c.AllowLocalHosts) != 0 {
		t.Fatalf("explicit host values lost: %v", err)
	}
	t.Setenv("ZOT_WEB_USER_AGENT", "old-env")
	t.Setenv("TERVA_EXT_WEB_USER_AGENT", "new-env")
	t.Setenv("ZOT_WEB_ALLOW_LOCAL_HOSTS", "old.invalid")
	t.Setenv("TERVA_EXT_WEB_ALLOW_LOCAL_HOSTS", "new.invalid")
	c, err = Resolve(dir, "", host)
	if err != nil || c.UserAgent != "new-env" || !reflect.DeepEqual(c.AllowLocalHosts, []string{"new.invalid"}) {
		t.Fatalf("env precedence failed: %v", err)
	}
	t.Setenv("TERVA_EXT_WEB_USER_AGENT", "")
	c, err = Resolve(dir, "", host)
	if err != nil || c.UserAgent != "" {
		t.Fatal("empty explicit env lost")
	}
}
func TestResolveRejectsInvalidWithoutEcho(t *testing.T) {
	isolateResolveEnv(t)
	for _, tc := range []struct{ key, value string }{
		{"fetch_timeout_sec", "0"}, {"fetch_max_bytes", "33554433"}, {"fetch_inline_images", `"private-marker"`}, {"fetch_cache_ttl_sec", "null"}, {"allow_local_hosts", `"[\\\"bad/private-marker\\\"]"`}, {"searxng_url", `"https://user:private-marker@example.invalid"`}, {"configuration_source", `"private-marker"`},
	} {
		t.Run(tc.key, func(t *testing.T) {
			_, err := Resolve("", "", map[string]json.RawMessage{tc.key: json.RawMessage(tc.value)})
			if err == nil || strings.Contains(err.Error(), "private-marker") {
				t.Fatalf("unsafe/missing error: %v", err)
			}
		})
	}
	t.Setenv("ZOT_WEB_FETCH_TIMEOUT_SEC", "25")
	t.Setenv("TERVA_EXT_WEB_FETCH_TIMEOUT_SEC", "")
	if _, err := Resolve("", "", nil); err == nil {
		t.Fatal("invalid new env fell through")
	}
}
func TestResolveLegacyFileFailuresAndHostMode(t *testing.T) {
	isolateResolveEnv(t)
	data, install := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(data, "config.json"), []byte(`{"tavily_api_key":"private-marker",`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(install, "config.json"), []byte(`{"search_backend":"searxng"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(data, install, nil); err == nil || strings.Contains(err.Error(), "private-marker") {
		t.Fatal("invalid data file not safely rejected")
	}
	c, err := Resolve(data, install, map[string]json.RawMessage{"configuration_source": json.RawMessage(`"host"`)})
	if err != nil || c.SearchBackend != "tavily" {
		t.Fatalf("host mode read legacy: %v", err)
	}
}

func TestSecretImportPrecedenceRollbackAndMissingHost(t *testing.T) {
	isolateResolveEnv(t)
	dir := t.TempDir()
	// Fictional values are generated per test; no installed credential is read.
	legacyKey := "synthetic-legacy-" + filepath.Base(dir)
	hostKey := "synthetic-host-" + filepath.Base(dir)
	envKey := "synthetic-env-" + filepath.Base(dir)
	original, _ := json.Marshal(map[string]string{"tavily_api_key": legacyKey})
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(hostKey)
	values := map[string]json.RawMessage{"tavily_api_key": raw}
	check := func(want string) {
		t.Helper()
		c, err := Resolve(dir, "", values)
		if err != nil {
			t.Fatal(err)
		}
		if c.TavilyAPIKey != want {
			t.Fatal("credential source mismatch")
		}
		b, err := os.ReadFile(path)
		if err != nil || string(b) != string(original) {
			t.Fatal("legacy fixture changed")
		}
	}
	check(legacyKey)
	values["configuration_source"] = json.RawMessage(`"host"`)
	check(hostKey)
	check(hostKey) // opt-in/retry is idempotent
	t.Setenv("TAVILY_API_KEY", envKey)
	check(envKey)
	t.Setenv("TAVILY_API_KEY", "")
	check("")
	if err := os.Unsetenv("TAVILY_API_KEY"); err != nil {
		t.Fatal(err)
	}
	delete(values, "tavily_api_key")
	check("") // absent/undecryptable never resurrects legacy
	values["configuration_source"] = json.RawMessage(`"legacy"`)
	check(legacyKey)
}
