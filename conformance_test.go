//go:build conformance

// Protocol conformance harness. It builds the real ./terva-ext-web binary and drives
// it over stdio exactly as a host would — under the supported Terva host
// profile — asserting the basics that an in-process
// unit test cannot:
//
//   - the binary starts and completes the handshake: hello, tool registration,
//     a session_start subscription, and ready;
//   - a tool_call round-trips to a tool_result;
//   - a terva session_start is accepted without wedging the wire;
//   - the process shuts down cleanly on shutdown; and
//   - every byte the extension writes to stdout is a well-formed JSON frame
//     (the stdout-purity invariant — a stray Println would corrupt the wire,
//     and a buffer-capturing in-process test can't catch it).
//
// The profile records the supported Terva floor/current wire contract. The blocked-download regression
// also proves that session switches cannot redirect an in-flight save. The
// subprocess is race-instrumented, so these checks require cgo and a C compiler.
//
// Tagged `conformance` so it stays out of the default unit run (it shells out
// to `go build`). Run with `just conformance` or
// `go test -tags conformance -run Conformance .`; CI runs it as its own step.
package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"terva.sh/terva/packages/agent/extproto"
	"testing"
	"time"
)

// hostProfile captures the supported Terva wire contract.
type hostProfile struct {
	name              string
	protocolVersion   int
	tervaVersion      string
	sendsSessionStart bool
}

var hostProfiles = []hostProfile{
	{name: "terva-0.137.0-floor-and-current", protocolVersion: 6, tervaVersion: "0.137.0", sendsSessionStart: true},
}

func (p hostProfile) helloAck(dataDir string) map[string]any {
	ack := map[string]any{
		"type":             "hello_ack",
		"protocol_version": p.protocolVersion,
		"supported_events": []string{"session_start"},
		"provider":         "anthropic",
		"model":            "test",
		"cwd":              dataDir,
		"data_dir":         dataDir,
	}
	if p.tervaVersion != "" {
		ack["terva_version"] = p.tervaVersion
	}
	return ack
}

var (
	buildOnce sync.Once
	binPath   string
	buildErr  error
)

// webBinary builds ./terva-ext-web once (offline, against vendor/) and returns the
// path to the temp binary.
func webBinary(t *testing.T) string {
	t.Helper()
	buildOnce.Do(func() {
		dir, err := os.MkdirTemp("", "terva-ext-web-conformance")
		if err != nil {
			buildErr = err
			return
		}
		binPath = filepath.Join(dir, "terva-ext-web")
		cmd := exec.Command("go", "build", "-race", "-mod=vendor", "-o", binPath, ".")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			buildErr = fmt.Errorf("build terva-ext-web: %v\n%s", err, stderr.String())
		}
	})
	if buildErr != nil {
		t.Fatalf("%v", buildErr)
	}
	return binPath
}

const frameTimeout = 10 * time.Second

// driver runs one extension subprocess and speaks the host side of the wire.
type driver struct {
	t      *testing.T
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	lines  chan string
	stderr *bytes.Buffer
}

func startExtension(t *testing.T) *driver {
	t.Helper()
	cmd := exec.Command(webBinary(t))
	// Hermetic env: drop provider/config overrides so web_search is deterministically
	// "not configured" (its handler returns an error result without touching
	// the network), and point both home vars at a throwaway dir.
	home := t.TempDir()
	var env []string
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "ZOT_WEB_") || strings.HasPrefix(kv, "TERVA_EXT_WEB_") || strings.HasPrefix(kv, "TAVILY_API_KEY=") {
			continue
		}
		env = append(env, kv)
	}
	cmd.Env = append(env, "ZOT_HOME="+home, "TERVA_HOME="+home)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("stdin pipe: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start extension: %v", err)
	}

	d := &driver{
		t:      t,
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdinPipe),
		lines:  make(chan string, 64),
		stderr: &stderr,
	}
	go func() {
		sc := bufio.NewScanner(stdoutPipe)
		sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
		for sc.Scan() {
			d.lines <- sc.Text()
		}
		close(d.lines)
	}()
	return d
}

func (d *driver) send(frame map[string]any) {
	d.t.Helper()
	b, err := json.Marshal(frame)
	if err != nil {
		d.t.Fatalf("marshal frame: %v", err)
	}
	if _, err := d.stdin.Write(append(b, '\n')); err != nil {
		d.t.Fatalf("write to extension: %v", err)
	}
	if err := d.stdin.Flush(); err != nil {
		d.t.Fatalf("flush to extension: %v", err)
	}
}

