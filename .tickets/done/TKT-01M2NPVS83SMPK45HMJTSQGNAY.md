---
schema: 3
id: TKT-01M2NPVS83SMPK45HMJTSQGNAY
title: Translate the remaining modernization plan into tickets
type: chore
status: done
status_reason: null
priority: normal
due_on: null
labels: []
assignees: []
milestone: null
parent: null
origin: null
dependencies: []
blocks_on: none
references:
  - ref: file:docs/plans/modernization-backlog.md
    path: docs/plans/modernization-backlog.md
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
claim: null
archive: null
created_at: 2026-09-16T18:15:01Z
updated_at: 2026-09-16T18:19:50Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-tickets
  name: ""
extensions: {}
---

## Description

Convert outstanding work in docs/plans/terva-ext-web.md into scoped draft tickets with acceptance criteria, source references, and dependencies. Preserve completed cutover work as history and keep optional investigations distinct from release prerequisites.

## Acceptance criteria

- [x] Outstanding implementation batches and feature decisions are mapped to tickets
- [x] Draft tickets have criteria, dependencies, and source references
- [x] Plan index and ticket store pass strict validation

## Implementation plan

Map the remaining plan into a core modernization/release epic and a separate follow-up epic. File discrete SDK, cwd, security, conformance, configuration, credentials, release, and presentation investigations with testable criteria and typed source references. Record actual prerequisite edges without serializing independent work. Keep all implementation work draft and leave implementation plans for the agent who claims it. Add a plan index and validate coverage, graph integrity, generated epic index, and instruction stability.

## Notes

**agent:codex/modernization-tickets** at 2026-09-16T18:19:50Z

Translated the remaining split plan into 20 draft tickets: two epics and 18 scoped tasks/investigations. The core epic covers verified SDK selection, SDK integration, the independent cwd save race, authority/concurrency, host conformance, configuration, credentials, release validation/publication and legacy transition. The follow-up epic covers visibility, discovery, display, optional context, brokered writes, project cache isolation, startup progress and prebuilt fallback.

Separating follow-ups avoids making optional experiments prerequisites for the replacement release. SDK verification and the cwd bug have no dependencies and can be selected independently; configuration follows SDK integration without waiting for all conformance work. Release validation joins conformance and configuration/secret migration. No implementation ticket was promoted or given a speculative implementation plan.

Completed split, baseline, identity/packaging, remote setup and ticket-store work remain linked historical evidence rather than duplicate outstanding tickets. Session indexing and interception remain deferred, connector tunneling stays skipped, and unsupported cancellation/trust metadata are limitations rather than feature commitments. Selective ordering is covered in authority/concurrency. No new authorization for publication, credentials, remote settings or archival was inferred from filing tickets.

Added docs/plans/modernization-backlog.md as the coverage/dependency index; tickets own current status. Corrected stale plan and agent statements about baseline, remote and archive packaging completion. The initial strict check correctly detected a stale generated epic index; git ticket check --fix regenerated it. Final strict checks pass. Programmatic review confirmed all 20 tickets are unclaimed drafts with expected parent/dependency edges, real source references and unchecked acceptance criteria; implementation plans remain empty. Document links resolve and instruction regeneration is byte-for-byte stable. No code changed, so Go tests were not repeated.

## Summary

PR #2 merged as 80d3ec6. Filed 20 draft backlog tickets with acceptance criteria, source references and prerequisite edges, grouped into core modernization/release and follow-up epics. Added a linked coverage index and refreshed completed-work evidence. Strict ticket validation, graph/criteria/reference checks, local document links, instruction regeneration and diff checks pass. SDK verification and the cwd save fix are independent candidates for user-selected promotion; all future implementation work remains draft.
