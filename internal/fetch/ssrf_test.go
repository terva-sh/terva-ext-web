package fetch

import (
	"context"
	"net"
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

func TestExtractStripsHTML(t *testing.T) {
	page := `<html><head><style>.x{color:red}</style><title>T</title></head>` +
		`<body><script>evil()</script><h1>Hello</h1><p>World &amp; more</p></body></html>`
	got := extract("text/html", []byte(page))
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
