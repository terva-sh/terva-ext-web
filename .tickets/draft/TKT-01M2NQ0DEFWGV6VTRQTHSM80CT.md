---
schema: 3
id: TKT-01M2NQ0DEFWGV6VTRQTHSM80CT
title: Add validated host configuration and legacy precedence
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
  - TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R
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
  - ref: file:internal/fetch/fetch.go
    path: internal/fetch/fetch.go
  - ref: file:README.md
    path: README.md
claim: null
archive: null
created_at: 2026-09-16T18:17:32Z
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

Move backend, URL, limits, User-Agent and allowlist settings to manifest/SDK configuration with validated immutable runtime snapshots and atomic updates. During a bounded transition preserve the proposed precedence: new TERVA_EXT_WEB_* overrides, legacy ZOT_WEB_* overrides, explicitly configured host values, legacy file values, defaults. Host-resolved defaults are not proof of explicit configuration. Secret handling is a separate dependent ticket.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Document and test precedence, preserving existing legacy choices and TAVILY_API_KEY behavior without introducing unsafe secret defaults
- [ ] Expose supported nonsecret settings with validation and distinguish explicit host values from resolved defaults or record a tested alternative import mechanism
- [ ] Apply config updates atomically; in-flight operations keep a coherent snapshot and rejected updates remain visible while working settings survive
- [ ] Define and test cache invalidation for backend, allowlist and User-Agent changes so stale entries cannot bypass tightened policy
- [ ] Cover first setup, malformed settings, missing backend credentials, legacy files and concurrent updates with synthetic fixtures only

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff
