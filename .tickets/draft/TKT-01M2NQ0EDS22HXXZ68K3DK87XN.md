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
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
updated_at: 2026-09-16T18:22:56Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/ticket-labels
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

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
