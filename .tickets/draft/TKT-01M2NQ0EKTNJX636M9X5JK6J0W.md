---
schema: 3
id: TKT-01M2NQ0EKTNJX636M9X5JK6J0W
title: Report supported startup progress during source builds
type: task
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:launcher
  - area:presentation
assignees: []
milestone: null
parent: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
origin: null
dependencies:
  - TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K
blocks_on: none
references:
  - ref: file:docs/plans/terva-ext-web.md
    path: docs/plans/terva-ext-web.md
  - ref: file:run.sh
    path: run.sh
  - ref: file:.goreleaser.yaml
    path: .goreleaser.yaml
  - ref: file:conformance_test.go
    path: conformance_test.go
claim: null
archive: null
created_at: 2026-09-16T18:17:34Z
updated_at: 2026-09-16T18:31:38Z
created_by:
  id: agent:codex/modernization-tickets
  name: ""
updated_by:
  id: agent:codex/critical-path
  name: ""
extensions: {}
---

## Description

The offline launcher can spend time compiling before hello. Adopt the host-supported bootstrap progress mechanism only after verifying the launcher/startup contract; do not inject arbitrary status text into the JSON wire. Research skill archive inclusion is already complete.

Source: docs/plans/terva-ext-web.md. This is scoped backlog work, not an implementation plan; inspect current code and record the approach after claiming it.

## Acceptance criteria

- [ ] Verify the supported pre-hello progress API and host fallback before changing launcher output
- [ ] Report first-build/rebuild progress without corrupting protocol stdout or hanging unsupported hosts
- [ ] Preserve offline vendor compilation, useful stderr errors and immediate execution of prebuilt archives
- [ ] Exercise first build, rebuild, missing toolchain and prebuilt paths
- [ ] Do not assume unknown bootstrap frames are ignored: prove startup behavior on every supported host floor or deliberately gate/omit frames before hello when support cannot be known

## Definition of done

- [ ] Record decisions and validation evidence in the ticket; commit intended changes and pass git ticket check before handoff

## Notes

**agent:codex/critical-path** at 2026-09-16T18:31:38Z

Grooming review, 2026-09-16 (not an implementation plan).

Entry inputs: SDK verification supplies actual host-floor/bootstrap support; current offline launcher.

Expected deliverable: Launcher progress only under a proven compatible startup contract.

Validation: First build/rebuild/prebuilt/no-Go paths and old-host behavior with bounded startup time.

Readiness/coordination: Terva docs explicitly say an unknown bootstrap frame can be treated as malformed hello; compatibility cannot be handshaken before that frame.
