package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"terva.sh/terva/packages/agent/ext"
)

func TestRuntimeUpdateKeepsInflightCacheSeparate(t *testing.T) {
	isolateRuntimeEnv(t)
	for _, change := range []string{"allow_local_hosts", "user_agent", "search_backend"} {
		t.Run(change, func(t *testing.T) {
			started := make(chan struct{})
			release := make(chan struct{})
			var once sync.Once
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				close(started)
				<-release
				_, _ = w.Write([]byte("old content"))
			}))
			defer server.Close()
			defer once.Do(func() { close(release) })
			values := ext.Config{}
			var mu sync.Mutex
			state := runtimeStore{read: func() ext.Config { mu.Lock(); defer mu.Unlock(); return values }, notify: func(s string) { t.Errorf("unexpected notification %s", s) }}
			old := state.snapshot()
			done := make(chan error, 1)
			go func() { _, err := old.fetcher.Raw(context.Background(), server.URL, ""); done <- err }()
			select {
			case <-started:
			case <-time.After(frameTimeoutForRuntime):
				t.Fatal("old fetch did not start")
			}
			mu.Lock()
			values = ext.Config{}
			switch change {
			case "allow_local_hosts":
				values[change] = json.RawMessage(`"[]"`)
			case "user_agent":
				values[change] = json.RawMessage(`"new-agent"`)
			case "search_backend":
				values[change] = json.RawMessage(`"searxng"`)
			}
			mu.Unlock()
			next := state.snapshot()
			if next == old || next.fetcher == old.fetcher {
				t.Fatal("runtime/cache reused")
			}
			once.Do(func() { close(release) })
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(frameTimeoutForRuntime):
				t.Fatal("old fetch stalled")
			}
			if len(next.fetcher.CacheList()) != 0 {
				t.Fatal("old fetch populated new cache")
			}
			if change == "allow_local_hosts" {
				if _, err := next.fetcher.Raw(context.Background(), server.URL, ""); err == nil {
					t.Fatal("tightened allowlist reused old cached result")
				}
			}
		})
	}
}

const frameTimeoutForRuntime = 5 * time.Second

func TestRuntimeRejectsUpdateWithoutLosingWorkingSettings(t *testing.T) {
	isolateRuntimeEnv(t)
	values := ext.Config{}
	var notices []string
	state := runtimeStore{read: func() ext.Config { return values }, notify: func(s string) { notices = append(notices, s) }}
	old := state.snapshot()
	values = ext.Config{"fetch_max_bytes": json.RawMessage(`"private-marker"`)}
	if state.snapshot() != old {
		t.Fatal("rejected update discarded working runtime")
	}
	if len(notices) != 1 || strings.Contains(notices[0], "private-marker") {
		t.Fatal("diagnostic leaked values or missing")
	}
	if state.snapshot() != old || len(notices) != 1 {
		t.Fatal("repeated rejection changed state")
	}
}

func isolateRuntimeEnv(t *testing.T) {
	t.Helper()
	for _, kv := range os.Environ() {
		key, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(key, "TERVA_EXT_WEB_") || strings.HasPrefix(key, "ZOT_WEB_") || key == "TAVILY_API_KEY" {
			t.Setenv(key, "")
			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
		}
	}
}
func TestRuntimeConcurrentUpdates(t *testing.T) {
	isolateRuntimeEnv(t)
	var mu sync.Mutex
	values := ext.Config{}
	state := runtimeStore{read: func() ext.Config { mu.Lock(); defer mu.Unlock(); return values }, notify: func(string) { t.Error("unexpected rejection") }}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				if r := state.snapshot(); r.configErr != nil || r.fetcher == nil {
					t.Error("incoherent snapshot")
				}
			}
		}()
	}
	for i := 0; i < 30; i++ {
		raw, _ := json.Marshal(i)
		mu.Lock()
		values = ext.Config{"fetch_cache_ttl_sec": raw}
		mu.Unlock()
		state.snapshot()
	}
	wg.Wait()
}
