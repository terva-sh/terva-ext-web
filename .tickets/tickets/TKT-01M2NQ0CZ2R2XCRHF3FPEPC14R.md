---
schema: 3
id: TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
title: Migrate protocol integration to the verified Terva SDK
type: task
status: in-progress
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:sdk
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:go.mod
    path: go.mod
  - ref: file:go.sum
    path: go.sum
  - ref: file:vendor/modules.txt
    path: vendor/modules.txt
  - ref: file:main.go
    path: main.go
  - ref: file:extension.json
    path: extension.json
  - ref: file:conformance_test.go
    path: conformance_test.go
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
  - ref: file:docs/plans/sdk-verification.md
    path: docs/plans/sdk-verification.md
claim:
  actor: agent:codex/modernization-run
  branch: feat/terva-sdk
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: 94ef28e81bbdcea20be708e5f3b22fd25ef96e98
  session: null
  claimed_at: 2026-09-16T19:27:25Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T19:31:03Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Replace handwritten internal/proto framing, registration and dispatch with the verified shared ext/extproto SDK in a reviewable batch. Preserve fetch/search behavior, public identities, lifecycle handling, text/image results, and offline vendor builds. Configuration schema migration and credentials remain separate tickets.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Pin and vendor the verified published SDK without a local replace directive
- [ ] Preserve web hello/manifest identity, six web_* tools, web-cache, permissions and JSON-only stdout
- [ ] Use ordered session lifecycle and advertised event capabilities while preserving concurrent network reads
- [ ] Retire the handwritten protocol integration and adapt existing behavior tests without weakening SSRF, limits, paths or sanitization
- [ ] Run applicable just CI, race, vendor-sync and snapshot checks; document unsupported cancellation/trust features honestly
- [ ] Carry the session-switch save regression through SDK migration if it has landed; coordinate overlapping main.go edits and preserve legacy config until its own migration batch
- [ ] Update typed ticket references to retired protocol files when deleting them so git ticket check stays valid
- [ ] Upgrade go.mod, source-build documentation/launcher requirements and CI to Go 1.27+ together; verify the CI image and inspect combined transitive/vendor changes before accepting the published SDK pin

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Implementation plan

Baseline just ci passes before SDK edits. Replace main.go registration, results and lifecycle with published ext APIs; retain one Host().CWD snapshot per save and lazy legacy configuration. Move registration inspection to subprocess wire assertions because SDK has no Tools accessor. Retire handwritten package, adopt documented v0.137.0 protocol-6 harness with min_protocol 2, and preserve all application tests. Upgrade Go requirement/CI/launcher together, inspect module graph/vendor changes, run race/conformance/vendor/snapshot gates and verify Go 1.27 CI through the real runner. Keep actual-host and expanded conformance coverage in their following ticket.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Completed SDK verification with published version and contracts.

Expected deliverable: SDK-backed registration/session/result transport while retaining existing config behavior.

Validation: Baseline then just CI, vendor sync, race/conformance and five-target snapshot.

Readiness/coordination: main.go overlaps the cwd fix: prefer landing that small fix first, but it is not a semantic dependency; adapt its test if work overlaps.

**agent:codex/modernization-run** at 2026-09-16T19:25:50Z

User explicitly approved upgrading to Go 1.27 on 2026-09-16. Update module, CI and source-launch requirements together in the SDK batch.

**agent:codex/modernization-run** at 2026-09-16T19:29:49Z

Migrated to ext SDK with min_protocol 2, OnSession subscription and Host().CWD snapshot. Published v0.137.0 harness replaces obsolete zot/early-Terva profiles; actual host tests remain next batch. Registration authority checks moved from inaccessible in-process registry to real SDK wire (exact six tools, network-read, not read_only). Legacy configuration and manifest ask rules unchanged. Combined module selection upgrades x/image 0.42→0.44, x/net 0.55→0.56 and x/text 0.38→0.40, and introduces age/hpke/crypto/sys through SDK imports. Inspected resulting decoder/parser/normalization and vendor package changes; go mod verify and full uncached race+conformance suite pass. Direct CI registry probe requires authentication (401); actual CI execution will verify the approved Go 1.27 image. Removed retired protocol package and repointed typed references to main.go.

**agent:codex/modernization-run** at 2026-09-16T19:31:03Z

Post-removal CI caught a third handwritten-protocol test file left behind (command_test.go). Removed it and preserved application command/display/error coverage in a real subprocess conformance test; startup checks now also assert web-cache registration and command capability. This supersedes the earlier complete-removal statement until the rerun passes.
