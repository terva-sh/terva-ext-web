---
schema: 3
id: TKT-01M2NQ0D3RZAMV83QH6AAX3466
title: Keep download saves in the workspace captured at call start
type: bug
status: draft
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:sessions
  - area:downloads
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies: []
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:main.go
    path: main.go
  - ref: file:main_test.go
    path: main_test.go
  - ref: file:internal/proto/proto.go
    path: internal/proto/proto.go
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T18:34:36Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/critical-path
  name: ""
extensions: {}
---

## Description

web_fetch_raw and web_fetch_image read cwd before a network fetch and again at write time. A session switch can validate in one workspace and write in another. Fix this application-level race independently of the SDK rename/migration and carry the regression through SDK adoption.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Capture one workspace identity per save invocation and use it for both preflight and final write
- [ ] For both raw and image saves, block the fetch, switch the session cwd, resume it, and prove only the original workspace receives the file
- [ ] Keep overwrite, symlink, .git, path traversal and resource-limit protections intact
- [ ] Run meaningful session-switch regressions under the race detector
- [ ] Make each blocked-fetch test fail against the original double-read implementation; synchronize via channels/barriers rather than timing sleeps, and assert no file or parent directory is created in the new workspace

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Current raw/image handlers and existing path guard tests; no SDK decision is needed.

Expected deliverable: One captured workspace identity per save with deterministic blocked-fetch regression coverage.

Validation: Both raw and image tests fail on the old double-read and pass on the fix under race testing.

Readiness/coordination: Startable when selected. Existing main_test.go covers path helpers but not this in-flight session-switch scenario.
