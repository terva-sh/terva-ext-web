package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"terva-ext-web/internal/version"
)

// TestManifestVersionMatchesCode pins extension.json's version (what the host
// shows in `ext list`) equal to internal/version.Version (the hello frame, the
// User-Agent, --version). They are bumped together at release; this guard fails
// the build if they drift, which is how they silently disagreed before.
func TestManifestVersionMatchesCode(t *testing.T) {
	b, err := os.ReadFile("extension.json")
	if err != nil {
		t.Fatalf("read extension.json: %v", err)
	}
	var m struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("parse extension.json: %v", err)
	}
	if m.Version != version.Version {
		t.Errorf("extension.json version %q != internal/version.Version %q — bump them together", m.Version, version.Version)
	}
}

func TestSaveToWorkspaceWritesUnderCWD(t *testing.T) {
	cwd := t.TempDir()
	rel, err := saveToWorkspace(cwd, "assets/logo.png", []byte("png"), false)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if rel != filepath.Join("assets", "logo.png") {
		t.Errorf("rel = %q", rel)
	}
	got, err := os.ReadFile(filepath.Join(cwd, "assets", "logo.png"))
	if err != nil || string(got) != "png" {
		t.Errorf("file content = %q, err = %v", got, err)
	}
}

func TestSaveToWorkspaceRejectsEscape(t *testing.T) {
	cwd := t.TempDir()
	for _, p := range []string{"../escape.png", "a/../../escape.png", "/etc/passwd"} {
		if _, err := saveToWorkspace(cwd, p, []byte("x"), false); err == nil {
			t.Errorf("save_path %q should have been rejected", p)
		}
	}
	// Nothing should have been written outside cwd.
	if _, err := os.Stat(filepath.Join(filepath.Dir(cwd), "escape.png")); err == nil {
		t.Error("a file escaped the workspace")
	}
}

