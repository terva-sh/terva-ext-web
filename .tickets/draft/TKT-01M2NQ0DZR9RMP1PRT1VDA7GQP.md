---
schema: 3
id: TKT-01M2NQ0DZR9RMP1PRT1VDA7GQP
title: Coordinate legacy migration notice and authorized zot-web archival
type: task
status: draft
status_reason: null
priority: normal
due_on: null
labels: []
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:docs/plans/release-process.md
    path: docs/plans/release-process.md
  - ref: file:README.md
    path: README.md
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
updated_at: 2026-09-16T18:17:33Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-tickets
  name: ""
extensions: {}
---

## Description

After the replacement release is verified, coordinate a migration notice in the source zot-web repository. That repository owns its changes: inspect its instructions/status and create or export a linked task there as required. Archival is a remote settings change and requires a person to record explicit authorization in the implementing ticket. This draft is a tracking record, not that authorization.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Confirm published replacement installation and settings-preserving rollback before recommending migration
- [ ] Prepare and review the legacy migration notice with verified destination links and the shared web-identity warning
- [ ] Record the source-repository task or PR and leave unrelated source work untouched
- [ ] Archive only after explicit ticket authorization is recorded; otherwise keep that action outstanding and do not claim archival complete

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