// readFrame returns the next stdout frame, failing if the line is not a
// well-formed JSON object with a "type" (the stdout-purity check) or nothing
// arrives in time.
func (d *driver) readFrame() map[string]any {
	d.t.Helper()
	select {
	case line, ok := <-d.lines:
		if !ok {
			d.t.Fatalf("extension closed stdout unexpectedly\nstderr:\n%s", d.stderr.String())
		}
		var f map[string]any
		if err := json.Unmarshal([]byte(line), &f); err != nil {
			d.t.Fatalf("stdout line is not valid JSON (wire corruption): %q: %v", line, err)
		}
		if _, ok := f["type"]; !ok {
			d.t.Fatalf("stdout frame has no \"type\": %q", line)
		}
		return f
	case <-time.After(frameTimeout):
		d.t.Fatalf("timed out waiting for a frame\nstderr:\n%s", d.stderr.String())
		return nil
	}
}

// readUntil reads frames until one of type typ, returning all of them.
func (d *driver) readUntil(typ string) []map[string]any {
	d.t.Helper()
	var got []map[string]any
	for {
		f := d.readFrame()
		got = append(got, f)
		if f["type"] == typ {
			return got
		}
	}
}

// awaitToolResult reads until the tool_result for id, tolerating one-way frames
// (e.g. notify) in between.
func (d *driver) awaitToolResult(id string) map[string]any {
	d.t.Helper()
	for {
		f := d.readFrame()
		if f["type"] == "tool_result" && f["id"] == id {
			return f
		}
		if f["type"] == "shutdown_ack" {
			d.t.Fatalf("extension acked shutdown while awaiting tool_result %q", id)
		}
	}
}

// expectCleanExit drains trailing stdout and asserts the process exits without
// error within a deadline.
func (d *driver) expectCleanExit() {
	d.t.Helper()
	done := make(chan error, 1)
	go func() {
		for range d.lines { // wait for stdout EOF (process closing down)
		}
		done <- d.cmd.Wait()
	}()
	select {
	case err := <-done:
		if err != nil {
			d.t.Errorf("extension exited with error: %v\nstderr:\n%s", err, d.stderr.String())
		}
	case <-time.After(5 * time.Second):
		d.t.Errorf("extension did not exit after shutdown\nstderr:\n%s", d.stderr.String())
		_ = d.cmd.Process.Kill()
	}
}

func TestConformance(t *testing.T) {
	for _, p := range hostProfiles {
		p := p
		t.Run(p.name, func(t *testing.T) {
			d := startExtension(t)
			defer func() { _ = d.cmd.Process.Kill() }() // safety net on early failure

			// 1. Startup is host-agnostic: the extension sends hello +
			//    registrations + subscribe + ready before it ever sees
			//    hello_ack.
			assertStartup(t, d.readUntil("ready"))

			// 2. Identify as this host.
			dataDir := t.TempDir()
			d.send(p.helloAck(dataDir))

			// 3. terva delivers a session_start; a one-way event must not wedge
			//    the wire (we verify by still getting a tool_result below).
			if p.sendsSessionStart {
				d.send(map[string]any{
					"type": "event", "event": "session_start",
					"session_id": "sess-1", "cwd": dataDir, "project_id": "proj-1",
				})
			}

			// 4. A tool_call must round-trip to a tool_result. web_search is
			//    unconfigured here, so the result is an error — that still
			//    exercises real dispatch through to a reply.
			d.send(map[string]any{
				"type": "tool_call", "id": "call-1", "name": "web_search",
				"args": map[string]any{"query": "ping"},
			})
			res := d.awaitToolResult("call-1")
			if _, ok := res["content"]; !ok {
				t.Errorf("tool_result missing content: %v", res)
			}

			// 5. Clean shutdown.
			d.send(map[string]any{"type": "shutdown"})
			d.readUntil("shutdown_ack")
			d.expectCleanExit()
		})
	}
}

