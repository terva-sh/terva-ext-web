---
schema: 3
id: TKT-01M2P05BFER2BYNEAVRMBNQN5Y
title: Remove obsolete zot support and refresh Terva documentation
type: chore
status: in-progress
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:config
  - area:workflow
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies: []
blocks_on: none
references: []
claim:
  actor: agent:codex/terva-cleanup
  branch: cleanup/terva-only-docs
  worktree: /home/sothr/.t3/worktrees/terva-ext-web/t3code-7d71b129
  commit: aafa038811ec350c58930519054099a651e39c9d
  session: null
  claimed_at: 2026-09-16T20:57:39Z
  expires_at: null
archive: null
created_at: 2026-09-16T20:57:31Z
updated_at: 2026-09-16T21:11:32Z
created_by:
  id: agent:codex/terva-cleanup
  name: ""
updated_by:
  id: agent:codex/terva-cleanup
  name: ""
extensions: {}
---

## Description

User requested a review of zot-only support and an unslop documentation pass. Audit active behavior, remove obsolete compatibility implementation, update Terva setup/security/release guidance, and distinguish historical records from current requirements. Apply the installed unslop skill without rewriting generated ticket workflow or third-party documentation. Confirm the migration-support boundary before changing documented 0.4.x behavior.

## Acceptance criteria

- [ ] Remove obsolete zot-only code and references without weakening Terva safety checks
- [ ] Document current configuration and migration behavior with any changed policy recorded explicitly
- [ ] Apply unslop to maintained documentation and correct stale implementation and validation claims
- [ ] Pass relevant regression and ticket checks and record native revalidation needs

## Implementation plan

Baseline just ci, then remove unused config.Load and replace its obsolete clamping tests with production Resolve coverage. Update operator errors to Terva configuration names. Review remaining compatibility against the documented 0.4.x migration contract, preserving active migration inputs unless the user selects their removal. Apply the installed /home/sothr/.claude/skills/unslop/SKILL.md to README, project AGENTS prose, maintained guides and the bundled skill. Replace obsolete release instructions with the current process and label historical evidence. Correct yolo permissions, cache scope and completed native validation claims. Preserve generated workflow, ticket history, vendored code and recorded evidence. Run just ci, document-link checks and native revalidation for changed source/archives.

## Notes

**agent:codex/terva-cleanup** at 2026-09-16T20:59:54Z

Unchanged-source just ci passed. Removed unused config.Load, which still implemented zot-only environment parsing and silent clamping; replaced its tests with production Resolve coverage for defaults, all settings, numeric bounds through each source, file order and allowlist semantics. The active migration resolver remains pending the user preference. Also replaced configure-searxng path guessing and config.json overwrites with the published Terva ext config API, verified in v0.137.0 docs/extensions.md. This recipe now changes only search_backend and searxng_url; private non-loopback allowlist entries are explicit operator configuration.

**agent:codex/terva-cleanup** at 2026-09-16T21:04:33Z

User explicitly selected removal of active legacy configuration support, superseding the earlier 0.4.x migration-period decision. Remove file reads, ZOT_WEB aliases and configuration_source; Resolve accepts only the host map plus current environment overrides. Keep TAVILY_API_KEY and data_secrets:true because old files may remain. No installed file is read, rewritten or deleted. Rollback now requires the old installation, with mutually exclusive enablement. Added negative tests for retired variables and subprocess fixtures proving valid/malformed old files are ignored and preserved.

**agent:codex/terva-cleanup** at 2026-09-16T21:10:17Z

Post-change just ci passed: vet/format, race tests, six-tool subprocess conformance including ignored old-file fixtures, published host-driver/policy checks and clean vendor regeneration. All 18 maintained Markdown files have valid local links; generated AGENTS ticket workflow is byte-identical to its previous version. Unslop pass removed stale dual-host and release-cut instructions, corrected yolo/cache/image-limit claims, and rewrote current setup, migration, direction and critical-path documentation. Native rerun will bind this changed manifest/source/skill archive to fresh evidence.

**agent:codex/terva-cleanup** at 2026-09-16T21:11:32Z

Final review tightened the retired-file regression fixture: it selects an unconfigured SearXNG backend as well as a synthetic key, so even a regression that restores file loading cannot send a request to the real Tavily service. The focused subprocess regression passes. Restart native validation for this final test revision.
