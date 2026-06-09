package fetch

import (
	"context"
	"net"
	"net/url"
	"strings"
	"testing"

	"git.local.sothr.com/warricksothr/zot-web/internal/config"
)

func TestIsBlockedIP(t *testing.T) {
	cases := map[string]bool{
		"127.0.0.1":       true,  // loopback
		"10.1.2.3":        true,  // private
		"192.168.0.1":     true,  // private
		"172.16.5.5":      true,  // private
		"169.254.169.254": true,  // link-local / cloud metadata
		"100.64.1.1":      true,  // CGNAT
		"::1":             true,  // loopback v6
		"fc00::1":         true,  // ULA
		"0.0.0.0":         true,  // unspecified
		"8.8.8.8":         false, // public
		"1.1.1.1":         false, // public
		"93.184.216.34":   false, // public
	}
	for s, want := range cases {
		ip := net.ParseIP(s)
		if ip == nil {
			t.Fatalf("bad test ip %q", s)
		}
		if got := isBlockedIP(ip); got != want {
			t.Errorf("isBlockedIP(%s) = %v, want %v", s, got, want)
		}
	}
}

func TestAllowListPermitted(t *testing.T) {
	a := ParseAllowList([]string{"localhost", "10.0.0.0/24", "192.168.1.50", "grafana.internal"})
	check := func(host, ip string, want bool) {
		t.Helper()
		if got := a.permitted(host, net.ParseIP(ip)); got != want {
			t.Errorf("permitted(%q, %s) = %v, want %v", host, ip, got, want)
		}
	}
	check("example.com", "8.8.8.8", true)       // public always ok
	check("evil.internal", "172.16.0.1", false) // unlisted private blocked
	check("box", "10.0.1.5", false)             // outside allowlisted /24
	check("box", "192.168.1.51", false)         // not the allowlisted exact IP
	check("localhost", "127.0.0.1", true)       // hostname allowlist
	check("grafana.internal", "10.5.5.5", true) // hostname allowlist
	check("GRAFANA.INTERNAL", "10.5.5.5", true) // case-insensitive
	check("box", "10.0.0.5", true)              // inside the /24
	check("box", "192.168.1.50", true)          // exact IP
}

// testClient builds a Client suitable for render/extract unit tests (no network
// is exercised by render itself).
func testClient() *Client {
	return New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5}, ParseAllowList(nil))
}

// TestRenderReadability feeds a realistic article wrapped in nav/script/footer
// chrome and asserts the readability+markdown pipeline keeps the body and drops
// the noise. Assertions hold whether the readability path or the heuristic
// fallback runs (both strip <script> and keep visible text).
func TestRenderReadability(t *testing.T) {
	page := `<html><head><title>Widget Guide</title><style>.a{color:red}</style></head>
<body>
<nav><a href="/">Home</a> <a href="/login">Log in</a></nav>
<script>tracker();</script>
<article>
<h1>All About Widgets</h1>
<p>Widgets are small components that do useful things. This paragraph explains the basics in enough detail to read like real article content.</p>
<p>The second paragraph continues the discussion, giving the readability algorithm enough text to recognize this as the page's main content.</p>
<p>A third paragraph ensures there is sufficient length for extraction to succeed reliably.</p>
</article>
<footer>Copyright 2026 Widget Co.</footer>
</body></html>`
	u, _ := url.Parse("https://example.com/widgets")
	p := testClient().render(u, "text/html", []byte(page))
	out := p.Title + "\n" + p.Markdown
	if !strings.Contains(out, "Widgets are small components") {
		t.Errorf("missing article body in %q", out)
	}
	if !strings.Contains(out, "second paragraph") {
		t.Errorf("missing later paragraph in %q", out)
	}
	if !strings.Contains(out, "Widget Guide") {
		t.Errorf("missing title in %q", out)
	}
	if strings.Contains(out, "tracker()") || strings.Contains(out, "color:red") {
		t.Errorf("script/style not stripped: %q", out)
	}
}

// TestRenderTable confirms the GFM table plugin is active: a real <table>
// becomes a pipe table rather than linearized text.
func TestRenderTable(t *testing.T) {
	page := `<html><body><article>
<h1>Specs</h1>
<p>The following table lists the specifications in a structured form for the reader.</p>
<table>
<tr><th>Name</th><th>Value</th></tr>
<tr><td>Width</td><td>10cm</td></tr>
<tr><td>Height</td><td>20cm</td></tr>
</table>
<p>That concludes the specifications section of this document about the product.</p>
</article></body></html>`
	u, _ := url.Parse("https://example.com/specs")
	p := testClient().render(u, "text/html", []byte(page))
	if !strings.Contains(p.Markdown, "|") || !strings.Contains(p.Markdown, "Width") {
		t.Errorf("expected a pipe table with cell content, got:\n%s", p.Markdown)
	}
}

// TestRenderImages checks that <img> URLs are indexed out to [image:N]
// placeholders (with alt), resolved to absolute URLs, and that data: URIs are
// skipped.
func TestRenderImages(t *testing.T) {
	page := `<html><body><article>
<h1>Gallery</h1>
<p>This article includes a picture to demonstrate the image extraction behavior end to end.</p>
<p><img src="/pics/cat.png" alt="A cat"> some text after the image to keep the paragraph long enough.</p>
<p>Here is an inline data image that must be ignored: <img src="data:image/gif;base64,R0lGOD999"> and more trailing text.</p>
</article></body></html>`
	u, _ := url.Parse("https://example.com/gallery")
	p := testClient().render(u, "text/html", []byte(page))
	if len(p.Images) != 1 {
		t.Fatalf("expected exactly 1 indexed image (data: skipped), got %d: %+v", len(p.Images), p.Images)
	}
	if p.Images[0].URL != "https://example.com/pics/cat.png" {
		t.Errorf("relative src not resolved to absolute: %q", p.Images[0].URL)
	}
	if !strings.Contains(p.Markdown, "[image:1: A cat]") {
		t.Errorf("placeholder with alt missing from markdown:\n%s", p.Markdown)
	}
	if strings.Contains(p.Markdown, "cat.png") || strings.Contains(p.Markdown, "{{IMG") {
		t.Errorf("raw url or sentinel leaked into markdown:\n%s", p.Markdown)
	}
}