func TestSaveToWorkspaceOverwritePolicy(t *testing.T) {
	cwd := t.TempDir()
	if _, err := saveToWorkspace(cwd, "f.png", []byte("v1"), false); err != nil {
		t.Fatalf("first save: %v", err)
	}
	if _, err := saveToWorkspace(cwd, "f.png", []byte("v2"), false); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected no-clobber error, got %v", err)
	}
	if _, err := saveToWorkspace(cwd, "f.png", []byte("v2"), true); err != nil {
		t.Fatalf("overwrite save: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(cwd, "f.png"))
	if string(got) != "v2" {
		t.Errorf("after overwrite content = %q, want v2", got)
	}
}

func TestSaveToWorkspaceNoCWD(t *testing.T) {
	if _, err := saveToWorkspace("", "f.png", []byte("x"), false); err == nil {
		t.Error("empty cwd should be rejected")
	}
}

func TestSaveToWorkspaceRejectsGitDir(t *testing.T) {
	cwd := t.TempDir()
	// .git/ itself
	if _, err := saveToWorkspace(cwd, ".git", []byte("x"), false); err == nil {
		t.Error("should reject bare .git path")
	}
	// .git/config
	if _, err := saveToWorkspace(cwd, ".git/config", []byte("x"), false); err == nil {
		t.Error("should reject .git/config")
	}
	// .git/hooks/some-hook
	if _, err := saveToWorkspace(cwd, ".git/hooks/pre-commit", []byte("x"), false); err == nil {
		t.Error("should reject files under .git/")
	}
	// But .gitignore or .gitattributes should still work
	if _, err := saveToWorkspace(cwd, ".gitignore", []byte("x"), false); err != nil {
		t.Errorf(".gitignore should be allowed, got: %v", err)
	}
}

func TestSaveToWorkspaceRejectsSymlinkParentEscape(t *testing.T) {
	cwd := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(cwd, "out")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := saveToWorkspace(cwd, filepath.Join("out", "file.txt"), []byte("x"), false); err == nil {
		t.Fatal("symlinked parent should be rejected")
	}
	if _, err := os.Stat(filepath.Join(outside, "file.txt")); err == nil {
		t.Fatal("file was written through symlink outside workspace")
	}
}

func TestSaveToWorkspaceRejectsFinalSymlink(t *testing.T) {
	cwd := t.TempDir()
	outside := filepath.Join(t.TempDir(), "target.txt")
	if err := os.WriteFile(outside, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(cwd, "link.txt")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := saveToWorkspace(cwd, "link.txt", []byte("replace"), true); err == nil {
		t.Fatal("final symlink should be rejected even with overwrite=true")
	}
	got, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "original" {
		t.Fatalf("symlink target was modified: %q", got)
	}
}

func TestCheckSavePathMatchesSavePolicy(t *testing.T) {
	cwd := t.TempDir()
	for _, p := range []string{"../escape.png", "a/../../escape.png", "/etc/passwd", ".git/config"} {
		if err := checkSavePath(cwd, p, false); err == nil {
			t.Errorf("preflight should reject %q", p)
		}
	}
	if err := os.WriteFile(filepath.Join(cwd, "f.png"), []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkSavePath(cwd, "f.png", false); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("preflight no-clobber: err = %v", err)
	}
	if err := checkSavePath(cwd, "f.png", true); err != nil {
		t.Errorf("preflight with overwrite should pass: %v", err)
	}
	if err := checkSavePath(cwd, "new/dir/f.png", false); err != nil {
		t.Errorf("preflight of a fresh nested path should pass: %v", err)
	}
}

func TestCheckSavePathHasNoSideEffects(t *testing.T) {
	cwd := t.TempDir()
	if err := checkSavePath(cwd, "a/b/c.png", false); err != nil {
		t.Fatalf("preflight: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cwd, "a")); !os.IsNotExist(err) {
		t.Error("preflight created parent directories")
	}
}

func TestCheckSavePathRejectsSymlinks(t *testing.T) {
	cwd := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(cwd, "out")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := checkSavePath(cwd, "out/file.txt", false); err == nil {
		t.Error("preflight should reject a symlinked parent")
	}
	if err := os.Symlink(filepath.Join(outside, "t.txt"), filepath.Join(cwd, "link.txt")); err != nil {
		t.Skip("symlinks unavailable")
	}
	if err := checkSavePath(cwd, "link.txt", true); err == nil {
		t.Error("preflight should reject a symlink target")
	}
}

func TestVersionString(t *testing.T) {
	got := versionString()
	if !strings.HasPrefix(got, "terva-ext-web "+version.Version) {
		t.Errorf("versionString() = %q, want prefix %q", got, "terva-ext-web "+version.Version)
	}
	if !strings.Contains(got, runtime.GOOS+"/"+runtime.GOARCH) {
		t.Errorf("versionString() = %q, missing platform", got)
	}
}

// Manifest defaults would erase the distinction between explicit host settings
// and application fallbacks; keep that provenance invariant reviewable.
func TestManifestConfigHasNoDefaults(t *testing.T) {
	b, err := os.ReadFile("extension.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Config []map[string]any `json:"config"`
	}
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Config) < 12 {
		t.Fatal("configuration fields missing")
	}
	for _, f := range manifest.Config {
		if _, ok := f["default"]; ok {
			t.Errorf("field %s adds a provenance-erasing default", f["key"])
		}
	}
}

func TestManifestSecretContract(t *testing.T) {
	b, err := os.ReadFile("extension.json")
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		DataSecrets *bool                        `json:"data_secrets"`
		Config      []struct{ Key, Type string } `json:"config"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.DataSecrets == nil || !*m.DataSecrets {
		t.Fatal("legacy credential files require data_secrets true")
	}
	found := false
	for _, field := range m.Config {
		if field.Key == "tavily_api_key" {
			found = field.Type == "secret"
		}
	}
	if !found {
		t.Fatal("Tavily key must use a secret field")
	}
}

// Windows accepts root-relative, drive-relative and reserved-device paths that
// IsAbs alone cannot reject. Both preflight and write must reject them.
func TestSavePathRejectsWindowsNonlocalPaths(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("native Windows path semantics")
	}
	cwd := t.TempDir()
	for _, name := range []string{`\rooted.txt`, `C:relative.txt`, `C:\absolute.txt`, `\\server\share\file`, "NUL", "COM1", "file:stream"} {
		if err := checkSavePath(cwd, name, true); err == nil {
			t.Errorf("preflight accepted %q", name)
		}
		if _, err := saveToWorkspace(cwd, name, []byte("x"), true); err == nil {
			t.Errorf("save accepted %q", name)
		}
	}
}

func TestSavePathRejectsGitMetadataAliases(t *testing.T) {
	cwd := t.TempDir()
	for _, name := range []string{".GIT/config", ".Git/hooks/pre-commit", ".git./config", ".git /config"} {
		if err := checkSavePath(cwd, name, true); err == nil {
			t.Errorf("preflight accepted %q", name)
		}
		if _, err := saveToWorkspace(cwd, name, []byte("x"), true); err == nil {
			t.Errorf("save accepted %q", name)
		}
	}
	entries, err := os.ReadDir(cwd)
	if err != nil || len(entries) != 0 {
		t.Fatalf("rejected paths changed workspace: %v, %v", entries, err)
	}
}
