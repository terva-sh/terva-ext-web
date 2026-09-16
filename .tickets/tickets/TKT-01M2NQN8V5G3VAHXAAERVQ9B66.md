---
schema: 3
id: TKT-01M2NQN8V5G3VAHXAAERVQ9B66
title: Groom the modernization critical path and readiness gaps
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
references: []
claim:
  actor: agent:codex/critical-path
  branch: docs/critical-path-grooming
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: f49bfdb6fa2d17a94803f57be71fbfd23c1268a8
  session: null
  claimed_at: 2026-09-16T18:29:18Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:28:56Z
updated_at: 2026-09-16T18:29:18Z
created_by:
  id: agent:codex/critical-path
  name: ""
updated_by:
  id: agent:codex/critical-path
  name: ""
extensions: {}
---

## Description

Review the merged backlog against current source and documented Terva contracts, identify release prerequisite paths and parallel work, close ticket-definition gaps, and record unanswered decisions without starting implementation.

## Acceptance criteria

- [ ] Release dependency paths and independent starting work are documented
- [ ] Every modernization ticket has reviewed entry inputs, outputs and validation
- [ ] Missing prerequisite decisions have scoped tickets and dependency edges
- [ ] Strict ticket validation and dependency/link checks pass

## Implementation plan

Review current source, Terva contracts and all draft ticket bodies. Model release dependencies without invented durations. Add only concrete missing readiness decisions, resolve ambiguous completion criteria, and record per-ticket entry inputs, deliverables and validation as grooming notes rather than implementation plans. Keep future work draft; document recommended first selections and external prerequisites. Refresh plan/label counts and validate graph integrity, coverage, references and generated indexes.
