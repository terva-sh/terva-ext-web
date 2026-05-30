// Package search implements the pluggable web_search backend. v1 ships Tavily
// (default) and SearXNG (self-hosted alternate) behind one Provider interface;
// adding Brave/Serper/Exa later is one new file each.
package search

import (
	"context"
	"fmt"
	"strings"

	"git.local.sothr.com/warricksothr/zot-web/internal/config"
)

// Result is one search hit.
type Result struct {
	Title   string
	URL     string
	Snippet string
}

// Provider is a web-search backend.
type Provider interface {
	Search(ctx context.Context, query string, count int) ([]Result, error)
}

// New builds the provider selected by cfg.SearchBackend.
func New(cfg config.Config) (Provider, error) {
	switch cfg.SearchBackend {
	case "tavily":
		if cfg.TavilyAPIKey == "" {
			return nil, fmt.Errorf("tavily backend selected but no API key (set TAVILY_API_KEY or tavily_api_key in config.json)")
		}
		return &tavily{key: cfg.TavilyAPIKey}, nil
	case "searxng":
		if cfg.SearxngURL == "" {
			return nil, fmt.Errorf("searxng backend selected but no instance URL (set ZOT_WEB_SEARXNG_URL or searxng_url in config.json)")
		}
		return &searxng{base: strings.TrimRight(cfg.SearxngURL, "/")}, nil
	default:
		return nil, fmt.Errorf("unknown search backend %q (use \"tavily\" or \"searxng\")", cfg.SearchBackend)
	}
}

// Format renders results as a compact markdown list the model can read and
// chain into web_fetch. URLs are kept so the model can cite and follow them.
func Format(query string, results []Result) string {
	if len(results) == 0 {
		return fmt.Sprintf("No results for %q.", query)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Search results for %q:\n", query)
	for i, r := range results {
		fmt.Fprintf(&b, "\n%d. %s\n   %s\n", i+1, strings.TrimSpace(r.Title), r.URL)
		if s := oneLine(r.Snippet); s != "" {
			fmt.Fprintf(&b, "   %s\n", s)
		}
	}
	return strings.TrimSpace(b.String())
}

func oneLine(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	const max = 300
	if len([]rune(s)) > max {
		s = string([]rune(s)[:max]) + "…"
	}
	return s
}

// clampCount keeps result counts sane.
func clampCount(n int) int {
	if n <= 0 {
		return 5
	}
	if n > 10 {
		return 10
	}
	return n
}
