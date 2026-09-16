---
schema: 3
id: TKT-01M2NYSRTKARHFGRF6XFEMM8A6
title: Fix native Windows save paths and cache eviction ordering
type: bug
status: done
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:security
  - area:cache
  - area:validation
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies: []
blocks_on: none
references: []
claim: null
archive: null
created_at: 2026-09-16T20:33:43Z
updated_at: 2026-09-16T20:51:31Z
created_by:
  id: agent:codex/modernization-run
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Native rehearsal https://github.com/terva-sh/terva-ext-web/actions/runs/35146845707 exposed Windows root-relative save paths accepted by filepath.IsAbs and nondeterministic timestamp-based LRU eviction. Fix without weakening safety tests or adding sleeps. Review case-insensitive Git metadata aliases alongside native path handling.

## Acceptance criteria

- [x] Reject nonlocal save paths before filesystem changes and protect Git metadata aliases
- [x] Preserve deterministic LRU order independently of clock resolution
- [x] Pass native Windows source and archive validation and Linux regression checks

## Implementation plan

Use filepath.IsLocal for native lexical containment and reject case-insensitive/trailing-dot Git metadata aliases. Replace time-based LRU rank with a lock-protected access sequence while retaining timestamps for TTL and display. Preserve existing behavior tests and add path alias regressions. Adding sleeps would hide a production ordering bug, so it is rejected. Rerun native validation on the exact changed candidate.

## Notes

**agent:codex/modernization-run** at 2026-09-16T20:35:09Z

Implemented native locality checks plus Git metadata spelling protection and lock-ordered cache LRU sequence. Existing entry/byte eviction tests exercise the production bug; added preflight/write Windows path and Git alias regressions. Local just ci passed, including race, subprocess conformance, published host driver and vendor checks. Native rerun remains required.

**agent:codex/modernization-run** at 2026-09-16T20:39:19Z

Second rehearsal 35147435097 passed the original Windows path and LRU tests but exposed an incorrect new test expectation: COM1.txt is a regular filename on modern Windows. Verified Go 1.27 internal/filepathlite/path_windows.go uses RtlIsDosDeviceName_U for names with extensions because Windows 11 permits them. Corrected the fixture to bare COM1, which remains reserved. Keep native IsLocal semantics rather than inventing a broader filename restriction.

## Summary

Fixed native save locality/Git metadata aliases and deterministic LRU ordering. Local just ci and all five source/archive native targets pass at 7fde307 in run 35148398553. Reports and rationale: docs/validation/native-2026-09-16.json and docs/plans/platform-validation.md. PR #16 carries the changes.
