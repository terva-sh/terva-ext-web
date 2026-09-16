---
schema: 3
id: TKT-01M2NPH16EMETVKTTQ6VE5X5CQ
title: Integrate git ticket into the agent workflow
type: chore
status: done
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
references:
  - ref: file:AGENTS.md
    path: AGENTS.md
  - ref: file:justfile
    path: justfile
claim: null
archive: null
created_at: 2026-09-16T18:09:08Z
updated_at: 2026-09-16T18:10:41Z
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

- [x] Repository ticket store is initialized and committed
- [x] AGENTS.md includes generated and repository-specific ticket guidance
- [x] Ticket validation and instruction regeneration checks pass

## Implementation plan

Keep the generated workflow block tool-owned. Replace obsolete no-store guidance with repository scope, plan references, and handoff checks; add a just recipe for strict read-only store validation and document merge-driver setup for each clone. Verify regeneration preserves local guidance and run ticket validation.

## Notes

**agent:codex/ticket-workflow** at 2026-09-16T18:09:52Z

PR #1 merged as 2ac463c after successful remote lint/test checks. Store initialization is committed in 80a2b86. Repository additions stay outside the generated block so regeneration is safe. Added a local strict just recipe rather than extending Go CI: CI does not provision git-ticket, and selecting a distribution/pin is separate work. Retained default store configuration instead of inventing label or actor allowlists. Installed the merge driver for this clone; future clones must install it explicitly.

## Summary

Initialized the project store and generated AGENTS.md workflow, documented project scope and plan references, added strict just ticket-check validation, and installed/documented the clone-local merge driver. Instruction regeneration is idempotent; strict read-only validation and git diff --check pass. Code and Go dependencies are unchanged. Setup is complete on chore/ticket-workflow for separate review after the identity cutover merge.