// TestRenderInlineImages: with inline images configured, URLs stay in the
// markdown and nothing is indexed.
func TestRenderInlineImages(t *testing.T) {
	c := New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5, FetchInlineImages: true}, ParseAllowList(nil))
	page := `<html><body><article><h1>G</h1>
<p>A paragraph long enough for readability to treat this as the article content here.</p>
<p><img src="https://cdn.example.com/x.png" alt="x"> trailing text to lengthen the paragraph body.</p>
</article></body></html>`
	u, _ := url.Parse("https://example.com/g")
	p := c.render(u, "text/html", []byte(page))
	if len(p.Images) != 0 {
		t.Errorf("inline mode should not index images, got %+v", p.Images)
	}
	if !strings.Contains(p.Markdown, "x.png") {
		t.Errorf("inline mode should keep the image URL in markdown:\n%s", p.Markdown)
	}
}

// TestExtractHeuristicFallback covers the tag-stripper directly: it must drop
// script/style and unescape entities for bodies readability can't handle.
func TestExtractHeuristicFallback(t *testing.T) {
	page := `<html><head><style>.x{color:red}</style></head>` +
		`<body><script>evil()</script><h1>Hello</h1><p>World &amp; more</p></body></html>`
	got := heuristicExtract([]byte(page))
	if !strings.Contains(got, "Hello") {
		t.Errorf("missing visible text in %q", got)
	}
	if !strings.Contains(got, "World & more") {
		t.Errorf("entities not unescaped: %q", got)
	}
	if strings.Contains(got, "evil()") || strings.Contains(got, "color:red") {
		t.Errorf("script/style not stripped: %q", got)
	}
}

// TestRenderNonHTML returns non-HTML bodies verbatim (no title, no images).
func TestRenderNonHTML(t *testing.T) {
	u, _ := url.Parse("https://example.com/data.txt")
	p := testClient().render(u, "text/plain", []byte("  plain body  "))
	if p.Title != "" {
		t.Errorf("non-HTML should have no title, got %q", p.Title)
	}
	if p.Markdown != "plain body" {
		t.Errorf("non-HTML body should pass through trimmed, got %q", p.Markdown)
	}
}

// TestFetchOffsetWindowing drives the metadata block + offset paging using a
// pre-seeded cache entry (so no network is touched).
func TestFetchOffsetWindowing(t *testing.T) {
	c := New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5, FetchCacheMaxEntries: 4, FetchCacheTTLSec: 60}, ParseAllowList(nil))
	u, _ := url.Parse("https://example.com/big")
	body := strings.Repeat("0123456789", 50) // 500 runes
	c.cache.put(page{URL: u.String(), Title: "Big", Markdown: body, ContentType: "text/html; charset=utf-8"})

	out, err := c.Fetch(context.Background(), u.String(), 100, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Chars: 0-100 of 500") {
		t.Errorf("missing/incorrect char metadata:\n%s", out)
	}
	if !strings.Contains(out, "Content-Type: text/html; charset=utf-8") {
		t.Errorf("missing content-type metadata:\n%s", out)
	}
	if !strings.Contains(out, "continue with offset=100") {
		t.Errorf("missing continuation hint:\n%s", out)
	}

	// Final window: no continuation hint when we reach the end.
	out2, err := c.Fetch(context.Background(), u.String(), 100, 450)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out2, "Chars: 450-500 of 500") {
		t.Errorf("incorrect tail window metadata:\n%s", out2)
	}
	if strings.Contains(out2, "continue with offset") {
		t.Errorf("should not advertise more chars at the end:\n%s", out2)
	}
}

// TestFetchBlocksPrivate exercises the real DialContext gate: a loopback target
// with no allowlist must be refused before any connection is made.
func TestFetchBlocksPrivate(t *testing.T) {
	c := New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5}, ParseAllowList(nil))
	_, err := c.Fetch(context.Background(), "http://127.0.0.1:1/", 0, 0)
	if err == nil {
		t.Fatal("expected loopback fetch to be blocked")
	}
	if !strings.Contains(err.Error(), "blocked") {
		t.Fatalf("expected an SSRF block error, got: %v", err)
	}
}

// TestFetchAllowlistAllowsPrivate confirms an allowlisted loopback gets PAST the
// SSRF gate (it then fails to connect to the closed port, which is fine — the
// point is the error is a connection error, not a block).
func TestFetchAllowlistAllowsPrivate(t *testing.T) {
	c := New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5}, ParseAllowList([]string{"127.0.0.1"}))
	_, err := c.Fetch(context.Background(), "http://127.0.0.1:1/", 0, 0)
	if err != nil && strings.Contains(err.Error(), "blocked") {
		t.Fatalf("allowlisted loopback should not be SSRF-blocked: %v", err)
	}
}

func TestFetchRejectsNonHTTPScheme(t *testing.T) {
	c := New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5}, ParseAllowList(nil))
	_, err := c.Fetch(context.Background(), "ftp://example.com/x", 0, 0)
	if err == nil || !strings.Contains(err.Error(), "scheme") {
		t.Fatalf("expected scheme rejection, got: %v", err)
	}
}
