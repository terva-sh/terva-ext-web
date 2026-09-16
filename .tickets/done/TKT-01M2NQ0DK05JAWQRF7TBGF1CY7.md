---
schema: 3
id: TKT-01M2NQ0DK05JAWQRF7TBGF1CY7
title: Migrate Tavily configuration secrets without losing legacy settings
type: task
status: done
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:config
  - area:security
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0DEFWGV6VTRQTHSM80CT
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:extension.json
    path: extension.json
  - ref: file:internal/config/config.go
    path: internal/config/config.go
  - ref: file:internal/config/config_test.go
    path: internal/config/config_test.go
  - ref: file:main.go
    path: main.go
  - ref: file:README.md
    path: README.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
  - ref: file:docs/plans/config-provenance.md
    path: docs/plans/config-provenance.md
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

Use a host secret-config field for a user-provided Tavily key; the runtime secret broker is for credentials acquired at runtime and is not automatically needed. Preserve TAVILY_API_KEY. Define an explicit, recoverable legacy import rather than delete/rewrite user credential files as a side effect. A masked field is not proof of encryption; data_secrets must describe the actual directory.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [x] Test synthetic legacy secret import, host/env precedence and rollback without logging or committing secret values
- [x] Preserve credential files unless the user explicitly authorizes a concrete migration action; document duplicate/legacy secret disposition
- [x] Do not set data_secrets=false while a legacy credential file can remain; verify the declaration against resulting data
- [x] Document host encryption prerequisites and config-secret versus runtime-broker contracts; only require protocol 6 if broker use is justified
- [x] Verify errors, migration notes, tool output and logs do not disclose synthetic secret values
- [x] Use only synthetic credentials when testing idempotent imports, duplicate legacy/host values, unset/reset behavior and rollback; assert fixtures remain recoverable and secret values never reach captured logs

## Definition of done

- [x] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Implementation plan

Expose a default-free host secret field only in explicit host mode; retain provider env precedence and legacy files for rollback. Declare data_secrets true because files may remain. Exercise synthetic duplicate/import/retry/unset/rollback and error-output checks. Remove raw Tavily error body snippets and redact echoed configured key from provider results/errors; keep status guidance. Document host encryption setup separately from secret-field masking; do not use runtime broker.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: Working config migration using approved import semantics.

Expected deliverable: Secret-field migration, truthful data_secrets policy, recoverable legacy state and documented encryption prerequisites.

Validation: Synthetic import/retry/unset/rollback cases and no-secret-output assertions.

Readiness/coordination: No runtime broker needed merely to read a user-supplied Tavily key; no credential rotation or source-file deletion is implied.

**agent:codex/modernization-run** at 2026-09-16T19:55:31Z

Added default-free Tavily secret config field gated by explicit host mode, provider env precedence and truthful data_secrets=true. Legacy files are unchanged and rollback remains possible. Synthetic precedence/import/retry/unset/rollback tests pass. Provider-error snippets could echo credentials, so removed them while retaining HTTP guidance; configured-key echoes are redacted from result fields and printable errors with error identity preserved. Added subprocess stdout/stderr secret-diagnostic tests and manifest invariants. No real credentials read or changed. README documents optional host encryption and broker distinction.

**agent:codex/modernization-run** at 2026-09-16T19:58:10Z

Full local just ci and uncached race/conformance suite pass, including actual published host driver and synthetic credential stdout/stderr/echo checks. PR #13 carries secret migration; remote checks are pending. Follow-up ddgr idea filed separately at user request. Release review also filed TKT-01M2NWQSBVHEYHFW8C1FD5ZTVX (Bound tool results to the supported host frame limit) as a prerequisite; it does not change credential scope.

## Summary

Host secret field, opt-in migration, env precedence, missing-key behavior and unchanged legacy files verified with synthetic tests. data_secrets=true and documented host encryption/broker boundaries. Provider errors and credential echoes are sanitized; local CI passes. PR #13 awaits remote gate and merge.
