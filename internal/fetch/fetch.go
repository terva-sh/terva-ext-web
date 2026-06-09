// Package fetch implements web_fetch: an SSRF-guarded HTTP client plus
// main-content extraction. Because the model chooses the URL, fetching is the
// extension's main attack surface — every connection is validated against the
// private-range block (with the configurable local allowlist) in a custom
// DialContext that dials the validated IP and re-runs on each redirect hop.
package fetch

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	readability "github.com/go-shiori/go-readability"
	xhtml "golang.org/x/net/html"

	"git.local.sothr.com/warricksothr/zot-web/internal/config"
)

const userAgent = "zot-web/0.1 (+https://git.local.sothr.com/warricksothr/zot-web)"

// Client is a reusable SSRF-guarded fetcher. Rendered pages are cached so a
// web_images call following a web_fetch needs no network.
type Client struct {
	http         *http.Client
	maxBytes     int64
	allow        AllowList
	inlineImages bool
	cache        *cache
}

// New builds a Client whose dialer refuses private/reserved destinations unless
// the allowlist permits them.
func New(cfg config.Config, allow AllowList) *Client {
	c := &Client{
		maxBytes:     cfg.FetchMaxBytes,
		allow:        allow,
		inlineImages: cfg.FetchInlineImages,
		cache:        newCache(time.Duration(cfg.FetchCacheTTLSec)*time.Second, cfg.FetchCacheMaxEntries),
	}
	dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}

	tr := &http.Transport{
		// Validate at dial time: resolve the host ourselves, pick the first
		// permitted IP, and connect to THAT ip — closing the DNS-rebinding /
		// TOCTOU gap. The transport re-invokes this for every redirect host.
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			for _, ipa := range ips {
				if c.allow.permitted(host, ipa.IP) {
					return dialer.DialContext(ctx, network, net.JoinHostPort(ipa.IP.String(), port))
				}
			}
			return nil, fmt.Errorf("blocked: %q resolves only to private/reserved addresses; add it to allow_local_hosts to permit", host)
		},
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		MaxIdleConns:          10,
	}
	c.http = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(cfg.FetchTimeoutSec) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}
	return c
}

// Fetch retrieves raw (http/https only) and returns its main content as
// Markdown, prefixed with a metadata block. maxChars caps the returned window
// (default 20000); offset skips that many characters into the rendered page so
// callers can page through dense documents. The page is rendered once and
// cached; image URLs are replaced with `[image:N]` placeholders (unless inline
// images are configured) and retrievable via Images.
func (c *Client) Fetch(ctx context.Context, raw string, maxChars, offset int) (string, error) {
	u, err := parseURL(raw)
	if err != nil {
		return "", err
	}
	p, err := c.load(ctx, u)
	if err != nil {
		return "", err
	}

	if maxChars <= 0 {
		maxChars = 20000
	}
	full := []rune(p.Markdown)
	total := len(full)
	start := min(max(offset, 0), total)
	end := min(start+maxChars, total)
	window := string(full[start:end])

	var b strings.Builder
	// Header: article title (when readability found one) above the source URL.
	if p.Title != "" {
		fmt.Fprintf(&b, "# %s\n%s\n", p.Title, u.String())
	} else {
		fmt.Fprintf(&b, "# %s\n", u.String())
	}
	if p.FinalURL != "" && p.FinalURL != u.String() {
		fmt.Fprintf(&b, "Final-URL: %s\n", p.FinalURL)
	}
	if p.ContentType != "" {
		fmt.Fprintf(&b, "Content-Type: %s\n", p.ContentType)
	}
	fmt.Fprintf(&b, "Chars: %d-%d of %d\n", start, end, total)
	if !c.inlineImages && len(p.Images) > 0 {
		fmt.Fprintf(&b, "Images: %d (shown as [image:N]; resolve with web_images)\n", len(p.Images))
	}
	b.WriteString("\n")
	b.WriteString(window)

	if end < total {
		fmt.Fprintf(&b, "\n\n…[%d more chars; continue with offset=%d]", total-end, end)
	} else if p.BodyTruncated {
		b.WriteString("\n\n…[the source response was capped at the byte limit before rendering]")
	}
	return b.String(), nil
}

// Images returns the images found on raw (resolved to absolute URLs). It serves
// a cached render when available, fetching only on a cold cache.
func (c *Client) Images(ctx context.Context, raw string) ([]Image, error) {
	u, err := parseURL(raw)
	if err != nil {
		return nil, err
	}
	p, err := c.load(ctx, u)
	if err != nil {
		return nil, err
	}
	return p.Images, nil
}

