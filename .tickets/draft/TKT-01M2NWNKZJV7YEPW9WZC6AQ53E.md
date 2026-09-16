---
schema: 3
id: TKT-01M2NWNKZJV7YEPW9WZC6AQ53E
title: Evaluate ddgr DuckDuckGo search as a zero-config default
type: spike
status: draft
status_reason: null
priority: normal
due_on: null
labels:
  - initiative:terva-modernization
  - scope:follow-up
  - area:config
  - area:discovery
  - area:security
assignees: []
milestone: null
parent: TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG
origin: null
dependencies:
  - TKT-01M2NQ0DEFWGV6VTRQTHSM80CT
blocks_on: none
references: []
claim: null
archive: null
created_at: 2026-09-16T19:56:30Z
updated_at: 2026-09-16T19:56:44Z
created_by:
  id: agent:codex/modernization-run
  name: ""
updated_by:
  id: agent:codex/modernization-run
  name: ""
extensions: {}
---

## Description

Consider adopting the ddgr-backed DuckDuckGo search idea from ggbalaazs/zot-web so a fresh installation can search without provider credentials when ddgr is installed. User requested this be filed for later, not added to the current release batch.

Reference: https://github.com/terva-sh/zot-web/compare/release...ggbalaazs:zot-web:release

GitHub comparison inspected on 2026-09-16: one commit ahead, including new internal/search/duckduckgo.go and tests plus config, documentation and wiring changes. Review implementation and license/attribution before borrowing; port the idea to the current SDK/config runtime rather than copying legacy wiring.

Explore automatic ddgr availability detection and a sensible fresh-install default. Preserve explicitly configured Tavily/SearXNG and existing legacy choices. Define behavior when ddgr is absent, removed, fails, times out, or produces malformed/oversized output; do not auto-install software. Review subprocess execution and egress separately because ddgr will not automatically inherit the Go HTTP client's SSRF/proxy/allowlist controls.

## Acceptance criteria

- [ ] Review referenced fork and record reusable behavior, license/attribution and an adopt/defer/reject decision
- [ ] Specify fresh-install default selection when ddgr is present/absent while preserving explicit and legacy backend choices
- [ ] Map query/count/freshness/domain/depth options and result normalization; report unsupported filters honestly
- [ ] Define argv-only invocation, timeout/cancellation, stdout/stderr bounds, executable lookup and egress-policy implications
- [ ] Define hermetic fake-ddgr tests for success, missing executable, nonzero exit, malformed/oversized output and backend precedence

## Definition of done

- [ ] Record decision and file scoped implementation follow-ups; this remains outside the first release critical path

## Notes

**agent:codex/modernization-run** at 2026-09-16T19:56:44Z

Reference commit at inspection: 41d9d65c180f1bb11af7a958630a50f0c0c8ef73 (feat: add duckduckgo search as default). User explicitly requested later consideration; keep draft and do not gate current release.
