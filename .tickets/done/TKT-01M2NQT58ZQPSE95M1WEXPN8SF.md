---
schema: 3
id: TKT-01M2NQT58ZQPSE95M1WEXPN8SF
title: Define release test matrix and secure platform test access
type: spike
status: done
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
claim: null
archive: null
created_at: 2026-09-16T18:31:36Z
updated_at: 2026-09-16T19:27:03Z
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

- [x] Map each advertised OS/architecture to a real runtime test environment and operator or record missing access explicitly; cross-compilation alone is not runtime evidence
- [x] Define which tests can use emulation and which need native host integration, with an explicit support decision for Windows Bash/MSYS launch
- [x] Specify a reproducible test report recording candidate commit, binary checksums, host version, platform, source/archive mode and pass/fail evidence
- [x] Inventory release automation, secret names/availability, and permissions without reading or recording secret values; link any missing access/provisioning prerequisite
- [x] Do not mark the spike complete while required platform access or an explicit support-scope decision is missing
- [x] Define toolchain requirements and isolated source/archive/legacy fixture slots; final floor/current Terva versions come from SDK verification at release validation, so access preparation can proceed independently

## Definition of done

- [x] Record reviewed decisions, reusable validation inputs and outstanding external requirements in the ticket; pass strict ticket validation

## Implementation plan

Inventory documented build targets, available local executables and Forgejo runtime/runner metadata using read-only filtered output. Record a five-target runtime matrix and a candidate-bound test report template with source/archive/legacy fixture slots. Keep Go 1.27 and published Terva v0.137.0 as the verified version inputs. Identify actual missing access without inventing hosts or treating cross-builds as runtime tests; ask the user for the missing platform environment or an explicit support decision and record an external prerequisite if unavailable.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:33:36Z

Entry inputs: current five-target build matrix, Linux docker CI and documented Linux amd64 archive evidence. Deliverable: platform access/support matrix and reusable candidate-report template. Validation: identify executable target environments, fixture/runner ownership, and missing access; do not provision infrastructure or reduce platform support as an implicit grooming action. Start access preparation independently; final host versions are supplied by SDK verification before release validation.

**agent:codex/critical-path** at 2026-09-16T18:44:48Z

Promoted to ready at the user's explicit request before merging PR #4. Grooming confirmed no unfinished prerequisite dependencies and concrete entry inputs/acceptance criteria. This ticket is available to claim; no implementation work has started.

**agent:codex/modernization-run** at 2026-09-16T19:26:49Z

User resolved access policy: Linux amd64 is primary; release mirror GitHub runners test other platforms preferably on tags. docs/plans/platform-validation.md records native matrix, tag/manual rehearsal gates, report and fixture requirements, sibling references and missing mirror verification. This explicit scope/timing decision unblocks development; no non-Linux runtime pass is claimed. Go 1.27 is approved. Existing Forgejo secret-name inventory is empty; no values read.

## Summary

Defined platform policy and reproducible report in docs/plans/platform-validation.md. User selected Linux amd64 development and GitHub tag-time cross-platform testing. TKT-01M2NTZ8KJ9CXBV1RTMZVXRP77 (Wire GitHub tag-time native release validation) carries mirror verification and execution prerequisites; release validation depends on it. No runtime evidence claimed for other platforms.
