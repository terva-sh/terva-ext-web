---
schema: 3
id: TKT-01M2NQ0EPWH96WEK8BDQ2HB6Y7
title: Evaluate checksum-verified prebuilt fallback for source installs
type: spike
status: draft
status_reason: null
priority: low
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:launcher
  - area:release
  - area:security
assignees: []
milestone: null
parent: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
origin: null
dependencies:
  - TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:run.sh
    path: run.sh
  - ref: file:scripts/release.sh
    path: scripts/release.sh
  - ref: file:.goreleaser.yaml
    path: .goreleaser.yaml
  - ref: file:docs/plans/release-process.md
    path: docs/plans/release-process.md
claim: null
archive: null
created_at: 2026-09-16T18:17:34Z
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

Evaluate fetching a compatible prebuilt archive when source installation lacks a Go toolchain. Keep this separate from progress reporting. Use actual published release metadata and explicit integrity, platform and failure policy rather than trusting an unverified download.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Compare source-only and prebuilt-fallback tradeoffs using the published artifact layout
- [ ] Define checksum verification, version/platform matching, offline/error behavior and rollback before any execution of downloads
- [ ] Ensure any proposal preserves research skill/manifest identity and avoids credential/configuration mutation
- [ ] Record adopt/defer/reject rationale and file a separate implementation ticket if selected

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
