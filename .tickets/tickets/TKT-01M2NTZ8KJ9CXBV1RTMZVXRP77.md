---
schema: 3
id: TKT-01M2NTZ8KJ9CXBV1RTMZVXRP77
title: Wire GitHub tag-time native release validation
type: task
status: in-progress
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
dependencies: []
blocks_on: none
references: []
claim:
  actor: agent:codex/modernization-run
  branch: ci/github-native-validation
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: 6550a0528eb415e5da5b79f163fec95d197fd960
  session: null
  claimed_at: 2026-09-16T20:01:49Z
  expires_at: null
archive: null
created_at: 2026-09-16T19:26:49Z
updated_at: 2026-09-16T20:01:49Z
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

- [ ] Verify mirror URL and access from documented configuration or user input; record no credentials
- [ ] Native tag/manual jobs exercise all five target archives and actual Terva launch, including Windows Git Bash, with candidate-bound reports
- [ ] Keep validation read-only and publication dependent on passing jobs; manual runs never publish

## Definition of done

- [ ] Record workflow validation and native run evidence or explicit external blockers

## Implementation plan

Add GitHub-only version-tag/manual validation workflow following sibling conventions. Build five GoReleaser archives once on Linux, then verify checksums and execute source/archive conformance plus actual host-driver launch on native Linux amd64/arm64, macOS amd64/arm64 and Windows amd64 runners. Upload candidate-bound reports with least-privilege permissions and no publication. Verify the empty user-supplied mirror, push reviewed source after Forgejo gates, and use manual rehearsal before any release tag. Report unavailable runner or host launch as an explicit blocker.

## Notes

**agent:codex/modernization-run** at 2026-09-16T20:01:49Z

User supplied https://github.com/terva-sh/terva-ext-web and git@github.com:terva-sh/terva-ext-web.git. Verified GitHub API identifies a public empty repository and authenticated push access; git ls-remote succeeds with no refs. This worktree initially had only origin, so added the supplied URL as local mirror. No remote repository settings changed and no credential contents read.