// parseURL validates and normalizes a model-supplied URL.
func parseURL(raw string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme %q (only http/https)", u.Scheme)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("url has no host")
	}
	return u, nil
}

// load returns the rendered page for u, from cache when fresh, otherwise by
// fetching and rendering (and caching the result).
func (c *Client) load(ctx context.Context, u *url.URL) (page, error) {
	key := u.String()
	if p, ok := c.cache.get(key); ok {
		return p, nil
	}
	f, err := c.download(ctx, u)
	if err != nil {
		return page{}, err
	}
	p := c.render(u, f.contentType, f.body)
	p.URL = key
	p.FinalURL = f.finalURL
	p.ContentType = f.contentType
	p.Status = f.status
	p.BodyTruncated = f.truncated
	c.cache.put(p)
	return p, nil
}

// fetched is the raw result of an SSRF-guarded GET.
type fetched struct {
	body        []byte
	truncated   bool // body hit the byte cap
	contentType string
	finalURL    string // after redirects
	status      int
}

// download performs the SSRF-guarded GET and returns the (byte-capped) body.
func (c *Client) download(ctx context.Context, u *url.URL) (fetched, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fetched{}, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,text/plain;q=0.9,*/*;q=0.8")

	resp, err := c.http.Do(req)
	if err != nil {
		return fetched{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fetched{}, fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, u)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, c.maxBytes+1))
	if err != nil {
		return fetched{}, err
	}
	f := fetched{
		body:        body,
		contentType: resp.Header.Get("Content-Type"),
		finalURL:    resp.Request.URL.String(),
		status:      resp.StatusCode,
	}
	if int64(len(body)) > c.maxBytes {
		f.body = body[:c.maxBytes]
		f.truncated = true
	}
	return f, nil
}

// render turns a response body into a page: readability isolates the main
// article, image URLs are indexed out to `[image:N]` placeholders, and
// html-to-markdown (with GFM tables) renders the result. On any failure it
// falls back to the heuristic tag-stripper. Non-HTML bodies pass through.
func (c *Client) render(u *url.URL, contentType string, body []byte) page {
	ct := strings.ToLower(contentType)
	isHTML := strings.Contains(ct, "html") ||
		(ct == "" && strings.Contains(strings.ToLower(string(body)), "<html"))
	if !isHTML {
		return page{Markdown: strings.TrimSpace(string(body))}
	}

	art, err := readability.FromReader(bytes.NewReader(body), u)
	if err == nil {
		node := art.Node
		if node == nil {
			node, _ = xhtml.Parse(strings.NewReader(art.Content))
		}
		if node != nil {
			var images []Image
			if !c.inlineImages {
				images = indexImages(node, u)
			}
			if md, err := convertNode(node); err == nil {
				if md = applyPlaceholders(strings.TrimSpace(md), images); md != "" {
					return page{Title: strings.TrimSpace(art.Title), Markdown: md, Images: images}
				}
			}
		}
		// Readability found content but markdown conversion produced nothing;
		// use its plain-text rendering rather than dropping to the heuristic.
		if t := strings.TrimSpace(art.TextContent); t != "" {
			return page{Title: strings.TrimSpace(art.Title), Markdown: t}
		}
	}
	return page{Markdown: heuristicExtract(body)}
}

// convertNode renders an HTML node to Markdown with CommonMark + GFM tables.
// Tables use mirror span cells so rowspan/colspan headers (e.g. infoboxes)
// repeat their value into spanned cells rather than leaving blanks.
func convertNode(n *xhtml.Node) (string, error) {
	conv := converter.NewConverter(converter.WithPlugins(
		base.NewBasePlugin(),
		commonmark.NewCommonmarkPlugin(),
		table.NewTablePlugin(
			table.WithSpanCellBehavior(table.SpanBehaviorMirror),
			table.WithPresentationTables(false),
		),
	))
	b, err := conv.ConvertNode(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// heuristicExtract is the fallback tag-stripper for bodies readability can't
// parse: it removes script/style, drops remaining tags, and unescapes entities.
var (
	reScriptStyle = regexp.MustCompile(`(?is)<(?:script|style|noscript|template)\b[^>]*>.*?</(?:script|style|noscript|template)\s*>`)
	reTag         = regexp.MustCompile(`(?s)<[^>]+>`)
	reInlineWS    = regexp.MustCompile(`[ \t]+`)
	reBlankLines  = regexp.MustCompile(`\n{3,}`)
)

func heuristicExtract(body []byte) string {
	s := string(body)
	s = reScriptStyle.ReplaceAllString(s, "\n")
	s = reTag.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = reInlineWS.ReplaceAllString(s, " ")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	s = strings.Join(lines, "\n")
	s = reBlankLines.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}
