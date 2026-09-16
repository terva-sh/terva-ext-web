---
schema: 3
id: TKT-01M2NQN8V5G3VAHXAAERVQ9B66
title: Groom the modernization critical path and readiness gaps
type: chore
status: done
status_reason: null
priority: normal
due_on: null
labels:
  - area:workflow
assignees: []
milestone: null
parent: null
origin: null
dependencies: []
blocks_on: none
references:
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim: null
archive: null
created_at: 2026-09-16T18:28:56Z
updated_at: 2026-09-16T18:34:37Z
created_by:
  id: agent:codex/critical-path
  name: ""
updated_by:
  id: agent:codex/critical-path
  name: ""
extensions: {}
---

## Description

Review the merged backlog against current source and documented Terva contracts, identify release prerequisite paths and parallel work, close ticket-definition gaps, and record unanswered decisions without starting implementation.

## Acceptance criteria

- [x] Release dependency paths and independent starting work are documented
- [x] Every modernization ticket has reviewed entry inputs, outputs and validation
- [x] Missing prerequisite decisions have scoped tickets and dependency edges
- [x] Strict ticket validation and dependency/link checks pass

## Implementation plan

Review current source, Terva contracts and all draft ticket bodies. Model release dependencies without invented durations. Add only concrete missing readiness decisions, resolve ambiguous completion criteria, and record per-ticket entry inputs, deliverables and validation as grooming notes rather than implementation plans. Keep future work draft; document recommended first selections and external prerequisites. Refresh plan/label counts and validate graph integrity, coverage, references and generated indexes.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:34:37Z

PR #3 merged as f49bfdb after successful push/PR lint checks. Reviewed all 20 existing modernization tickets and current code plus Terva docs/SDK at local revision 7f754b9bb7c6754284dc4f3cc5fa6a96525de49d; the unrelated Terva working-tree ticket remained untouched.

Release prerequisites form two joins after published SDK verification: SDK -> authority -> conformance, and SDK plus a config import decision -> configuration -> credentials. The independent cwd fix joins at conformance; platform access joins at release validation. Neither a follow-up feature nor legacy archival is a prerequisite of publication. There are no estimates or known platform capacity sufficient to claim a duration-based critical path.

Added two scoped core spikes for concrete missing inputs: resolved-config provenance/import semantics and platform/runtime test access. Config.Has only tests presence of resolved defaults/user values; treating it as explicit user choice would break migration precedence. The CI/build setup demonstrates cross-builds, not native runtime access for all five targets. Configuration now depends on the provenance decision; release validation depends on the platform matrix. Platform preparation can begin before SDK selection because final host-version slots are filled by SDK verification later.

Each original ticket has grooming notes for entry inputs, deliverable, validation and coordination, plus targeted observable criteria. Added tests/decisions for inherited TAVILY_API_KEY in subprocess environments, stale config fetches repopulating cache, session-boundary search withdrawal, unknown bootstrap-frame rejection, candidate-specific release validation, safe raw/binary host-write semantics and checksum authenticity. No implementation plans were prefilled, no draft was promoted, and no SDK/publication/credential action was performed.

Updated the backlog index and label counts for 22 drafts (13 core and 9 follow-up, including epics), and added docs/plans/modernization-critical-path.md with the graph, first selections, observed evidence and unresolved inputs. Strict store checks pass. Programmatic validation confirmed all tickets have grooming notes and unchecked criteria, no claims/plans, valid source/document references, acyclic dependencies, exactly three independent release roots and no follow-up/legacy publication gates. Generated instructions are unchanged. No application code changed, so Go tests were not rerun.

## Summary

Merged PR #3 and groomed all modernization work. Added configuration-import and runtime-platform decision gates, tightened existing criteria, and documented the release prerequisite graph with independent SDK verification, cwd repair and platform preparation starting points. All 22 implementation/epic tickets remain draft; external SDK availability and platform access are still unverified inputs owned by explicit tickets. Strict validation, dependency/reference checks and instruction regeneration pass.
