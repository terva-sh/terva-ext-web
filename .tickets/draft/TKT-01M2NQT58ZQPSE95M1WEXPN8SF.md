---
schema: 3
id: TKT-01M2NQT58ZQPSE95M1WEXPN8SF
title: Define release test matrix and secure platform test access
type: spike
status: draft
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
claim: null
archive: null
created_at: 2026-09-16T18:31:36Z
updated_at: 2026-09-16T18:31:36Z
created_by:
  id: agent:codex/critical-path
  name: ""
updated_by:
  id: agent:codex/critical-path
  name: ""
extensions: {}
---

## Description

The build matrix advertises linux amd64/arm64, darwin amd64/arm64 and windows amd64. Existing CI runs on a Linux docker runner; only Linux amd64 archive launch is recorded. Resolve runtime test access early, before implementation reaches release validation. Keep all five targets as the current expectation; narrowing support requires a recorded user decision. This spike inventories and specifies test access, not authorization to provision infrastructure or change remote settings.

## Acceptance criteria

- [ ] Map each advertised OS/architecture to a real runtime test environment and operator or record missing access explicitly; cross-compilation alone is not runtime evidence
- [ ] Define which tests can use emulation and which need native host integration, with an explicit support decision for Windows Bash/MSYS launch
- [ ] Inventory required host/toolchain versions and isolated source/archive/legacy fixtures; tie the host floor to the SDK verification outcome once available
- [ ] Specify a reproducible test report recording candidate commit, binary checksums, host version, platform, source/archive mode and pass/fail evidence
- [ ] Inventory release automation, secret names/availability, and permissions without reading or recording secret values; link any missing access/provisioning prerequisite
- [ ] Do not mark the spike complete while required platform access or an explicit support-scope decision is missing

## Definition of done

- [ ] Record reviewed decisions, reusable validation inputs and outstanding external requirements in the ticket; pass strict ticket validation
