// Command zot-web is a zot extension that gives the agent web tools:
//
//	web_search(query, count?)        -> ranked results (title, url, snippet)
//	web_fetch(url, max_chars?, ...)  -> the page's main content as Markdown
//	web_images(url)                  -> resolve a page's [image:N] placeholders
//	web_links(url)                   -> every link on a page (absolute URL + text)
//	web_fetch_image(url, ...)        -> an image for multimodal viewing / save to disk
//	web_fetch_raw(url, save_path)    -> the page's unrendered source, saved to a file
//
// Search is pluggable (Tavily default, SearXNG alternate). Fetching is
// SSRF-guarded with a configurable local-address allowlist. See README.md.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

const linksSchema = `{
  "type": "object",
  "properties": {
    "url": {"type": "string", "description": "URL of a page (ideally one already retrieved with web_fetch)."}
  },
  "required": ["url"]
}`

const webFetchRawSchema = `{
  "type": "object",
  "properties": {
    "url": {"type": "string", "description": "Absolute http(s) URL to fetch."},
    "save_path": {"type": "string", "description": "Workspace-relative path to write the unrendered page source to (e.g. \"tmp/thread.html\"). Must stay within the workspace; parent directories are created as needed."},
    "overwrite": {"type": "boolean", "description": "Allow overwriting save_path if it already exists (default false)."}
  },
  "required": ["url", "save_path"]
}`

const webFetchImageSchema = `{
  "type": "object",
  "properties": {
    "url": {"type": "string", "description": "Absolute http(s) URL of an image (PNG, JPEG, GIF, or WebP)."},
    "max_dimension": {"type": "integer", "description": "If set, downscale so the image's longest edge is at most this many pixels (preserves aspect ratio, never upscales). Use this to bring an oversized image under the size limit.", "minimum": 1},
    "save_path": {"type": "string", "description": "Optional workspace-relative path to write the image to (e.g. \"assets/logo.png\"). Must stay within the workspace; parent directories are created as needed."},
    "overwrite": {"type": "boolean", "description": "Allow overwriting save_path if it already exists (default false)."},
    "inject": {"type": "boolean", "description": "Whether to return the image to you for viewing (default true). Set false to only download/save it without spending context on the pixels."}
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
		rl       = newRateLimiter(10) // 10 burst, refilled per-tool at different rates
	)
	ensure := func() {
		once.Do(func() {
			cfg := config.Load(e.Host().DataDir)
			fetcher = fetch.New(cfg, fetch.ParseAllowList(cfg.AllowLocalHosts))
			provider, provErr = search.New(cfg, fetcher.HTTPClient())
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
			if !rl.allow("web_search", 5*time.Second) {
				return proto.Errorf("web_search: rate limit reached; wait a few seconds")
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
				return proto.Errorf("search failed: %v", logSSRF(e, err))
			}
			return proto.Text(search.Format(in.Query, results))
		})

	e.Tool("web_fetch",
		"Fetch a web page (http/https) and return its main text content. Results are cached briefly: paging with offset (or repeating the call) within that window reads the same snapshot, so it won't drift mid-read; after the cache expires a re-fetch may differ, with new content typically appended at the end. Private/internal addresses are blocked unless explicitly allowlisted.",
		json.RawMessage(fetchSchema),
		func(args json.RawMessage) proto.Result {
			ensure()
			if !rl.allow("web_fetch", 2*time.Second) {
				return proto.Errorf("web_fetch: rate limit reached; wait a few seconds")
			}
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
				return proto.Errorf("fetch failed: %v", logSSRF(e, err))
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
				return proto.Errorf("web_images failed: %v", logSSRF(e, err))
			}
			return proto.Text(fetch.FormatImages(in.URL, imgs))
		})

	e.Tool("web_links",
		"List every hyperlink on a page (absolute URL plus anchor text). Use to enumerate a page's outbound links without scraping the fetched text yourself. Cheap when the page was recently fetched (served from cache).",
		json.RawMessage(linksSchema),
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
			links, err := fetcher.Links(ctx, in.URL)
			if err != nil {
				return proto.Errorf("web_links failed: %v", logSSRF(e, err))
			}
			return proto.Text(fetch.FormatLinks(in.URL, links))
		})

	e.Tool("web_fetch_raw",
		"Fetch a page and save its UNRENDERED source (HTML/JSON/text, exactly as the server sent it) to a workspace file for you to grep or parse yourself. A fallback for when web_fetch/web_images/web_links don't surface what you need. Served from the same cache as web_fetch. Private/internal addresses are blocked unless explicitly allowlisted.",
		json.RawMessage(webFetchRawSchema),
		func(args json.RawMessage) proto.Result {
			ensure()
			var in struct {
				URL       string `json:"url"`
				SavePath  string `json:"save_path"`
				Overwrite bool   `json:"overwrite"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return proto.Errorf("invalid args: %v", err)
			}
			if strings.TrimSpace(in.URL) == "" {
				return proto.Errorf("url is required")
			}
			if strings.TrimSpace(in.SavePath) == "" {
				return proto.Errorf("save_path is required")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
			defer cancel()
			raw, err := fetcher.Raw(ctx, in.URL)
			if err != nil {
				return proto.Errorf("web_fetch_raw failed: %v", logSSRF(e, err))
			}
			rel, werr := saveToWorkspace(e.Host().CWD, in.SavePath, raw.Body, in.Overwrite)
			if werr != nil {
				return proto.Errorf("fetched the page but could not save it: %v", werr)
			}
			ctype := strings.TrimSpace(raw.ContentType)
			if ctype == "" {
				ctype = "unknown type"
			}
			var meta strings.Builder
			fmt.Fprintf(&meta, "Saved unrendered source to %s\n%s, %d bytes", rel, ctype, len(raw.Body))
			if raw.FinalURL != "" && raw.FinalURL != in.URL {
				fmt.Fprintf(&meta, " (final: %s)", raw.FinalURL)
			}
			if raw.Truncated {
				meta.WriteString("\n…source was capped at the fetch byte limit before saving")
			}
			return proto.Text(meta.String())
		})

	e.Tool("web_fetch_image",
		"Fetch an image (PNG/JPEG/GIF/WebP) by URL and return it for you to view, and/or save it into the workspace. Use max_dimension to downscale a large image. Private/internal addresses are blocked unless explicitly allowlisted.",
		json.RawMessage(webFetchImageSchema),
		func(args json.RawMessage) proto.Result {
			ensure()
			var in struct {
				URL          string `json:"url"`
				MaxDimension int    `json:"max_dimension"`
				SavePath     string `json:"save_path"`
				Overwrite    bool   `json:"overwrite"`
				Inject       *bool  `json:"inject"`
			}
			if err := json.Unmarshal(args, &in); err != nil {
				return proto.Errorf("invalid args: %v", err)
			}
			if strings.TrimSpace(in.URL) == "" {
				return proto.Errorf("url is required")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
			defer cancel()
			img, err := fetcher.FetchImage(ctx, in.URL, in.MaxDimension)
			if err != nil {
				// ImageTooLargeError's message already tells the model how to
				// resubmit (with a suggested max_dimension), so pass it through.
				var tooBig *fetch.ImageTooLargeError
				if errors.As(err, &tooBig) {
					return proto.Errorf("web_fetch_image failed: %v", err)
				}
				return proto.Errorf("web_fetch_image failed: %v", logSSRF(e, err))
			}

			var meta strings.Builder
			fmt.Fprintf(&meta, "Fetched %s", in.URL)
			if img.FinalURL != "" && img.FinalURL != in.URL {
				fmt.Fprintf(&meta, " (final: %s)", img.FinalURL)
			}
			fmt.Fprintf(&meta, "\n%s, %d×%d, %.1f KiB", img.MimeType, img.Width, img.Height, float64(len(img.Data))/1024)
			if img.Resized {
				fmt.Fprintf(&meta, " (resized from %d×%d)", img.OrigW, img.OrigH)
			}

			if strings.TrimSpace(in.SavePath) != "" {
				rel, werr := saveToWorkspace(e.Host().CWD, in.SavePath, img.Data, in.Overwrite)
				if werr != nil {
					return proto.Errorf("fetched the image but could not save it: %v", werr)
				}
				fmt.Fprintf(&meta, "\nSaved to %s", rel)
			}

			inject := in.Inject == nil || *in.Inject
			if inject {
				return proto.Image(img.MimeType, img.Data, meta.String())
			}
			return proto.Text(meta.String())
		})

	if err := e.Run(); err != nil {
		e.Logf("fatal: %v", err)
	}
}

