---
schema: 3
id: TKT-01M2NPVS83SMPK45HMJTSQGNAY
title: Translate the remaining modernization plan into tickets
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
  actor: agent:codex/modernization-tickets
  branch: docs/modernization-tickets
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: 80d3ec6222d6f2a1e6643ce544211e1601a68240
  session: null
  claimed_at: 2026-09-16T18:15:12Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:15:01Z
updated_at: 2026-09-16T18:15:12Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-tickets
  name: ""
extensions: {}
---

## Description

Convert outstanding work in docs/plans/terva-ext-web.md into scoped draft tickets with acceptance criteria, source references, and dependencies. Preserve completed cutover work as history and keep optional investigations distinct from release prerequisites.

## Acceptance criteria

- [ ] Outstanding implementation batches and feature decisions are mapped to tickets
- [ ] Draft tickets have criteria, dependencies, and source references
- [ ] Plan index and ticket store pass strict validation

## Implementation plan

Map the remaining plan into a core modernization/release epic and a separate follow-up epic. File discrete SDK, cwd, security, conformance, configuration, credentials, release, and presentation investigations with testable criteria and typed source references. Record actual prerequisite edges without serializing independent work. Keep all implementation work draft and leave implementation plans for the agent who claims it. Add a plan index and validate coverage, graph integrity, generated epic index, and instruction stability.
