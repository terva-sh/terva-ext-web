// Package version holds the extension's version string, shared by the
// protocol hello, the default User-Agent, and the --version flag.
package version

// Version is the terva-ext-web release version. It is a var (not a const) so
// release builds can override it with the git tag via
//
//	-ldflags "-X terva-ext-web/internal/version.Version=..."
//
// (see .goreleaser.yaml). Source builds (run.sh, just build) report this
// committed default — bump it (together with extension.json's "version",
// which a test pins equal) when cutting a release.
var Version = "0.4.0"