// saveToWorkspace writes data to savePath resolved under the workspace cwd. It
// refuses absolute paths and any path that escapes the workspace, creates
// parent directories within it, and (unless overwrite) refuses to clobber an
// existing file. Returns the cleaned workspace-relative path written.
func saveToWorkspace(cwd, savePath string, data []byte, overwrite bool) (string, error) {
	if strings.TrimSpace(cwd) == "" {
		return "", fmt.Errorf("no workspace directory available to save into")
	}
	if filepath.IsAbs(savePath) {
		return "", fmt.Errorf("save_path must be relative to the workspace, not absolute")
	}
	root, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}
	target := filepath.Join(root, filepath.Clean(savePath))
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("save_path escapes the workspace")
	}
	// Refuse writes into .git/ — a prompt-injected model could overwrite
	// .git/config or other control files.
	if strings.HasPrefix(rel, ".git"+string(filepath.Separator)) || rel == ".git" {
		return "", fmt.Errorf("writing to .git/ is not permitted")
	}
	if !overwrite {
		if _, err := os.Stat(target); err == nil {
			return "", fmt.Errorf("%s already exists (set overwrite=true to replace it)", rel)
		}
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return "", err
	}
	return rel, nil
}

// logSSRF logs the full SSRF block details to the extension log when an error
// chain contains an SSRFBlockedError, and returns the model-safe message.
func logSSRF(e *proto.Extension, err error) string {
	var ssrf *fetch.SSRFBlockedError
	if errors.As(err, &ssrf) {
		e.Logf("%s", ssrf.Full())
		return ssrf.Error()
	}
	return err.Error()
}

// rateLimiter is a simple per-tool token-bucket rate limiter: it allows
// burst tools per toolKey with a refill rate of refillSec seconds.
type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]int
	burst   int
}

func newRateLimiter(burst int) *rateLimiter {
	return &rateLimiter{
		buckets: map[string]int{},
		burst:   burst,
	}
}

// allow reports whether a call for key is within limits, consuming one token.
// The bucket refills by one token every refillSec seconds (lazily on each
// call).
func (rl *rateLimiter) allow(key string, refillSec time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	tokens, ok := rl.buckets[key]
	if !ok {
		tokens = rl.burst
	}
	if tokens <= 0 {
		return false
	}
	rl.buckets[key] = tokens - 1
	// Start a goroutine to refill one token after refillSec.
	go func(k string) {
		time.Sleep(refillSec)
		rl.mu.Lock()
		if n := rl.buckets[k]; n < rl.burst {
			rl.buckets[k] = n + 1
		}
		rl.mu.Unlock()
	}(key)
	return true
}
