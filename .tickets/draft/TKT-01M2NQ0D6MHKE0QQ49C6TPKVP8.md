---
schema: 3
id: TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8
title: Preserve network and write authority during SDK adoption
type: task
status: draft
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:sdk
  - area:security
  - area:downloads
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:main.go
    path: main.go
  - ref: file:extension.json
    path: extension.json
  - ref: file:main_test.go
    path: main_test.go
  - ref: file:internal/fetch/fetch.go
    path: internal/fetch/fetch.go
  - ref: file:docs/plans/untrusted-web-content.md
    path: docs/plans/untrusted-web-content.md
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

The wire carries one authority string, while download tools have both network and workspace effects. Preserve network gating and ask defaults; do not relabel writers as workspace-mutation alone and lose egress gating. Review selective SDK ordering for stateful commands/writes without serializing network reads. Keep untrusted content protections independent of SDK result metadata.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Document a verified authority/permission mapping for all six tools and combined write/network effects
- [ ] Show writer ask defaults and network gating remain effective on supported host modes and user overrides remain host-controlled
- [ ] Record whether selective Sequential use is needed, with concurrency evidence
- [ ] Retain and run SSRF/redirect, path, byte/image/cache limit and output-sanitization regressions; do not claim an SDK trust field that does not exist
- [ ] Record host-mode and user-override test cases for both writers and all readers, including overwrite refusal and denied-network behavior; preserve application SSRF defenses even when host permission allows the tool

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Migrated SDK registrations and verified host authority semantics.

Expected deliverable: An explicit reader/writer permission matrix and ordering decision.

Validation: Host modes/user overrides plus SSRF, write-path and concurrency regression tests.

Readiness/coordination: Do not conflate single wire authority with combined effects; preserve existing network-read plus ask behavior unless a proven alternative is reviewed.
