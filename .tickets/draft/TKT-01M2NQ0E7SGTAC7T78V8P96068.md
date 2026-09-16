---
schema: 3
id: TKT-01M2NQ0E7SGTAC7T78V8P96068
title: Add web request display hints and cache/backend status
type: task
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:presentation
assignees: []
milestone: null
parent: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
origin: null
dependencies:
  - TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
  - TKT-01M2NQ0DEFWGV6VTRQTHSM80CT
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:main.go
    path: main.go
  - ref: file:conformance_test.go
    path: conformance_test.go
  - ref: file:README.md
    path: README.md
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

Use supported WithDisplay subjects for queries/URLs and an appropriate cache/backend status surface. Keep web-cache text responses for clients without panels or widgets. Display metadata must not be mistaken for model provenance or permissions.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Show useful query/URL subjects and cache/backend state using verified host APIs
- [ ] Retain text command output and capability fallback for clients without rich display
- [ ] Test error/unconfigured states and ensure no credential values appear in display text
- [ ] Keep presentation separate from model tool output, permission decisions and fetched-content trust

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
