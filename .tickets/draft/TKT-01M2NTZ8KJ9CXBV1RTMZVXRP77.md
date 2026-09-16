---
schema: 3
id: TKT-01M2NTZ8KJ9CXBV1RTMZVXRP77
title: Wire GitHub tag-time native release validation
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
dependencies: []
blocks_on: none
references: []
claim: null
archive: null
created_at: 2026-09-16T19:26:49Z
updated_at: 2026-09-16T19:26:49Z
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
