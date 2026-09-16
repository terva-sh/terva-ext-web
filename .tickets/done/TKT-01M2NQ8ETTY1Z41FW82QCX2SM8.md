---
schema: 3
id: TKT-01M2NQ8ETTY1Z41FW82QCX2SM8
title: Label modernization tickets for consistent filtering
type: chore
status: done
status_reason: null
priority: normal
due_on: null
labels:
  - area:workflow
assignees: []
milestone: null
parent: null
origin: null
dependencies: []
blocks_on: none
references:
  - ref: file:docs/ticket-labels.md
    path: docs/ticket-labels.md
  - ref: file:.tickets/config.yml
    path: .tickets/config.yml
claim: null
archive: null
created_at: 2026-09-16T18:21:56Z
updated_at: 2026-09-16T18:24:54Z
created_by:
  id: agent:codex/ticket-labels
  name: ""
updated_by:
  id: agent:codex/ticket-labels
  name: ""
extensions: {}
---

## Description

Review and label the modernization backlog with a documented vocabulary, validate CLI and JSON filtering, and keep existing status, priority and dependency metadata unchanged.

## Acceptance criteria

- [x] All 20 modernization backlog tickets have consistent grouping labels
- [x] The label vocabulary and useful CLI filters are documented
- [x] Strict ticket checks and filter validation pass without changing work scope

## Implementation plan

Use namespaced initiative, scope and area labels. Apply one shared modernization initiative and exactly one core/follow-up scope to every backlog ticket, with a small set of relevant area labels. Declare the vocabulary in store configuration and document semantics/filter examples outside generated instructions. Use revision-guarded CLI updates; verify every label filter and unchanged ticket content/relationships against a captured baseline. Add the change to the existing open backlog PR.

## Notes

**agent:codex/ticket-labels** at 2026-09-16T18:24:54Z

Applied a namespaced label vocabulary to all 20 modernization backlog tickets. initiative:terva-modernization selects the complete backlog, scope:core selects 11 tickets, and scope:follow-up selects 9 (counts include epics). Area labels describe direct work topics, not every dependency. The vocabulary is declared in .tickets/config.yml so unknown labels are caught by strict validation. Workflow maintenance has area:workflow and is excluded from the modernization grouping.

Kept status, priority and type as native fields rather than duplicating them in labels. Preserved ticket bodies, acceptance criteria, parent/dependency edges and claims. Updates used the captured ticket revisions. Labels are applied to children explicitly because parent labels are not inherited. The label taxonomy and examples are in docs/ticket-labels.md and linked from AGENTS.md and the backlog index.

Verified each single-label query against expected membership and compared every original ticket field except labels and mutation metadata with its pre-change value. Repeated --label flags use OR: scope:core plus area:security returns 14 tickets, while their intersection is two tasks. Documented Python 3 JSON intersection and grouping export examples because jq is absent here; all documented CLI examples executed successfully. Regenerating instructions leaves AGENTS.md unchanged. Strict ticket checks and whitespace checks pass. The existing open backlog PR #3 carries these additions; no implementation ticket was promoted.

## Summary

Labeled all 20 modernization tickets with initiative, core/follow-up scope and relevant area labels. Declared the vocabulary and documented CLI/JSON grouping, including repeated-label OR semantics. All expected label memberships and documented filters pass; ticket content and relationships are unchanged. Strict validation and generated-instruction stability pass. Changes are included in the existing backlog PR #3.
