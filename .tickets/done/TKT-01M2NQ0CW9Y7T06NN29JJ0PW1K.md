---
schema: 3
id: TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K
title: Verify a published Terva SDK and supported host floor
type: spike
status: done
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
  - ref: file:docs/plans/terva-host-features-2026-06.md
    path: docs/plans/terva-host-features-2026-06.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
  - ref: file:docs/plans/sdk-verification.md
    path: docs/plans/sdk-verification.md
  - ref: file:main.go
    path: main.go
claim: null
archive: null
created_at: 2026-09-16T18:17:32Z
updated_at: 2026-09-16T19:29:33Z
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

- [x] Record a retrievable published module version, provenance, required Go version and API evidence
- [x] Document minimum Terva/protocol support: ordered sessions require protocol 2; host tools require 3; runtime broker use would require 6
- [x] Document available events, display/visibility fallbacks, and actual limits on per-call cancellation and trust metadata
- [x] Record any release/vendor/toolchain compatibility constraints for the migration
- [x] Deliver a versioned capability matrix and commands that reproduce module retrieval from a clean environment; include module checksum/provenance, supported Go/host versions and launcher bootstrap compatibility

## Definition of done

- [x] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

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

**agent:codex/modernization-run** at 2026-09-16T19:22:29Z

Verified public v0.137.0 from an empty cache through proxy.golang.org and sum.golang.org, including provenance/checksums. Isolated vendor probe builds and speaks hello/min_protocol/registration/shutdown on Linux; offline cross-builds pass for all five advertised targets. Selected v0.137.0 as initial supported host/test baseline, with protocol 2 correctness floor and optional capability gates. Real-host conformance remains its own ticket. Public module requires Go 1.27.0: source launcher/README/CI must upgrade coherently in the migration. Full evidence and reproducible retrieval commands are in docs/plans/sdk-verification.md.

## Summary

Selected verified published terva.sh/terva v0.137.0 with checksums and origin evidence. Recorded capability matrix, v0.137.0 supported-host baseline, protocol 2 correctness floor and honest cancellation/trust limits. Isolated module verification, vendoring, Linux wire probe and offline five-target SDK builds pass. Migration must upgrade Go 1.25 requirements/CI to 1.27 and rerun combined dependency checks; no repository dependency change was made.
