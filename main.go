// Command zot-web is a zot extension that gives the agent two tools:
//
//	web_search(query, count?)  -> ranked results (title, url, snippet)
//	web_fetch(url, max_chars?) -> the page's main content as text
//
// Search is pluggable (Tavily default, SearXNG alternate). Fetch is SSRF-guarded
// with a configurable local-address allowlist. See README.md.
package main

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"git.local.sothr.com/warricksothr/zot-web/internal/config"
	"git.local.sothr.com/warricksothr/zot-web/internal/fetch"
	"git.local.sothr.com/warricksothr/zot-web/internal/proto"
	"git.local.sothr.com/warricksothr/zot-web/internal/search"
)

const searchSchema = `{
  "type": "object",
  "properties": {
    "query": {"type": "string", "description": "The search query."},
    "count": {"type": "integer", "description": "Number of results (default 5, max 10).", "minimum": 1, "maximum": 10}
  },
  "required": ["query"]
}`

const fetchSchema = `{
  "type": "object",
  "properties": {
    "url": {"type": "string", "description": "Absolute http(s) URL to fetch."},
    "max_chars": {"type": "integer", "description": "Max characters of extracted text to return (default 20000)."},
    "offset": {"type": "integer", "description": "Skip this many characters into the page, to continue reading after a previous truncated fetch (default 0).", "minimum": 0}
  },
  "required": ["url"]
}`

const imagesSchema = `{
  "type": "object",
  "properties": {
    "url": {"type": "string", "description": "URL of a page already retrieved with web_fetch."}
  },
  "required": ["url"]
}`

func main() {
	e := proto.New("web", "0.1.0")

	// Providers are built lazily on first tool call, by which point the
	// hello_ack (and thus data_dir for config.json) has arrived.
	var (
		once     sync.Once
		provider search.Provider
		provErr  error
		fetcher  *fetch.Client
	)
	ensure := func() {
		once.Do(func() {
			cfg := config.Load(e.Host().DataDir)
			provider, provErr = search.New(cfg)
			fetcher = fetch.New(cfg, fetch.ParseAllowList(cfg.AllowLocalHosts))
		})
	}

	e.Tool("web_search",
		"Search the web and return ranked results (title, URL, snippet). Use for current events, facts, documentation, or to find pages to read with web_fetch.",
		json.RawMessage(searchSchema),
		func(args json.RawMessage) proto.Result {
			ensure()
			if provErr != nil {
				return proto.Errorf("web_search is not configured: %v", provErr)
			}
			var in struct {
				Query string `json:"query"`
				Count int    `json:"count"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return proto.Errorf("invalid args: %v", err)
			}
			if strings.TrimSpace(in.Query) == "" {
				return proto.Errorf("query is required")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
			defer cancel()
			results, err := provider.Search(ctx, in.Query, in.Count)
			if err != nil {
				return proto.Errorf("search failed: %v", err)
			}
			return proto.Text(search.Format(in.Query, results))
		})

	e.Tool("web_fetch",
		"Fetch a web page (http/https) and return its main text content. Private/internal addresses are blocked unless explicitly allowlisted.",
		json.RawMessage(fetchSchema),
		func(args json.RawMessage) proto.Result {
			ensure()
			var in struct {
				URL      string `json:"url"`
				MaxChars int    `json:"max_chars"`
				Offset   int    `json:"offset"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return proto.Errorf("invalid args: %v", err)
			}
			if strings.TrimSpace(in.URL) == "" {
				return proto.Errorf("url is required")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
			defer cancel()
			text, err := fetcher.Fetch(ctx, in.URL, in.MaxChars, in.Offset)
			if err != nil {
				return proto.Errorf("fetch failed: %v", err)
			}
			return proto.Text(text)
		})

	e.Tool("web_images",
		"List the image URLs on a page that web_fetch represented as [image:N] placeholders. Cheap when the page was recently fetched (it is served from cache).",
		json.RawMessage(imagesSchema),
		func(args json.RawMessage) proto.Result {
			ensure()
			var in struct {
				URL string `json:"url"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return proto.Errorf("invalid args: %v", err)
			}
			if strings.TrimSpace(in.URL) == "" {
				return proto.Errorf("url is required")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
			defer cancel()
			imgs, err := fetcher.Images(ctx, in.URL)
			if err != nil {
				return proto.Errorf("web_images failed: %v", err)
			}
			return proto.Text(fetch.FormatImages(in.URL, imgs))
		})

	if err := e.Run(); err != nil {
		e.Logf("fatal: %v", err)
	}
}
