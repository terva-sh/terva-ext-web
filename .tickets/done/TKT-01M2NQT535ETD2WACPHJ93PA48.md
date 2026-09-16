---
schema: 3
id: TKT-01M2NQT535ETD2WACPHJ93PA48
title: Resolve configuration provenance and legacy import semantics
type: spike
status: done
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:config
  - area:security
  - area:sdk
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies:
  - TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K
blocks_on: none
references:
  - ref: file:internal/config/config.go
    path: internal/config/config.go
  - ref: file:internal/config/config_test.go
    path: internal/config/config_test.go
  - ref: file:extension.json
    path: extension.json
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:docs/plans/modernization-critical-path.md
    path: docs/plans/modernization-critical-path.md
claim: null
archive: null
created_at: 2026-09-16T18:31:36Z
updated_at: 2026-09-16T19:34:37Z
created_by:
  id: agent:codex/critical-path
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

The inspected Terva SDK exposes Config as resolved raw JSON values: manifest defaults overlaid with user values. Config.Has tests presence, not explicit user choice. The planned precedence cannot be implemented by treating presence as an explicit host value. Resolve that contract against the selected published SDK/host before nonsecret or secret migration. This decision can proceed alongside SDK integration after version verification; it must not require deploying new host behavior without a linked prerequisite.

## Acceptance criteria

- [x] Record how the selected published host distinguishes explicit settings from defaults, or prove that the distinction is unavailable
- [x] Choose and document a deterministic import/precedence policy using only available APIs; compare missing, explicitly default-valued, zero/false/empty and invalid values
- [x] Define a field-by-field precedence table including allowlist replace-versus-append behavior, new/legacy env variables and TAVILY_API_KEY
- [x] Specify opt-in/idempotent legacy import and rollback behavior without silently rewriting credentials or inventing host provenance
- [x] Provide synthetic migration fixtures/cases for config and secret implementation; record any host change as an explicit blocker and linked ticket

## Definition of done

- [x] Record reviewed decisions, reusable validation inputs and outstanding external requirements in the ticket; pass strict ticket validation

## Implementation plan

Inspect published v0.137.0 resolver, manifest schema and SDK Config API against existing config loader. Establish a schema strategy that makes explicit values distinguishable without invented metadata; specify legacy transition, env precedence, validation and synthetic fixtures. Account separately for omitted undecryptable secrets and opt-in host-only operation; do not read or write real credentials.

## Notes

**agent:codex/critical-path** at 2026-09-16T18:33:36Z

Entry inputs: verified published SDK/host capability report and the existing config loader behavior. Deliverable: a deterministic import/precedence decision and synthetic fixture table for both nonsecret and secret migration. Validation: cover absent, explicit-default, zero/false/empty, legacy/file/env, allowlist append/replace and retry/rollback cases. This decision runs alongside SDK integration and gates config implementation; local SDK presence alone is not evidence of explicit setting provenance.

**agent:codex/modernization-run** at 2026-09-16T19:33:41Z

Published resolver proves Config.Has has no provenance when manifest defaults exist. Decision: omit schema defaults and keep application defaults, use explicit configuration_source host mode for opt-in credential migration and disabling legacy reads. Host drops undecryptable secrets, so legacy mode deliberately retains legacy credential ownership until opt-in; host mode never resurrects a legacy key. Host form blank means unset for nonsecrets and keep for secrets. docs/plans/config-provenance.md records full field precedence, allowlist exception, validation, immutable runtime/cache semantics, rollback and synthetic cases. Manual host-form import wins over automatic file/credential rewriting; no new host API needed.

## Summary

Recorded available resolver/form/SDK contract and deterministic no-manifest-default policy in docs/plans/config-provenance.md. Explicit host mode opts into credential import and disables legacy reads. Field table, env and allowlist semantics, immutable updates, synthetic cases and rollback require no new host API. Config and secret tickets consume this decision; production remains legacy until those batches.
