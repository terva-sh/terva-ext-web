---
schema: 3
id: TKT-01M2NQ0EGSJJMK5CESDRDQV7MM
title: Decide whether private page caches need project isolation
type: spike
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:cache
  - area:sessions
  - area:security
assignees: []
milestone: null
parent: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
origin: null
dependencies:
  - TKT-01M2NQ0DEFWGV6VTRQTHSM80CT
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:internal/fetch/fetch.go
    path: internal/fetch/fetch.go
  - ref: file:main.go
    path: main.go
  - ref: file:internal/config/config.go
    path: internal/config/config.go
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
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

The page cache is process-global. Decide whether session/project transitions should segregate private page snapshots beyond the mandatory configuration invalidation in the config ticket. Use verified project/session identities rather than guessed directories.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Document cache lifecycle and privacy risks across session, project and configuration transitions
- [ ] Evaluate scoping or invalidation alternatives for cached text, raw bytes, images and links, including resource bounds
- [ ] Record the selected policy with scenarios that prove tightened egress cannot be bypassed
- [ ] File separate implementation work if additional isolation is warranted; do not duplicate config invalidation work
- [ ] State whether isolation is keyed by project, session or config generation, and define transitions including no-session and same-project session switches; specify tests without changing product behavior in this spike

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Config invalidation design and verified project/session identity behavior.

Expected deliverable: A policy decision on additional project/session isolation plus implementation criteria if needed.

Validation: Cross-project, same-project/new-session, no-session and config-generation cases within cache limits.

Readiness/coordination: Distinct from mandatory stale-config cache invalidation; do not make this exploratory ticket a hidden release dependency.
