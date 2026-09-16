---
schema: 3
id: TKT-01M2NZD0DP9MTFSJ62F8BZV96V
title: Invoke Bash explicitly for native Windows host launch
type: bug
status: done
status_reason: null
priority: high
due_on: null
labels:
  - initiative:terva-modernization
  - scope:core
  - area:launcher
  - area:validation
assignees: []
milestone: null
parent: TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5
origin: null
dependencies: []
blocks_on: none
references: []
claim: null
archive: null
created_at: 2026-09-16T20:44:14Z
updated_at: 2026-09-16T20:51:32Z
created_by:
  id: agent:codex/modernization-run
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Native Windows rehearsal 35147849983 passes extension tests but the published Terva driver rejects exec ./run.sh as an invalid Win32 application. Change the manifest to exec bash with args ./run.sh, preserving the existing launcher and documented Git Bash prerequisite. Do not bypass the real host test or alter its manifest.

## Acceptance criteria

- [x] Source and archive manifests launch through the actual Terva driver on all five targets
- [x] Retain no-Go archive launch, path-with-spaces support and offline source build behavior

## Implementation plan

Use the published host manifest Exec/Args contract to invoke Bash explicitly. This works on Unix and Windows with Git Bash and preserves a single launcher. A separate Windows executable launcher would duplicate build logic; bypassing the launcher only in tests would conceal the installation bug. Validate the unchanged actual-driver fixture locally and in all native jobs.

## Notes

**agent:codex/modernization-run** at 2026-09-16T20:44:47Z

Changed the source manifest to exec bash with ./run.sh as an argument; archive packaging carries the same manifest. Actual published host-driver Linux test passes without fixture substitutions. Native validation rerun required to prove Windows behavior.

## Summary

Manifest explicitly invokes Bash with ./run.sh. Actual published Terva driver passes source and archive launch on all five targets, including Windows Git Bash and paths with spaces; archive launch without Go passes. Candidate 7fde307, run 35148398553, reports in docs/validation/native-2026-09-16.json. No test-only manifest substitution.
