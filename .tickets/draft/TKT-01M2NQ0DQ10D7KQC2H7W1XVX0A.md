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
  - TKT-01M2NQT58ZQPSE95M1WEXPN8SF
  - TKT-01M2NTZ8KJ9CXBV1RTMZVXRP77
  - TKT-01M2NWQSBVHEYHFW8C1FD5ZTVX
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
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
updated_at: 2026-09-16T19:58:10Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-run
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
- [ ] Use the approved platform/access matrix and bind the complete validation report to one candidate commit and checksums; identify revalidation triggers after source or packaging changes
- [ ] Test installation discovery and rollback with preserved synthetic old settings and verify the two web identities are never enabled together

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Conformance/config/secrets completion and the new release matrix/access decision.

Expected deliverable: Candidate-bound runtime, source install, archive install, upgrade and rollback evidence.

Validation: Run on the agreed advertised platform matrix and supported host versions; verify actual manifests, skills and checksums.

Readiness/coordination: Only Linux amd64 launch is established; the matrix spike exposes missing access early. Missing access is not a passing test.
