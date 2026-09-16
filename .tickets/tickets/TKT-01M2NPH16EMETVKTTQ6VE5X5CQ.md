---
schema: 3
id: TKT-01M2NPH16EMETVKTTQ6VE5X5CQ
title: Integrate git ticket into the agent workflow
type: chore
status: in-progress
status_reason: null
priority: normal
due_on: null
labels: []
assignees: []
milestone: null
parent: null
origin: null
dependencies: []
blocks_on: none
references: []
claim:
  actor: agent:codex/ticket-workflow
  branch: chore/ticket-workflow
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: 80a2b86c274a7999f18eb7b14a45f4fe649896f2
  session: null
  claimed_at: 2026-09-16T18:09:22Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:09:08Z
updated_at: 2026-09-16T18:09:22Z
created_by:
  id: agent:codex/ticket-workflow
  name: ""
updated_by:
  id: agent:codex/ticket-workflow
  name: ""
extensions: {}
---

## Description

Complete the user-requested repository ticket-store setup after merging PR #1. Preserve generated instructions and add repository-specific work tracking and validation guidance.

## Acceptance criteria

- [ ] Repository ticket store is initialized and committed
- [ ] AGENTS.md includes generated and repository-specific ticket guidance
- [ ] Ticket validation and instruction regeneration checks pass

## Implementation plan

Keep the generated workflow block tool-owned. Replace obsolete no-store guidance with repository scope, plan references, and handoff checks; add a just recipe for strict read-only store validation and document merge-driver setup for each clone. Verify regeneration preserves local guidance and run ticket validation.