func assertStartup(t *testing.T, frames []map[string]any) {
	t.Helper()
	var hello, subscribe map[string]any
	command := false
	toolAuth := map[string]string{}
	for _, f := range frames {
		switch f["type"] {
		case "hello":
			hello = f
		case "register_command":
			command = command || f["name"] == "web-cache"
		case "register_tool":
			name, _ := f["name"].(string)
			auth, _ := f["authority"].(string)
			toolAuth[name] = auth
			if f["read_only"] == true {
				t.Errorf("network tool %q must not be read_only", name)
			}
		case "subscribe":
			subscribe = f
		}
	}

	if !command {
		t.Error("web-cache command not registered")
	}
	if hello == nil {
		t.Fatal("no hello frame in startup")
	}
	if name, _ := hello["name"].(string); name != "web" {
		t.Errorf("hello name = %q, want web", name)
	}
	if hello["min_protocol"] != float64(2) {
		t.Errorf("min_protocol = %v, want 2", hello["min_protocol"])
	}
	if len(toolAuth) != 6 {
		t.Errorf("registered %d tools, want 6", len(toolAuth))
	}
	caps := toStringSet(hello["capabilities"])
	for _, want := range []string{"tools", "events", "commands"} {
		if !caps[want] {
			t.Errorf("hello capabilities missing %q (got %v)", want, hello["capabilities"])
		}
	}

	for _, name := range []string{"web_search", "web_fetch", "web_images", "web_links", "web_fetch_raw", "web_fetch_image"} {
		auth, ok := toolAuth[name]
		if !ok {
			t.Errorf("tool %q not registered", name)
			continue
		}
		if auth != "network-read" {
			t.Errorf("tool %q authority = %q, want network-read", name, auth)
		}
	}

	if subscribe == nil {
		t.Fatal("no subscribe frame (extension must subscribe to session_start)")
	}
	if !toStringSet(subscribe["events"])["session_start"] {
		t.Errorf("subscribe events missing session_start (got %v)", subscribe["events"])
	}
}

func toStringSet(v any) map[string]bool {
	out := map[string]bool{}
	arr, _ := v.([]any)
	for _, e := range arr {
		if s, ok := e.(string); ok {
			out[s] = true
		}
	}
	return out
}

// TestConformanceSessionSwitchSaves exercises the production handlers while a
// download is blocked. No timing sleeps: the HTTP request and command response
// establish that preflight ran in A and session B was processed before saving.
func TestConformanceSessionSwitchSaves(t *testing.T) {
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		tool, contentType string
		body              []byte
	}{
		{"web_fetch_raw", "text/plain", []byte("original raw bytes\n")},
		{"web_fetch_image", "image/png", encoded.Bytes()},
	} {
		t.Run(tc.tool, func(t *testing.T) {
			started, release := make(chan struct{}), make(chan struct{})
			var startOnce, releaseOnce sync.Once
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				startOnce.Do(func() { close(started) })
				select {
				case <-release:
				case <-r.Context().Done():
					return
				}
				w.Header().Set("Content-Type", tc.contentType)
				_, _ = w.Write(tc.body)
			}))
			t.Cleanup(func() {
				releaseOnce.Do(func() { close(release) })
				server.Close()
			})
			d := startExtension(t)
			defer func() { _ = d.cmd.Process.Kill() }()
			assertStartup(t, d.readUntil("ready"))
			first, second := t.TempDir(), t.TempDir()
			profile := hostProfiles[0]
			d.send(profile.helloAck(t.TempDir()))
			d.send(map[string]any{"type": "event", "event": "session_start", "session_id": "first", "cwd": first})
			d.send(map[string]any{
				"type": "tool_call", "id": "download", "name": tc.tool,
				"args": map[string]any{"url": server.URL, "save_path": "downloads/result", "inject": false},
			})
			select {
			case <-started:
			case <-time.After(frameTimeout):
				t.Fatal("download did not reach the HTTP server")
			}
			for _, cwd := range []string{first, second} {
				if _, err := os.Stat(filepath.Join(cwd, "downloads")); !os.IsNotExist(err) {
					t.Fatalf("preflight created downloads in %s: %v", cwd, err)
				}
			}
			d.send(map[string]any{"type": "event", "event": "session_start", "session_id": "second", "cwd": second})
			// Frames are read in order. A response to the next command proves
			// the preceding session event was handled before we release HTTP.
			d.send(map[string]any{"type": "command_invoked", "id": "barrier", "name": "web-cache", "args": ""})
			frames := d.readUntil("command_response")
			if frames[len(frames)-1]["id"] != "barrier" {
				t.Fatal("unexpected command response")
			}
			releaseOnce.Do(func() { close(release) })
			result := d.awaitToolResult("download")
			if result["is_error"] == true {
				t.Errorf("download failed: %v", result)
			}
			got, err := os.ReadFile(filepath.Join(first, "downloads", "result"))
			if err != nil || !bytes.Equal(got, tc.body) {
				t.Errorf("original workspace bytes = %q, err = %v; want %q", got, err, tc.body)
			}
			if _, err := os.Stat(filepath.Join(second, "downloads")); !os.IsNotExist(err) {
				t.Errorf("session switch created a download directory in the new workspace: %v", err)
			}
			d.send(map[string]any{"type": "shutdown"})
			d.readUntil("shutdown_ack")
			d.expectCleanExit()
		})
	}
}

