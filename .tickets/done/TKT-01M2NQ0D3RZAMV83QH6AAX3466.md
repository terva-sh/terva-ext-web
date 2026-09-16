---
schema: 3
id: TKT-01M2NQ0D3RZAMV83QH6AAX3466
title: Keep download saves in the workspace captured at call start
type: bug
status: done
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
updated_at: 2026-09-16T19:17:33Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

web_fetch_raw and web_fetch_image read cwd before a network fetch and again at write time. A session switch can validate in one workspace and write in another. Fix this application-level race independently of the SDK rename/migration and carry the regression through SDK adoption.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [x] Capture one workspace identity per save invocation and use it for both preflight and final write
- [x] For both raw and image saves, block the fetch, switch the session cwd, resume it, and prove only the original workspace receives the file
- [x] Keep overwrite, symlink, .git, path traversal and resource-limit protections intact
- [x] Run meaningful session-switch regressions under the race detector
- [x] Make each blocked-fetch test fail against the original double-read implementation; synchronize via channels/barriers rather than timing sleeps, and assert no file or parent directory is created in the new workspace

## Definition of done

- [x] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Implementation plan

Capture cwd once in each download handler before preflight and retain it through save. Add real-subprocess conformance regressions for raw and image downloads: a local HTTP server signals request entry and blocks on a channel; send a new session event and a subsequent command as a processing barrier before releasing the response; verify exact bytes only under the original cwd and no directory in the new cwd. Build the conformance subprocess with the race detector so the regression exercises instrumented production code. First prove the tests fail on unchanged handlers, then apply the narrow fix and run just CI plus tagged race conformance. Keep SDK migration separate.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Current raw/image handlers and existing path guard tests; no SDK decision is needed.

Expected deliverable: One captured workspace identity per save with deterministic blocked-fetch regression coverage.

Validation: Both raw and image tests fail on the old double-read and pass on the fix under race testing.

Readiness/coordination: Startable when selected. Existing main_test.go covers path helpers but not this in-flight session-switch scenario.

**agent:codex/critical-path** at 2026-09-16T18:44:48Z

Promoted to ready at the user's explicit request before merging PR #4. Grooming confirmed no unfinished prerequisite dependencies and concrete entry inputs/acceptance criteria. This ticket is available to claim; no implementation work has started.

**agent:codex/modernization-run** at 2026-09-16T19:15:53Z

Unchanged-source just ci passed with installed mise Go 1.27.1 (Go remains absent from default PATH). The existing protocol owns unexported handler/session state, so a real subprocess test avoids adding test-only protocol setters and survives the later SDK migration. The user authorized autonomous execution of the agreed order; record deviations and stop on an external blocker.

**agent:codex/modernization-run** at 2026-09-16T19:17:32Z

Implemented one cwd capture at entry to each raw/image download handler, used for preflight and final write. Existing save helpers, resource limits, permission declarations and SSRF defenses are unchanged.

Added a real-subprocess regression for both tools with a local HTTP server and channel barriers. The request reaching the server proves preflight occurred in the original workspace; a command response following session_start proves the new session was processed before HTTP was released. Tests compare saved bytes in the old workspace and require that no download directory exists in the new workspace. No sleeps are used for ordering.

The new tests failed against the original handlers for both tools: no output existed under the original workspace and the new workspace contained the download directory. After the fix, both pass under go test -race -tags conformance -run TestConformanceSessionSwitchSaves -count=1. Full just ci passes, including vet/format, race tests, old conformance profiles, new regressions and vendor synchronization. Go 1.27.1 was used from the documented mise installation.

Small scope deviation: pulled forward subprocess environment isolation from the conformance ticket (remove TAVILY_API_KEY and both legacy/new override namespaces) because these network tests must be hermetic. Also enabled -race on the subprocess build: -race on the parent test alone would not instrument the actual extension. Existing CI already installs the required C toolchain. Conformance coverage in the later SDK ticket remains outstanding.

## Summary

Both download handlers retain the cwd captured at handler entry. Deterministic real-subprocess raw/image regressions fail on the original implementation and pass on the fix with race instrumentation. Full just ci passes; all existing path/resource/security defenses remain intact. SDK migration is unchanged.
