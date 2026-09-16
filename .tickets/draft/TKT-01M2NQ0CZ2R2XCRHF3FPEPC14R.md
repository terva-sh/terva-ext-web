---
schema: 3
id: TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
title: Migrate protocol integration to the verified Terva SDK
type: task
status: draft
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
  - ref: file:internal/proto/proto.go
    path: internal/proto/proto.go
  - ref: file:extension.json
    path: extension.json
  - ref: file:conformance_test.go
    path: conformance_test.go
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

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Completed SDK verification with published version and contracts.

Expected deliverable: SDK-backed registration/session/result transport while retaining existing config behavior.

Validation: Baseline then just CI, vendor sync, race/conformance and five-target snapshot.

Readiness/coordination: main.go overlaps the cwd fix: prefer landing that small fix first, but it is not a semantic dependency; adapt its test if work overlaps.
