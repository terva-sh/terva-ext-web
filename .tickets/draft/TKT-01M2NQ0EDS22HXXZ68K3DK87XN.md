---
schema: 3
id: TKT-01M2NQ0EDS22HXXZ68K3DK87XN
title: Assess host-brokered saves without weakening download semantics
type: spike
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:downloads
  - area:security
assignees: []
milestone: null
parent: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
origin: null
dependencies:
  - TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
  - TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:main.go
    path: main.go
  - ref: file:main_test.go
    path: main_test.go
  - ref: file:docs/plans/untrusted-web-content.md
    path: docs/plans/untrusted-web-content.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
updated_at: 2026-09-16T18:34:36Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/critical-path
  name: ""
extensions: {}
---

## Description

Investigate v3 HostToolCall for text saves so host policy can participate. Do not assume host text-write tools can preserve arbitrary downloaded bytes or binary images. Keep direct guarded saves unless overwrite, path, size, raw-byte and session semantics can be matched.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Compare actual host write contracts against overwrite, symlink/path, byte limits, session identity and raw-byte requirements
- [ ] Assess network and workspace permission effects together and capability/version fallbacks
- [ ] Record an evidence-based adopt/defer/reject decision separately for raw text and binary images
- [ ] If adoption is viable, file scoped implementation/tests rather than weakening existing guards in this investigation
- [ ] Include non-UTF8 raw responses and binary image fixtures when comparing host text-write semantics; retaining direct guarded saves is an acceptable documented spike outcome

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: SDK HostToolCall and authority decision.

Expected deliverable: Byte/path/permission comparison for text and binary saves with a feasibility conclusion.

Validation: Synthetic raw/non-UTF8/image, overwrite and escape cases against actual host write contracts.

Readiness/coordination: Do not start a migration just because the API exists; file implementation only if semantic equivalence is demonstrated.
