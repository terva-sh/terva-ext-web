---
schema: 3
id: TKT-01M2P05BFER2BYNEAVRMBNQN5Y
title: Remove obsolete zot support and refresh Terva documentation
type: chore
status: done
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
claim: null
archive: null
created_at: 2026-09-16T20:57:31Z
updated_at: 2026-09-16T21:18:30Z
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

- [x] Remove obsolete zot-only code and references without weakening Terva safety checks
- [x] Document current configuration and migration behavior with any changed policy recorded explicitly
- [x] Apply unslop to maintained documentation and correct stale implementation and validation claims
- [x] Pass relevant regression and ticket checks and record native revalidation needs

## Implementation plan

Remove unused config.Load and all active legacy configuration paths as explicitly requested: Resolve takes only Terva host settings, TERVA_EXT_WEB overrides and TAVILY_API_KEY. Replace silent-clamping and legacy-mode tests with production resolver bounds/precedence checks and negative regressions for ignored variables/files. Use the published Terva configuration CLI in the SearXNG recipe. Apply /home/sothr/.claude/skills/unslop/SKILL.md to maintained documentation and the bundled research skill, correcting obsolete release instructions, permission/cache claims and status. Preserve generated workflow, vendored code, installed files and historical evidence. Pass local CI and link checks, then run native source/archive validation on all five targets and record candidate-bound reports in PR #17.

## Notes

**agent:codex/terva-cleanup** at 2026-09-16T20:59:54Z

Unchanged-source just ci passed. Removed unused config.Load, which still implemented zot-only environment parsing and silent clamping; replaced its tests with production Resolve coverage for defaults, all settings, numeric bounds through each source, file order and allowlist semantics. The active migration resolver remains pending the user preference. Also replaced configure-searxng path guessing and config.json overwrites with the published Terva ext config API, verified in v0.137.0 docs/extensions.md. This recipe now changes only search_backend and searxng_url; private non-loopback allowlist entries are explicit operator configuration.

**agent:codex/terva-cleanup** at 2026-09-16T21:04:33Z

User explicitly selected removal of active legacy configuration support, superseding the earlier 0.4.x migration-period decision. Remove file reads, ZOT_WEB aliases and configuration_source; Resolve accepts only the host map plus current environment overrides. Keep TAVILY_API_KEY and data_secrets:true because old files may remain. No installed file is read, rewritten or deleted. Rollback now requires the old installation, with mutually exclusive enablement. Added negative tests for retired variables and subprocess fixtures proving valid/malformed old files are ignored and preserved.

**agent:codex/terva-cleanup** at 2026-09-16T21:10:17Z

Post-change just ci passed: vet/format, race tests, six-tool subprocess conformance including ignored old-file fixtures, published host-driver/policy checks and clean vendor regeneration. All 18 maintained Markdown files have valid local links; generated AGENTS ticket workflow is byte-identical to its previous version. Unslop pass removed stale dual-host and release-cut instructions, corrected yolo/cache/image-limit claims, and rewrote current setup, migration, direction and critical-path documentation. Native rerun will bind this changed manifest/source/skill archive to fresh evidence.

**agent:codex/terva-cleanup** at 2026-09-16T21:11:32Z

Final review tightened the retired-file regression fixture: it selects an unconfigured SearXNG backend as well as a synthetic key, so even a regression that restores file loading cannot send a request to the real Tavily service. The focused subprocess regression passes. Restart native validation for this final test revision.

## Summary

Removed the unused loader and active legacy file/environment/mode support. Terva host configuration and current env overrides are the only inputs; existing files remain untouched. SearXNG helper uses the host API. Applied the requested unslop skill to maintained docs and research guidance, corrected safety/status claims, and updated release migration criteria. Baseline and final just ci, documentation links and strict ticket checks pass. Native run 35150980130 passed all five source/archive targets at b162e831b076d85084f89c701293fdfbb436f129; ten reports and checksums are preserved in docs/validation/native-terva-only-2026-09-16.json. Code and documentation commits are separate in Forgejo PR #17.
