---
schema: 3
id: TKT-01M2NQ0DANEHS615M5VXNGS7W5
title: Validate SDK conformance against supported Terva hosts
type: task
status: draft
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
claim: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T18:31:37Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/critical-path
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

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: SDK integration, cwd regression and authority decision; versioned real hosts from SDK verification.

Expected deliverable: Updated subprocess harness and separate real-host evidence for supported floor/current versions.

Validation: Frame recovery, image/text, concurrency, sessions and bounded clean shutdown with JSON-only stdout.

Readiness/coordination: The existing harness strips ZOT_WEB_* but inherits TAVILY_API_KEY. Fix fixture isolation in this ticket, not by reading real user credentials.
