---
schema: 3
id: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
title: Improve Terva web tool discovery and presentation
type: epic
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:discovery
  - area:presentation
  - area:launcher
assignees: []
milestone: null
parent: null
origin: null
dependencies: []
blocks_on: children
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:README.md
    path: README.md
claim: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T18:31:38Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/critical-path
  name: ""
extensions: {}
---

## Description

Track recommended next features and optional investigations from the feature assessment. They do not gate the first replacement release. A spike is complete when it records an evidence-based decision and any separately scoped follow-up; filing it does not commit to implementation. Keep deferred session indexing/interception and skipped connector tunneling out of the active backlog.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Each child task has validation evidence or an explicit disposition
- [ ] Exploratory children record adopt/defer/reject decisions and scoped follow-up work
- [ ] Capability fallbacks preserve working text tools and existing security boundaries
- [ ] Keep SDK/config prerequisites explicit and document adopt/defer/reject outcomes for spikes without promoting their implementation follow-ups automatically

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Approved selected child work and its actual SDK/config prerequisites.

Expected deliverable: Incremental product improvements or evidence-backed spike dispositions.

Validation: Per-child criteria and capability fallback without disrupting released core behavior.

Readiness/coordination: All children remain draft. These do not block first replacement publication.