func TestConformanceCommands(t *testing.T) {
	d := startExtension(t)
	defer func() { _ = d.cmd.Process.Kill() }()
	assertStartup(t, d.readUntil("ready"))
	d.send(hostProfiles[0].helloAck(t.TempDir()))
	for _, tc := range []struct{ args, action, text string }{
		{"", "display", "web cache is empty"},
		{"clear", "display", "web cache cleared (0 entries dropped)"},
		{"invalid", "noop", "unknown argument"},
	} {
		d.send(map[string]any{"type": "command_invoked", "id": "cmd", "name": "web-cache", "args": tc.args})
		frames := d.readUntil("command_response")
		f := frames[len(frames)-1]
		field := "display"
		if tc.action == "noop" {
			field = "error"
		}
		text, _ := f[field].(string)
		if f["id"] != "cmd" || f["action"] != tc.action || !strings.Contains(text, tc.text) {
			t.Errorf("command %q: %v", tc.args, f)
		}
	}
	d.send(map[string]any{"type": "shutdown"})
	d.readUntil("shutdown_ack")
	d.expectCleanExit()
}

// Network calls remain concurrent; serializing both downloads would deadlock
// at the second request because neither response is released until both arrive.
func TestConformanceConcurrentDownloads(t *testing.T) {
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	var once sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started <- struct{}{}
		<-release
		_, _ = w.Write([]byte("download"))
	}))
	defer server.Close()
	defer once.Do(func() { close(release) })
	d := startExtension(t)
	defer func() { _ = d.cmd.Process.Kill() }()
	assertStartup(t, d.readUntil("ready"))
	cwd := t.TempDir()
	d.send(hostProfiles[0].helloAck(cwd))
	for _, id := range []string{"one", "two"} {
		d.send(map[string]any{"type": "tool_call", "id": id, "name": "web_fetch_raw", "args": map[string]any{"url": server.URL + "/" + id, "save_path": id}})
	}
	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(frameTimeout):
			t.Fatal("downloads serialized or stalled")
		}
	}
	once.Do(func() { close(release) })
	seen := map[string]bool{}
	for len(seen) < 2 {
		f := d.readFrame()
		if f["type"] != "tool_result" {
			continue
		}
		id, _ := f["id"].(string)
		if (id != "one" && id != "two") || seen[id] || f["is_error"] == true {
			t.Fatalf("unexpected result: %v", f)
		}
		seen[id] = true
		b, err := os.ReadFile(filepath.Join(cwd, id))
		if err != nil || string(b) != "download" {
			t.Fatalf("download %s: %v", id, err)
		}
	}
	d.send(map[string]any{"type": "shutdown"})
	d.readUntil("shutdown_ack")
	d.expectCleanExit()
}

