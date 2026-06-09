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

// TestExtractReadability feeds a realistic article wrapped in nav/script/footer
// chrome and asserts the readability+markdown pipeline keeps the body and drops
// the noise. Assertions hold whether the readability path or the heuristic
// fallback runs (both strip <script> and keep visible text).
func TestExtractReadability(t *testing.T) {
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
	title, text := extract(u, "text/html", []byte(page))
	out := title + "\n" + text
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

// TestExtractNonHTML returns non-HTML bodies verbatim (no title).
func TestExtractNonHTML(t *testing.T) {
	u, _ := url.Parse("https://example.com/data.txt")
	title, text := extract(u, "text/plain", []byte("  plain body  "))
	if title != "" {
		t.Errorf("non-HTML should have no title, got %q", title)
	}
	if text != "plain body" {
		t.Errorf("non-HTML body should pass through trimmed, got %q", text)
	}
}

// TestFetchBlocksPrivate exercises the real DialContext gate: a loopback target
// with no allowlist must be refused before any connection is made.
func TestFetchBlocksPrivate(t *testing.T) {
	c := New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5}, ParseAllowList(nil))
	_, err := c.Fetch(context.Background(), "http://127.0.0.1:1/", 0)
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
	_, err := c.Fetch(context.Background(), "http://127.0.0.1:1/", 0)
	if err != nil && strings.Contains(err.Error(), "blocked") {
		t.Fatalf("allowlisted loopback should not be SSRF-blocked: %v", err)
	}
}

func TestFetchRejectsNonHTTPScheme(t *testing.T) {
	c := New(config.Config{FetchMaxBytes: 1 << 20, FetchTimeoutSec: 5}, ParseAllowList(nil))
	_, err := c.Fetch(context.Background(), "ftp://example.com/x", 0)
	if err == nil || !strings.Contains(err.Error(), "scheme") {
		t.Fatalf("expected scheme rejection, got: %v", err)
	}
}
