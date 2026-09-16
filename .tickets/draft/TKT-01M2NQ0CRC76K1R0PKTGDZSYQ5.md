---
schema: 3
id: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
title: Complete Terva modernization and replacement release
type: epic
status: draft
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:sdk
  - area:config
  - area:release
assignees: []
milestone: null
parent: null
origin: null
dependencies: []
blocks_on: children
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:docs/plans/identity-packaging.md
    path: docs/plans/identity-packaging.md
  - ref: file:docs/plans/release-process.md
    path: docs/plans/release-process.md
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

Track the remaining core batches from the split plan: SDK and correctness, configuration and credentials, then validated release and legacy transition. Identity/packaging and repository creation are already complete in PR #1; ticket workflow is complete in PR #2. Do not redo the fork, rename, or baseline. Feature experiments live in a separate epic. Creating this epic does not authorize publication, credential changes, or archival.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Core child tickets are completed with recorded validation evidence
- [ ] Supported-host installation, upgrade and rollback are verified before directing users to migrate
- [ ] Legacy transition is resolved with explicit authorization for any remote archival
- [ ] Track the publication milestone separately from subsequent legacy notice/archival completion; optional follow-up features are not hidden release gates

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Completed cutover/store work and the groomed core dependency graph.

Expected deliverable: A verified replacement release, followed by resolved legacy transition.

Validation: Child-ticket evidence; publication milestone is separate from epic-wide completion.

Readiness/coordination: No reliable duration estimates yet, so critical path means release prerequisite chains, not a schedule or delivery date.
