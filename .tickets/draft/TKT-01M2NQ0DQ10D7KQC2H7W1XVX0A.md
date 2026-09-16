---
schema: 3
id: TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A
title: Validate replacement installation, upgrades and release platforms
type: task
status: draft
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:release
  - area:validation
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0DANEHS615M5VXNGS7W5
  - TKT-01M2NQ0DEFWGV6VTRQTHSM80CT
  - TKT-01M2NQ0DK05JAWQRF7TBGF1CY7
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:run.sh
    path: run.sh
  - ref: file:.goreleaser.yaml
    path: .goreleaser.yaml
  - ref: file:extension.json
    path: extension.json
  - ref: file:skills/web-research/SKILL.md
    path: skills/web-research/SKILL.md
  - ref: file:docs/plans/release-process.md
    path: docs/plans/release-process.md
  - ref: file:README.md
    path: README.md
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
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

Close the remaining release validation gap: snapshots are built for five targets but only Linux amd64 was executed locally. Validate supported runtime launch, real Terva source/archive installs, existing-settings upgrades and rollback before publication. Use isolated synthetic installations; identify paths from ext list and host data, not repository basenames.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Validate each advertised platform/architecture at runtime or explicitly narrow the release support matrix with rationale, including Windows Bash requirements
- [ ] Verify matching manifest/hello/version, research skill discovery, checksums, offline source builds and no-toolchain archive launch
- [ ] Prove new install, legacy settings upgrade and rollback on supported Terva versions without enabling old and new web identities together
- [ ] Preserve old installation and settings until the replacement smoke test succeeds; use synthetic credentials and document exact reproducible steps
- [ ] Record test evidence for the release candidate and unresolved limitations in release documentation

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
