package hostcontract

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"terva.sh/terva/packages/agent/permissions"
	"terva.sh/terva/packages/core"
)

// Exercise the real published host policy using this extension's actual
// manifest, isolated from installed user config and credentials.
func TestAuthorityAndManifestPolicy(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TERVA_HOME", home)
	cwd := t.TempDir()
	manifest, err := os.ReadFile("../../extension.json")
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(home, "extensions", "web")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "extension.json"), manifest, 0600); err != nil {
		t.Fatal(err)
	}
	names := []string{"web_search", "web_fetch", "web_images", "web_links", "web_fetch_raw", "web_fetch_image"}
	if core.IsReadOnlyAuthority(string(core.AuthNetworkRead)) {
		t.Fatal("host treats network as local read")
	}
	for _, mode := range []string{"plan", "ask", "auto-edit", "workspace", "yolo"} {
		for _, override := range []string{"", "allow", "deny"} {
			t.Run(mode+"/"+override, func(t *testing.T) {
				cfg := map[string]any{}
				if override != "" {
					cfg["permissions"] = []map[string]string{{"tool": "web_*", "decision": override}}
				}
				data, err := json.Marshal(cfg)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(home, "config.json"), data, 0600); err != nil {
					t.Fatal(err)
				}
				p, warnings := permissions.BuildPolicy(permissions.Inputs{CWD: cwd, Approval: mode})
				if len(warnings) > 0 || p == nil {
					t.Fatalf("policy missing or invalid: %v", warnings)
				}
				for _, name := range names {
					want := core.VerdictAsk
					switch {
					case mode == "plan", override == "deny":
						want = core.VerdictDeny
					case override == "allow", mode == "yolo":
						want = core.VerdictAllow
					}
					got, _ := p.Evaluate(name, json.RawMessage(`{"url":"https://example.invalid","save_path":"out","overwrite":false}`))
					if got != want {
						t.Errorf("%s: verdict %v, want %v", name, got, want)
					}
				}
				for _, writer := range []string{"web_fetch_raw", "web_fetch_image"} {
					found := false
					for _, r := range p.Rules {
						if r.Tool == writer && r.Decision == core.RuleAsk && r.Source != "user" {
							found = true
						}
					}
					if !found {
						t.Errorf("manifest ask rule missing for %s", writer)
					}
				}
			})
		}
	}
}