func TestConformanceAllTools(t *testing.T) {
	var pngBody bytes.Buffer
	if err := png.Encode(&pngBody, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/search":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"results":[{"title":"fixture result","url":"https://example.invalid/fixture","content":"fixture snippet"}]}`))
		case "/image":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write(pngBody.Bytes())
		default:
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`<html><title>Fixture</title><body><article><p>Useful fixture content for web retrieval.</p><a href="/other">Fixture link</a><img src="/image" alt="Fixture image"></article></body></html>`))
		}
	}))
	defer server.Close()
	cwd := t.TempDir()
	cfg, _ := json.Marshal(map[string]any{"search_backend": "searxng", "searxng_url": server.URL})
	if err := os.WriteFile(filepath.Join(cwd, "config.json"), cfg, 0600); err != nil {
		t.Fatal(err)
	}
	d := startExtension(t)
	defer func() { _ = d.cmd.Process.Kill() }()
	assertStartup(t, d.readUntil("ready"))
	d.send(hostProfiles[0].helloAck(cwd))
	for _, tc := range []struct {
		name string
		args map[string]any
		want string
	}{
		{"web_search", map[string]any{"query": "fixture"}, "fixture result"},
		{"web_fetch", map[string]any{"url": server.URL}, "fixture content"},
		{"web_images", map[string]any{"url": server.URL}, "/image"},
		{"web_links", map[string]any{"url": server.URL}, "/other"},
		{"web_fetch_raw", map[string]any{"url": server.URL, "save_path": "raw.html"}, "raw.html"},
		{"web_fetch_image", map[string]any{"url": server.URL + "/image"}, "image"},
	} {
		d.send(map[string]any{"type": "tool_call", "id": tc.name, "name": tc.name, "args": tc.args})
		f := d.awaitToolResult(tc.name)
		if f["is_error"] == true {
			t.Fatalf("%s failed: %v", tc.name, f)
		}
		blocks, ok := f["content"].([]any)
		if !ok || len(blocks) == 0 {
			t.Fatalf("%s missing content", tc.name)
		}
		found := false
		for _, b := range blocks {
			block := b.(map[string]any)
			if tc.name == "web_fetch_image" && block["type"] == "image" {
				data, _ := block["data"].(string)
				raw, err := base64.StdEncoding.DecodeString(data)
				if err != nil || !bytes.Equal(raw, pngBody.Bytes()) || block["mime_type"] != "image/png" {
					t.Fatalf("invalid image block: %v", err)
				}
				found = true
			} else if text, ok := block["text"].(string); ok && tc.name != "web_fetch_image" && strings.Contains(text, tc.want) {
				found = true
			}
		}
		if !found {
			t.Errorf("%s content did not contain expected result", tc.name)
		}
	}
	d.send(map[string]any{"type": "shutdown"})
	d.readUntil("shutdown_ack")
	d.expectCleanExit()
}

func TestConformanceFrameRecovery(t *testing.T) {
	d := startExtension(t)
	defer func() { _ = d.cmd.Process.Kill() }()
	assertStartup(t, d.readUntil("ready"))
	d.send(hostProfiles[0].helloAck(t.TempDir()))
	// Send asynchronously with a deadline so a reader regression cannot wedge CI.
	sent := make(chan error, 1)
	go func() {
		_, err := d.stdin.WriteString("{invalid JSON\n" + strings.Repeat("x", extproto.MaxFrameBytes+1) + "\n")
		if err == nil {
			err = d.stdin.Flush()
		}
		sent <- err
	}()
	select {
	case err := <-sent:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(frameTimeout):
		t.Fatal("oversized frame blocked reader")
	}
	d.send(map[string]any{"type": "command_invoked", "id": "recovered", "name": "web-cache", "args": ""})
	frames := d.readUntil("command_response")
	if frames[len(frames)-1]["id"] != "recovered" {
		t.Fatal("frame recovery failed")
	}
	d.send(map[string]any{"type": "shutdown"})
	d.readUntil("shutdown_ack")
	d.expectCleanExit()
}

func TestConformanceShutdownDuringDownload(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(started); <-release }))
	defer server.Close()
	defer close(release)
	d := startExtension(t)
	defer func() { _ = d.cmd.Process.Kill() }()
	assertStartup(t, d.readUntil("ready"))
	cwd := t.TempDir()
	d.send(hostProfiles[0].helloAck(cwd))
	d.send(map[string]any{"type": "tool_call", "id": "blocked", "name": "web_fetch_raw", "args": map[string]any{"url": server.URL, "save_path": "never"}})
	select {
	case <-started:
	case <-time.After(frameTimeout):
		t.Fatal("request not started")
	}
	d.send(map[string]any{"type": "shutdown"})
	d.readUntil("shutdown_ack")
	d.expectCleanExit()
	if _, err := os.Stat(filepath.Join(cwd, "never")); !os.IsNotExist(err) {
		t.Fatalf("shutdown created file: %v", err)
	}
}

