---
schema: 3
id: TKT-01M2NQT58ZQPSE95M1WEXPN8SF
title: Define release test matrix and secure platform test access
type: spike
status: in-progress
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:release
  - area:validation
  - area:launcher
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies: []
blocks_on: none
references:
  - ref: file:.goreleaser.yaml
    path: .goreleaser.yaml
  - ref: file:.forgejo/workflows/ci.yml
    path: .forgejo/workflows/ci.yml
  - ref: file:.forgejo/workflows/release.yml
    path: .forgejo/workflows/release.yml
  - ref: file:docs/plans/release-process.md
    path: docs/plans/release-process.md
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim:
  actor: agent:codex/modernization-run
  branch: spike/release-test-access
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: 99f68ba27385b730c3cff28080407aa18f6ad660
  session: null
  claimed_at: 2026-09-16T19:23:13Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:31:36Z
updated_at: 2026-09-16T19:23:13Z
created_by:
  id: agent:codex/critical-path
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

The build matrix advertises linux amd64/arm64, darwin amd64/arm64 and windows amd64. Existing CI runs on a Linux docker runner; only Linux amd64 archive launch is recorded. Resolve runtime test access early, before implementation reaches release validation. Keep all five targets as the current expectation; narrowing support requires a recorded user decision. This spike inventories and specifies test access, not authorization to provision infrastructure or change remote settings.

## Acceptance criteria

- [ ] Map each advertised OS/architecture to a real runtime test environment and operator or record missing access explicitly; cross-compilation alone is not runtime evidence
- [ ] Define which tests can use emulation and which need native host integration, with an explicit support decision for Windows Bash/MSYS launch
- [ ] Specify a reproducible test report recording candidate commit, binary checksums, host version, platform, source/archive mode and pass/fail evidence
- [ ] Inventory release automation, secret names/availability, and permissions without reading or recording secret values; link any missing access/provisioning prerequisite
- [ ] Do not mark the spike complete while required platform access or an explicit support-scope decision is missing
- [ ] Define toolchain requirements and isolated source/archive/legacy fixture slots; final floor/current Terva versions come from SDK verification at release validation, so access preparation can proceed independently

## Definition of done

- [ ] Record reviewed decisions, reusable validation inputs and outstanding external requirements in the ticket; pass strict ticket validation

## Implementation plan

Inventory documented build targets, available local executables and Forgejo runtime/runner metadata using read-only filtered output. Record a five-target runtime matrix and a candidate-bound test report template with source/archive/legacy fixture slots. Keep Go 1.27 and published Terva v0.137.0 as the verified version inputs. Identify actual missing access without inventing hosts or treating cross-builds as runtime tests; ask the user for the missing platform environment or an explicit support decision and record an external prerequisite if unavailable.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:33:36Z

Entry inputs: current five-target build matrix, Linux docker CI and documented Linux amd64 archive evidence. Deliverable: platform access/support matrix and reusable candidate-report template. Validation: identify executable target environments, fixture/runner ownership, and missing access; do not provision infrastructure or reduce platform support as an implicit grooming action. Start access preparation independently; final host versions are supplied by SDK verification before release validation.

**agent:codex/critical-path** at 2026-09-16T18:44:48Z

Promoted to ready at the user's explicit request before merging PR #4. Grooming confirmed no unfinished prerequisite dependencies and concrete entry inputs/acceptance criteria. This ticket is available to claim; no implementation work has started.
