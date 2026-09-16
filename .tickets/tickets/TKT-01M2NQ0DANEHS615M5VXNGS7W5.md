---
schema: 3
id: TKT-01M2NQ0DANEHS615M5VXNGS7W5
title: Validate SDK conformance against supported Terva hosts
type: task
status: in-progress
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:sdk
  - area:sessions
  - area:validation
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
  - TKT-01M2NQ0D3RZAMV83QH6AAX3466
  - TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:conformance_test.go
    path: conformance_test.go
  - ref: file:main_test.go
    path: main_test.go
  - ref: file:justfile
    path: justfile
  - ref: file:.forgejo/workflows/ci.yml
    path: .forgejo/workflows/ci.yml
  - ref: file:README.md
    path: README.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim:
  actor: agent:codex/modernization-run
  branch: test/sdk-conformance
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: df228a854be9b4d3811a2293819aeb02a8c717e1
  session: null
  claimed_at: 2026-09-16T19:39:18Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T19:41:58Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Replace the inherited zot/early-Terva-only conformance contract with subprocess checks for the documented supported floor and current Terva. Distinguish a simulated host profile from evidence obtained using real supported host versions. Retire stock-zot-only profiles after documenting the new support contract.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Exercise hello, all six tools, web-cache, text/image results and clean shutdown over real subprocess stdio
- [ ] Verify malformed/oversized-frame recovery, concurrent calls, ordered sessions and blocked-fetch session-switch saves
- [ ] Test both documented host floor and current Terva; record versions and real-host smoke evidence
- [ ] Keep stdout JSON-only and preserve race/vet/format/vendor gates in just and CI
- [ ] Document capability fallback and support boundaries before removing zot-only profiles
- [ ] Remove inherited TAVILY_API_KEY, ZOT_WEB_* and TERVA_EXT_WEB_* from subprocess environments; use synthetic fixtures and local servers so tests cannot call real providers
- [ ] Use image fixtures and bounded timeout/deadlock checks for oversized/malformed frames and concurrent shutdown; record actual-host tests separately from simulated wire profiles

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Implementation plan

Extend subprocess conformance with successful local-fixture calls for all six tools including image content, malformed/oversized recovery and blocked-call shutdown. Preserve deterministic concurrency/session cases. Add actual published v0.137.0 host-driver launch to isolated host-contract module, distinguish this from simulated wire and full CLI install smoke. Floor/current coincide at verified v0.137.0. Run race/CI gates, record exact limits, and carry full installation scenarios to release validation.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: SDK integration, cwd regression and authority decision; versioned real hosts from SDK verification.

Expected deliverable: Updated subprocess harness and separate real-host evidence for supported floor/current versions.

Validation: Frame recovery, image/text, concurrency, sessions and bounded clean shutdown with JSON-only stdout.

Readiness/coordination: The existing harness strips ZOT_WEB_* but inherits TAVILY_API_KEY. Fix fixture isolation in this ticket, not by reading real user credentials.

**agent:codex/modernization-run** at 2026-09-16T19:17:32Z

The cwd-save fix pulled forward subprocess environment isolation (drop TAVILY_API_KEY, ZOT_WEB_* and TERVA_EXT_WEB_*) and enabled race instrumentation for the real extension subprocess. Preserve these and the deterministic raw/image session-switch regressions when migrating the harness to SDK/current-host conformance. This does not complete the broader conformance ticket.

**agent:codex/modernization-run** at 2026-09-16T19:41:58Z

Expanded real subprocess suite passes under race: all six tools with successful synthetic local responses and exact image bytes, malformed JSON and >4 MiB frame recovery, concurrent downloads, both session-switch saves and shutdown during blocked fetch. Actual published v0.137.0 extdriver launch also passes via isolated host-contract module: launcher, ready, six tools, command, ordered session/save, stdout monitor and shutdown. Floor/current coincide at verified v0.137.0. docs/plans/sdk-conformance.md distinguishes host-driver evidence from full CLI installation/upgrade smoke, which remains a release gate. No per-call cancellation claim.
