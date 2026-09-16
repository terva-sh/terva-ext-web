---
schema: 3
id: TKT-01M2NQ8ETTY1Z41FW82QCX2SM8
title: Label modernization tickets for consistent filtering
type: chore
status: in-progress
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
claim:
  actor: agent:codex/ticket-labels
  branch: docs/modernization-tickets
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: 85b110f482affaa71ba2d58bf91b09ba4a81d1c6
  session: null
  claimed_at: 2026-09-16T18:22:18Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:21:56Z
updated_at: 2026-09-16T18:24:04Z
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

- [ ] All 20 modernization backlog tickets have consistent grouping labels
- [ ] The label vocabulary and useful CLI filters are documented
- [ ] Strict ticket checks and filter validation pass without changing work scope

## Implementation plan

Use namespaced initiative, scope and area labels. Apply one shared modernization initiative and exactly one core/follow-up scope to every backlog ticket, with a small set of relevant area labels. Declare the vocabulary in store configuration and document semantics/filter examples outside generated instructions. Use revision-guarded CLI updates; verify every label filter and unchanged ticket content/relationships against a captured baseline. Add the change to the existing open backlog PR.
