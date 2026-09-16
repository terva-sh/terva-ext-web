package hostcontract

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"terva.sh/terva/packages/agent/extdriver"
	"terva.sh/terva/packages/agent/extproto"
)

type quietHooks struct{}

func (quietHooks) Notify(string, string, string)                                           {}
func (quietHooks) Submit(string)                                                           {}
func (quietHooks) SubmitSlash(string)                                                      {}
func (quietHooks) Insert(string)                                                           {}
func (quietHooks) Display(string, string)                                                  {}
func (quietHooks) ClearNotes(string)                                                       {}
func (quietHooks) OpenPanel(string, extproto.PanelSpec)                                    {}
func (quietHooks) UpdatePanel(string, string, string, []string, string, []extproto.Widget) {}
func (quietHooks) ClosePanel(string, string)                                               {}
func (quietHooks) RefreshStatus()                                                          {}
func (quietHooks) RefreshContext()                                                         {}
func (quietHooks) RefreshTools()                                                           {}

// This runs the published host's actual spawning/handshake/event/correlation
// code, not a simulated hello_ack. Full CLI installation remains a release test.
func TestPublishedHostDriverLaunch(t *testing.T) {
	for _, kv := range os.Environ() {
		key, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(key, "ZOT_WEB_") || strings.HasPrefix(key, "TERVA_EXT_WEB_") || key == "TAVILY_API_KEY" {
			t.Setenv(key, "")
			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
		}
	}
	home, cwd, install := t.TempDir(), t.TempDir(), t.TempDir()
	t.Setenv("TERVA_HOME", home)
	t.Setenv("ZOT_HOME", home)
	bin := "terva-ext-web"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-race", "-mod=vendor", "-o", filepath.Join(install, bin), ".")
	cmd.Dir = "../.."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	for _, name := range []string{"run.sh", "extension.json"} {
		b, err := os.ReadFile(filepath.Join("../..", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(install, name), b, 0700); err != nil {
			t.Fatal(err)
		}
	}
	b, err := os.ReadFile(filepath.Join(install, "extension.json"))
	if err != nil {
		t.Fatal(err)
	}
	var mf extdriver.Manifest
	if err := json.Unmarshal(b, &mf); err != nil {
		t.Fatal(err)
	}
	d := extdriver.New(home, cwd, "0.137.0", "fixture", "fixture", quietHooks{})
	var malformed atomic.Int32
	d.SetOnMalformedFrame(func(_, _, _ string) { malformed.Add(1) })
	if err := d.Load(ctx, install, mf); err != nil {
		t.Fatal(err)
	}
	defer d.Stop(3 * time.Second)
	d.WaitForReady(10 * time.Second)
	if d.ReadyCount() != 1 {
		t.Fatal("extension not ready")
	}
	tools := d.Tools()
	if len(tools) != 6 {
		t.Fatalf("registered %d tools", len(tools))
	}
	for _, tool := range tools {
		if tool.Authority != "network-read" || tool.ReadOnly {
			t.Errorf("incorrect authority on %s", tool.Name)
		}
	}
	result, err := d.Invoke(ctx, "web-cache", "clear", 5*time.Second)
	if err != nil || result.Action != "display" {
		t.Fatalf("command: %v %v", result, err)
	}
	next := t.TempDir()
	d.EmitEvent(extproto.EventFromHost{Type: "event", Event: "session_start", SessionID: "next", CWD: next})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("host fixture")) }))
	defer server.Close()
	args, _ := json.Marshal(map[string]any{"url": server.URL, "save_path": "result"})
	res, err := d.InvokeTool(ctx, "web_fetch_raw", args, 10*time.Second)
	if err != nil || res.IsError {
		t.Fatalf("fetch: %v %v", res, err)
	}
	saved, err := os.ReadFile(filepath.Join(next, "result"))
	if err != nil || string(saved) != "host fixture" {
		t.Fatalf("host session save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cwd, "result")); !os.IsNotExist(err) {
		t.Fatal("saved in old session")
	}
	d.Stop(3 * time.Second)
	if malformed.Load() != 0 {
		t.Fatalf("host observed %d malformed stdout frames", malformed.Load())
	}
	t.Log("published Terva v0.137.0 host driver (floor/current): launcher, ready, six tools, command, network save, session, shutdown passed")
}
