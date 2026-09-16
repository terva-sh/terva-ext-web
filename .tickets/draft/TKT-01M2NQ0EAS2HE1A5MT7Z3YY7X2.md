---
schema: 3
id: TKT-01M2NQ0EAS2HE1A5MT7Z3YY7X2
title: Evaluate small extension-authored context guidance
type: spike
status: draft
status_reason: null
priority: low
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:context
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
  - ref: file:skills/web-research/SKILL.md
    path: skills/web-research/SKILL.md
  - ref: file:main.go
    path: main.go
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

Optionally evaluate static context, RefreshContext and compaction hooks for backend availability and paging guidance. The bundled research skill may already suffice. Fetched page bodies remain untrusted tool data and must not be promoted into system context.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Compare the bundled skill with a bounded extension-authored guidance contribution and record adopt/defer/reject reasoning
- [ ] If adopted, capability-gate refresh/compaction behavior and test no unexpected prompt growth from frequent cache/status changes
- [ ] Keep all fetched page bodies out of system-context contributions and document the trust boundary

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
