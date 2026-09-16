---
schema: 3
id: TKT-01M2NWQSBVHEYHFW8C1FD5ZTVX
title: Bound tool results to the supported host frame limit
type: bug
status: in-progress
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:sdk
  - area:validation
  - area:security
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies: []
blocks_on: none
references: []
claim:
  actor: agent:codex/modernization-run
  branch: fix/result-frame-budget
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: cec0228a573372c5b5629085564031858cf13572
  session: null
  claimed_at: 2026-09-16T19:58:11Z
  expires_at: null
archive: null
created_at: 2026-09-16T19:57:41Z
updated_at: 2026-09-16T20:00:48Z
created_by:
  id: agent:codex/modernization-run
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Release review found a published SDK/host limit mismatch: extproto reads at most 4 MiB per frame but Encode does not bound output. Existing image configuration permits 5 MiB encoded image bytes before base64, so valid application results can be silently discarded by the host and time out. Large text/cache output also needs a wire budget. Preserve download limits and return actionable bounded errors instead of oversized frames.

## Acceptance criteria

- [ ] Reproduce an image/result exceeding the 4 MiB wire budget and return a bounded actionable result; keep save-only behavior available
- [ ] Bound serialized tool content and cache command output including JSON escaping/base64 overhead and reserve host envelope space
- [ ] Keep normal text/image output unchanged and run race/conformance plus actual host-driver checks

## Definition of done

- [ ] Record deviation and validation; merge separately before release validation

## Implementation plan

Bound marshaled content (including JSON/base64 expansion) below the published 4 MiB host ceiling with 4 KiB envelope reserve for generated host IDs. Wrap all registered tool results and cache command display output; reject oversized image injection before writes with resize/save-only guidance. Preserve configured fetch/download byte limits. Verify escaped text, large images and normal results, including actual host-driver large-image round-trip and save-only fallback.

## Notes

**agent:codex/modernization-run** at 2026-09-16T20:00:48Z

Verified the mismatch with a deterministic PNG over 3 MiB but below the application 5 MiB limit. Actual published host-driver test now receives an actionable bounded error for injection (no file written) and saves identical bytes with inject:false. Added serialized text/base64/escaping budget tests and bounded all tool results plus cache command responses with 4 KiB envelope reserve. Normal output unchanged. This is a release-review deviation; separate from merged secret PR #13.
