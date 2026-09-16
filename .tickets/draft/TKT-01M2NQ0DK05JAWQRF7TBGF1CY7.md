---
schema: 3
id: TKT-01M2NQ0DK05JAWQRF7TBGF1CY7
title: Migrate Tavily configuration secrets without losing legacy settings
type: task
status: draft
status_reason: null
priority: high
due_on: null
labels: []
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
claim: null
archive: null
created_at: 2026-09-16T18:17:33Z
updated_at: 2026-09-16T18:17:33Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/modernization-tickets
  name: ""
extensions: {}
---

## Description

Use a host secret-config field for a user-provided Tavily key; the runtime secret broker is for credentials acquired at runtime and is not automatically needed. Preserve TAVILY_API_KEY. Define an explicit, recoverable legacy import rather than delete/rewrite user credential files as a side effect. A masked field is not proof of encryption; data_secrets must describe the actual directory.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Test synthetic legacy secret import, host/env precedence and rollback without logging or committing secret values
- [ ] Preserve credential files unless the user explicitly authorizes a concrete migration action; document duplicate/legacy secret disposition
- [ ] Do not set data_secrets=false while a legacy credential file can remain; verify the declaration against resulting data
- [ ] Document host encryption prerequisites and config-secret versus runtime-broker contracts; only require protocol 6 if broker use is justified
- [ ] Verify errors, migration notes, tool output and logs do not disclose synthetic secret values

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
