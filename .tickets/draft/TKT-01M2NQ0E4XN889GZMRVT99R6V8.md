---
schema: 3
id: TKT-01M2NQ0E4XN889GZMRVT99R6V8
title: Measure whether web tools need eager discovery hints
type: spike
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:discovery
assignees: []
milestone: null
parent: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
origin: null
dependencies:
  - TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:main.go
    path: main.go
  - ref: file:skills/web-research/SKILL.md
    path: skills/web-research/SKILL.md
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

Assess Essential and host lazy-tool discovery before pinning anything. Existing discovery may be sufficient; do not mark all six tools essential by default. Compare search/fetch visibility and prompt cost against current behavior.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Record measurements or reproducible observations of default discovery and Essential hints under the host cap
- [ ] Recommend a minimal hint set or retain defaults, explaining the tradeoff
- [ ] Verify any adopted hints degrade on supported hosts and do not unexpectedly grow prompt context

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
