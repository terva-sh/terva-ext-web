# SDK conformance evidence

Supported floor and current baseline both refer to published Terva v0.137.0
(protocol 6); SDK requires min_protocol 2 for ordered session identity. Older
Terva and stock zot are not validated support targets. Optional capabilities
(display, withdrawal, host tool calls, runtime secret broker) are not required
or enabled by this migration. No optional event is assumed from version alone.

Two distinct automated layers now run in `just ci` and Forgejo:

- `conformance_test.go` builds the real extension with the race detector and
  simulates the supported wire contract. It checks hello/identity, exactly six
  network-read tools, command/subscription, local-fixture search/fetch/images/
  links/raw/image results including exact decoded image bytes, JSON-only
  stdout and shutdown. It sends malformed JSON and a frame above the SDK's
  4 MiB inbound limit and checks subsequent command recovery with a deadline.
  Blocked-download cases prove concurrent raw calls, workspace capture across
  session changes for both writers, and bounded shutdown without a save.
- `tests/host-contract/driver_test.go` imports the actual published host
  `extdriver`, launches our copied Bash launcher and race-instrumented binary
  from a temporary install, and tests ready, six registrations, command,
  session event, network save, shutdown and the host's malformed-frame monitor.
  This exercises the authoritative host wire implementation, not a fabricated
  hello_ack. The separate module records the host pin and dependency checksums.

Both pass on Linux amd64 with Go 1.27.1 (2026-09-16). The earlier authority
matrix also runs the real host permission resolver against our actual manifest.
Tests strip web/provider overrides, use temporary homes/data/workspaces and
local HTTP fixtures, and do not read installed config or call real providers.

Limits: this is not a full interactive Terva CLI or installed-release test.
Discovery/trust/UI setup, real source installation, upgrade and rollback remain
release-validation checks. Other OS/architectures have cross-build evidence
only until GitHub tag-time native validation. SDK shutdown exits the process;
it is not per-call cancellation, and active network handlers are not promised
to finish after shutdown. The test intentionally asserts no completed save
while the server remains blocked, not cancellation metadata the SDK lacks.

All legacy protocol-package tests were retired during SDK adoption; application
SSRF, redirect, resource/cache bounds, path and sanitization regressions remain.
