---
schema: 3
id: TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K
title: Verify a published Terva SDK and supported host floor
type: spike
status: in-progress
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:sdk
  - area:validation
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies: []
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:go.mod
    path: go.mod
  - ref: file:internal/proto/proto.go
    path: internal/proto/proto.go
  - ref: file:docs/plans/terva-host-features-2026-06.md
    path: docs/plans/terva-host-features-2026-06.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim:
  actor: agent:codex/modernization-run
  branch: spike/sdk-verification
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: 25ec03fbcf32fbfca545b07f32ff822416d2e986
  session: null
  claimed_at: 2026-09-16T19:18:01Z
  expires_at: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T19:18:01Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

The handwritten protocol is still in use. Verify a published module release before selecting the SDK pin; inspect v0.137.0 first because the recorded local tag has the needed APIs, but do not treat local tags or sibling template pins as proof of public availability. Reinspect sibling status/revision without changing it. Establish the minimum supported Terva version and capability policy from actual contracts.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Record a retrievable published module version, provenance, required Go version and API evidence
- [ ] Document minimum Terva/protocol support: ordered sessions require protocol 2; host tools require 3; runtime broker use would require 6
- [ ] Document available events, display/visibility fallbacks, and actual limits on per-call cancellation and trust metadata
- [ ] Record any release/vendor/toolchain compatibility constraints for the migration
- [ ] Deliver a versioned capability matrix and commands that reproduce module retrieval from a clean environment; include module checksum/provenance, supported Go/host versions and launcher bootstrap compatibility

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Implementation plan

Verify v0.137.0 through the public Go module proxy first using an isolated module cache and explicit proxy/sumdb settings, without editing repository dependencies. If unavailable, inspect the authoritative public repository/module metadata for other published versions and document a viable alternative or an actual publication blocker. For a verified version, inspect SDK/host APIs, Go requirements, version/capability negotiation, bootstrap behavior and checksums; record reproducible commands and an honest supported-host floor. Do not substitute a local tag or replace directive for a published release.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:37Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Current go.mod, recorded local Terva revision and the SDK requirements in the plan.

Expected deliverable: A retrievable published SDK pin and capability/host-floor matrix.

Validation: Clean module retrieval, API inspection and compatible toolchain evidence.

Readiness/coordination: Public module availability and exact oldest supported host are still unverified; a local tag is not sufficient.

**agent:codex/critical-path** at 2026-09-16T18:44:48Z

Promoted to ready at the user's explicit request before merging PR #4. Grooming confirmed no unfinished prerequisite dependencies and concrete entry inputs/acceptance criteria. This ticket is available to claim; no implementation work has started.
