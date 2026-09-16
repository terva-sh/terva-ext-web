---
schema: 3
id: TKT-01M2NTZ8KJ9CXBV1RTMZVXRP77
title: Wire GitHub tag-time native release validation
type: task
status: done
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
  - TKT-01M2NYSRTKARHFGRF6XFEMM8A6
  - TKT-01M2NZD0DP9MTFSJ62F8BZV96V
blocks_on: none
references: []
claim: null
archive: null
created_at: 2026-09-16T19:26:49Z
updated_at: 2026-09-16T20:51:48Z
created_by:
  id: agent:codex/modernization-run
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Implement docs/plans/platform-validation.md in the user-provided release mirror. Verify the exact GitHub destination and transport access without guessing or changing remote settings. Follow sibling GitHub-only tag/manual workflows; native tests precede asset publication. Linux amd64 remains the development gate.

## Acceptance criteria

- [x] Verify mirror URL and access from documented configuration or user input; record no credentials
- [x] Native tag/manual jobs exercise all five target archives and actual Terva launch, including Windows Git Bash, with candidate-bound reports
- [x] Keep validation read-only and publication dependent on passing jobs; manual runs never publish

## Definition of done

- [x] Record workflow validation and native run evidence or explicit external blockers

## Implementation plan

Add GitHub-only version-tag/manual validation workflow following sibling conventions. Build five GoReleaser archives once on Linux, then verify checksums and execute source/archive conformance plus actual host-driver launch on native Linux amd64/arm64, macOS amd64/arm64 and Windows amd64 runners. Upload candidate-bound reports with least-privilege permissions and no publication. Verify the empty user-supplied mirror, push reviewed source after Forgejo gates, and use manual rehearsal before any release tag. Report unavailable runner or host launch as an explicit blocker.

## Notes

**agent:codex/modernization-run** at 2026-09-16T20:01:49Z

User supplied https://github.com/terva-sh/terva-ext-web and git@github.com:terva-sh/terva-ext-web.git. Verified GitHub API identifies a public empty repository and authenticated push access; git ls-remote succeeds with no refs. This worktree initially had only origin, so added the supplied URL as local mirror. No remote repository settings changed and no credential contents read.

**agent:codex/modernization-run** at 2026-09-16T20:05:13Z

Implemented GitHub-only tag/manual native workflow with read-only permissions, Linux packaging and five native runner targets. Added candidate checksum/version/skill/no-Go launcher validation and reusable source/archive host-driver checks with JSON reports. Local Linux archive rehearsal passes, including installation path with spaces; native GitHub execution remains pending. Added .exe suffix for Windows subprocess test builds. Manual workflow never publishes; tag builds also skip publication until the separate approved release path.

**agent:codex/modernization-run** at 2026-09-16T20:29:38Z

Merged Forgejo PRs #14 and #15 after passing CI. Pushed reviewed commit 87ebe9610f0053bbe5caeaaba1b2c54e07fcfebb to the previously empty GitHub mirror main using an ordinary push, without tags or settings changes. Dispatched read-only manual rehearsal https://github.com/terva-sh/terva-ext-web/actions/runs/35146845707; native outcomes pending. Manual snapshot was chosen to test runner availability and archives without publishing a release.

## Summary

Verified and populated the user-supplied GitHub mirror without changing remote settings. Native manual run https://github.com/terva-sh/terva-ext-web/actions/runs/35148398553 passed packaging plus source and archive validation on Linux amd64/arm64, macOS amd64/arm64 and Windows amd64 at candidate 7fde307f42609844d87881101a13c14a4d53dcc3. All ten reports, including archive/binary hashes and runner metadata, are committed in docs/validation/native-2026-09-16.json. Native failures were repaired in linked tickets, with rationale in docs/plans/platform-validation.md. Manual and tag validation never publish; the final versioned release must rerun validation. Full CLI installation/upgrade/rollback remains the next release-validation ticket.
