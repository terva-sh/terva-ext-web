---
schema: 3
id: TKT-01M2NQ0DEFWGV6VTRQTHSM80CT
title: Add validated host configuration and legacy precedence
type: task
status: in-progress
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:config
  - area:cache
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
  - TKT-01M2NQT535ETD2WACPHJ93PA48
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:extension.json
    path: extension.json
  - ref: file:internal/config/config.go
    path: internal/config/config.go
  - ref: file:internal/config/config_test.go
    path: internal/config/config_test.go
  - ref: file:main.go
    path: main.go
  - ref: file:internal/fetch/fetch.go
    path: internal/fetch/fetch.go
  - ref: file:README.md
    path: README.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
  - ref: file:docs/plans/config-provenance.md
    path: docs/plans/config-provenance.md
claim:
  actor: agent:codex/modernization-run
  branch: feat/host-config
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: b36f628ceaf55d41827fb4977d9217f62317b246
  session: null
  claimed_at: 2026-09-16T19:42:48Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T19:46:28Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Move backend, URL, limits, User-Agent and allowlist settings to manifest/SDK configuration with validated immutable runtime snapshots and atomic updates. During a bounded transition preserve the proposed precedence: new TERVA_EXT_WEB_* overrides, legacy ZOT_WEB_* overrides, explicitly configured host values, legacy file values, defaults. Host-resolved defaults are not proof of explicit configuration. Secret handling is a separate dependent ticket.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Document and test precedence, preserving existing legacy choices and TAVILY_API_KEY behavior without introducing unsafe secret defaults
- [ ] Expose supported nonsecret settings with validation and distinguish explicit host values from resolved defaults or record a tested alternative import mechanism
- [ ] Apply config updates atomically; in-flight operations keep a coherent snapshot and rejected updates remain visible while working settings survive
- [ ] Define and test cache invalidation for backend, allowlist and User-Agent changes so stale entries cannot bypass tightened policy
- [ ] Cover first setup, malformed settings, missing backend credentials, legacy files and concurrent updates with synthetic fixtures only
- [ ] Consume the approved configuration-provenance decision and fixtures before implementing host/legacy precedence
- [ ] Test a blocked old-config fetch finishing after a backend/allowlist/User-Agent change; it must not repopulate or become visible through the new-config cache
- [ ] Define and test whether a tightening update affects already accepted requests versus later requests; do not claim cancellation of accepted calls without a supported mechanism

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Implementation plan

Implement config.Resolve with application defaults, absent-default manifest fields, explicit host/legacy mode, new/old environment precedence and sanitized validation. Retain Load for historical caller compatibility while production uses validated resolution. Add a mutex-protected immutable runtime snapshot per call and replace fetcher/provider/cache on valid updates; rejected updates retain prior runtime and notify without values. Test full precedence, initial errors for every network handler, concurrent updates and blocked old fetch cache isolation. Host Tavily secret field follows in its separate ticket.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: SDK integration plus the new provenance/import-contract decision and fixtures.

Expected deliverable: Immutable validated runtime configuration and explicit cache/update semantics.

Validation: Full precedence table, invalid-update retention, blocked-request updates and stale cache completion tests.

Readiness/coordination: Config.Has cannot distinguish explicit settings from default values in the inspected SDK; do not infer provenance from presence.

**agent:codex/modernization-run** at 2026-09-16T19:46:28Z

Implemented strict host/legacy/env resolution and immutable per-call runtimes with fresh clients/caches on updates. Unit/race tests pass for explicit defaults, zero/false/empty, namespace precedence, safe invalid settings, legacy failure/host bypass, rejected-update retention and old-fetch cache isolation after backend/allowlist/UA changes. Production now gates all six handlers on initial config errors. Host secret field remains a separate batch. Pausing edits briefly to resolve conformance CI failure before stacking more work.