func TestConformanceInvalidConfigBlocksEveryNetworkTool(t *testing.T) {
	cwd := t.TempDir()
	if err := os.WriteFile(filepath.Join(cwd, "config.json"), []byte(`{"fetch_max_bytes":"private-marker"}`), 0600); err != nil {
		t.Fatal(err)
	}
	d := startExtension(t)
	defer func() { _ = d.cmd.Process.Kill() }()
	assertStartup(t, d.readUntil("ready"))
	d.send(hostProfiles[0].helloAck(cwd))
	for _, name := range []string{"web_search", "web_fetch", "web_images", "web_links", "web_fetch_raw", "web_fetch_image"} {
		d.send(map[string]any{"type": "tool_call", "id": name, "name": name, "args": map[string]any{}})
		f := d.awaitToolResult(name)
		b, _ := json.Marshal(f)
		if f["is_error"] != true || !bytes.Contains(b, []byte("invalid setting fetch_max_bytes")) || bytes.Contains(b, []byte("private-marker")) {
			t.Errorf("%s did not safely block: %s", name, b)
		}
	}
	d.send(map[string]any{"type": "shutdown"})
	d.readUntil("shutdown_ack")
	d.expectCleanExit()
}

func TestConformanceHostConfigUpdates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(r.UserAgent())) }))
	defer server.Close()
	d := startExtension(t)
	defer func() { _ = d.cmd.Process.Kill() }()
	assertStartup(t, d.readUntil("ready"))
	cwd := t.TempDir()
	ack := hostProfiles[0].helloAck(cwd)
	ack["config"] = map[string]any{"configuration_source": "host", "user_agent": "first"}
	d.send(ack)
	call := func(id string) map[string]any {
		d.send(map[string]any{"type": "tool_call", "id": id, "name": "web_fetch_raw", "args": map[string]any{"url": server.URL, "save_path": id}})
		return d.awaitToolResult(id)
	}
	for _, ua := range []string{"first", "second"} {
		if ua == "second" {
			d.send(map[string]any{"type": "event", "event": "config_update", "config": map[string]any{"configuration_source": "host", "user_agent": ua}})
		}
		if f := call(ua); f["is_error"] == true {
			t.Fatalf("call failed: %v", f)
		}
		b, err := os.ReadFile(filepath.Join(cwd, ua))
		if err != nil || string(b) != ua {
			t.Fatalf("stale config/cache after %s: %q %v", ua, b, err)
		}
	}
	d.send(map[string]any{"type": "event", "event": "config_update", "config": map[string]any{"configuration_source": "host", "allow_local_hosts": "[]"}})
	if f := call("blocked"); f["is_error"] != true {
		t.Fatal("tightened allowlist reused cached response")
	}
	if _, err := os.Stat(filepath.Join(cwd, "blocked")); !os.IsNotExist(err) {
		t.Fatal("denied call wrote a file")
	}
	d.send(map[string]any{"type": "shutdown"})
	d.readUntil("shutdown_ack")
	d.expectCleanExit()
}

func TestConformanceSecretDiagnostics(t *testing.T) {
	for _, malformed := range []bool{false, true} {
		t.Run(fmt.Sprint(malformed), func(t *testing.T) {
			key := "synthetic-" + filepath.Base(t.TempDir())
			d := startExtension(t)
			defer func() { _ = d.cmd.Process.Kill() }()
			assertStartup(t, d.readUntil("ready"))
			ack := hostProfiles[0].helloAck(t.TempDir())
			var value any = key
			if malformed {
				value = map[string]string{"invalid": key}
			}
			ack["config"] = map[string]any{"configuration_source": "host", "tavily_api_key": value}
			d.send(ack)
			// Empty query fails before any provider request, even with a configured key.
			d.send(map[string]any{"type": "tool_call", "id": "secret-check", "name": "web_search", "args": map[string]any{"query": ""}})
			for _, f := range d.readUntil("tool_result") {
				b, _ := json.Marshal(f)
				if bytes.Contains(b, []byte(key)) {
					t.Fatal("credential reached protocol output")
				}
			}
			d.send(map[string]any{"type": "shutdown"})
			d.readUntil("shutdown_ack")
			d.expectCleanExit()
			if strings.Contains(d.stderr.String(), key) {
				t.Fatal("credential reached stderr")
			}
		})
	}
}
