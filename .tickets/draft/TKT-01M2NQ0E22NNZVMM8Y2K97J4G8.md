---
schema: 3
id: TKT-01M2NQ0E22NNZVMM8Y2K97J4G8
title: Withdraw and restore search when its backend is unavailable
type: task
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:discovery
  - area:config
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
  - ref: file:main.go
    path: main.go
  - ref: file:extension.json
    path: extension.json
  - ref: file:conformance_test.go
    path: conformance_test.go
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

Use supported v4 withdrawal/restore at appropriate session boundaries to hide only unavailable search. Backend configuration changes can race calls, so UI visibility cannot replace runtime readiness checks. Fetch tools should remain usable without search credentials.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Withdraw only web_search when unavailable on hosts supporting the capability, then restore it when usable
- [ ] Handle session boundaries and config transitions without stale visibility or removing working fetch tools
- [ ] Keep runtime readiness errors safe for in-flight calls and retain usable fallback behavior on older supported hosts
- [ ] Test absent credentials, backend recovery, config updates and unsupported withdrawal capability
- [ ] A config update changes runtime readiness immediately but withdrawal/restore is asserted at the next supported session boundary; test calls during that interval without mutating the prompt prefix off-boundary

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Config snapshot/readiness behavior and SDK session-scoped withdrawal.

Expected deliverable: Predictable available-search visibility with runtime errors between boundaries.

Validation: Absent/invalid credentials, backend recovery, update races and unsupported hosts.

Readiness/coordination: The inspected host guidance forbids changing cached-prefix visibility per config event; apply it at the next session boundary.
