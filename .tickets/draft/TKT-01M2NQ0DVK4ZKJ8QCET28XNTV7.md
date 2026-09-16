---
schema: 3
id: TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7
title: Publish and smoke-test the first terva-ext-web release
type: task
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:release
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:.goreleaser.yaml
    path: .goreleaser.yaml
  - ref: file:.forgejo/workflows/release.yml
    path: .forgejo/workflows/release.yml
  - ref: file:scripts/release.sh
    path: scripts/release.sh
  - ref: file:release.just
    path: release.just
  - ref: file:docs/plans/release-process.md
    path: docs/plans/release-process.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
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

Publication is currently disabled in GoReleaser and CI. After release validation, prepare the first 0.4.0 release for the verified terva-sh/terva-ext-web destination, rechecking metadata at execution time. Inventory needed release assets, CI secret names, branch protection and mirror settings without exposing credentials or copying inherited secrets. This draft records intended work, not authorization to change remote settings or publish yet.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Record release readiness and user approval before publication; obtain explicit ticket authorization for any required remote settings changes
- [ ] Confirm tag/source version and correct Forgejo destination, and deliberately enable the reviewed release path without reusing cut/* tags or publishing to zot-web
- [ ] Publish approved versioned archives/checksums, then download and verify the published artifacts and installation flow
- [ ] Update installation/release documentation with actual artifact URLs and retain rollback instructions
- [ ] Do not force-push, rotate/revoke credentials or modify the legacy repository as an incidental release step
- [ ] If release-enabling changes alter the validated candidate, rebuild and rerun affected checks on the exact tagged commit before publication; do not reuse stale validation evidence
- [ ] Record whether released artifacts match the tested checksums and retain a recovery plan that does not force-push tags

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: A validated candidate plus actual publication approval and verified destination.

Expected deliverable: Tagged first release and verified downloadable artifacts with recovery instructions.

Validation: Revalidate any changed candidate then verify downloaded checksums and smoke installation.

Readiness/coordination: Automation is still disabled. Reading metadata/secret names is preparatory; remote settings changes still need specific ticket authorization.
