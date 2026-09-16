# Network and workspace authority

All six tools retain `network-read`, with `read_only` absent/false. The wire
has one authority string, not combined effect flags. Raw saves and optional
image saves additionally carry manifest `ask` rules. Relabeling either writer
as workspace-mutation would misdescribe its egress; dropping its ask rule
would weaken the explicit write declaration. No SDK trust metadata exists.

| Tools | Effects | Manifest rule |
| --- | --- | --- |
| web_search, web_fetch, web_images, web_links | Network read, transient cache | Host mode/user policy |
| web_fetch_raw | Network read and guarded workspace write | ask |
| web_fetch_image | Network read, optional guarded workspace write | ask (also without save_path) |

Published Terva v0.137.0 host source: `build.ExtToolReadOnly` uses declared
authority before the legacy boolean; `ExtToolRegisters` excludes these tools
from plan mode. `permissions.BuildPolicy` orders user, project, then extension
suggestions. `core.PermissionPolicy.Evaluate` denies nonlocal tools in plan
before rule evaluation. Other modes honor explicit user allow/deny first.

| Mode | Default for all six | Explicit user allow | Explicit user deny |
| --- | --- | --- | --- |
| plan | deny / not registered | deny | deny |
| ask | ask | allow | deny |
| auto-edit | ask | allow | deny |
| workspace | ask (foreign network tools) | allow | deny |
| yolo | allow, including manifest ask | allow | deny |

Ask is a host prompt request, not an unconditional prohibition: yolo overrides
ask, and headless modes refuse when a prompt is required. Do not promise that
the extension can override user-selected host policy. Network authority is a
classification for approval, not a destination firewall. Application SSRF,
redirect, allowlist and resource guards apply even when host permission allows.
User rules are host-owned; this extension never rewrites them.

`tests/host-contract` runs the real published host policy with the actual
manifest, temporary host config and all five modes plus user overrides.
`conformance_test.go` inspects the actual six SDK registrations and runs two
blocked raw downloads concurrently, proving neither is serialized behind the
other. Session-switch regressions cover both writers. Existing fetch/path tests
cover denied destinations, redirects, byte/image/cache limits, sanitization,
symlink escapes and overwrite refusal.

Do not apply Sequential to downloads: their preflight and write guards are
per-call, writes use exclusive creation unless overwrite is explicitly true,
and concurrent network reads are useful. Cache operations already synchronize
internally. /web-cache clear is a best-effort cache eviction, not cancellation
of an active fetch; adding serial command dispatch would not make it one.
A separate conformance batch adds full actual-host launch and frame recovery.

Validation recorded 2026-09-16: all 15 mode/override policy combinations and
uncached race/conformance application suites passed locally; Forgejo PR #10
passed and merged. This is published host-code policy evidence, not a UI test.
