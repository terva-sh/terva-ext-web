# Modernization backlog

This index translates the outstanding work in [the split plan](terva-ext-web.md)
into repository tickets. The tickets own status, acceptance criteria, dependencies,
decisions and validation; this page is a navigation and coverage map. All 22
tickets comprise two epics and 20 scoped tasks/investigations. The three
independent starting tickets were promoted by user approval; the save fix is now
complete. Ticket status is authoritative for subsequent progress.
Implementation plans are intentionally left for the agent who claims each task.

See [the critical-path grooming review](modernization-critical-path.md) for
recommended first selections, evidence, decision gates, and outstanding inputs.

## Filtering the backlog

All tickets below carry `initiative:terva-modernization`, exactly one of
`scope:core` or `scope:follow-up`, and relevant `area:*` labels. See
[Ticket labels](../ticket-labels.md) for the vocabulary, CLI examples, and
JSON filtering (repeated `--label` options use OR).

## Completed foundation

- Exact history-preserving split and baseline: see [identity-packaging.md](identity-packaging.md).
- Identity and packaging: PR #1 merged as `2ac463c`; implementation `6e1d43c`.
- Ticket store and workflow: PR #2 merged as `80d3ec6`.
- Public Forgejo repository and origin verified; research skill archives and matching version stamps already implemented.
- These are historical evidence, not duplicate unfinished tickets.

## Core modernization and release

[TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5 (Complete Terva modernization and replacement release)](../../.tickets/draft/TKT-01M2NQ0CRC76K1R0PKTGDZSYQ5.md)

SDK verification, the save fix and platform-access planning can begin independently
once selected. Configuration also waits for the provenance/import decision.

| Work | Prerequisites |
| --- | --- |
| [TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K (Verify a published Terva SDK and supported host floor)](../../.tickets/done/TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K.md) | None |
| [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md) | [TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K (Verify a published Terva SDK and supported host floor)](../../.tickets/done/TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K.md) |
| [TKT-01M2NQ0D3RZAMV83QH6AAX3466 (Keep download saves in the workspace captured at call start)](../../.tickets/done/TKT-01M2NQ0D3RZAMV83QH6AAX3466.md) | None |
| [TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8 (Preserve network and write authority during SDK adoption)](../../.tickets/done/TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8.md) | [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md) |
| [TKT-01M2NQ0DANEHS615M5VXNGS7W5 (Validate SDK conformance against supported Terva hosts)](../../.tickets/tickets/TKT-01M2NQ0DANEHS615M5VXNGS7W5.md) | [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md), [TKT-01M2NQ0D3RZAMV83QH6AAX3466 (Keep download saves in the workspace captured at call start)](../../.tickets/done/TKT-01M2NQ0D3RZAMV83QH6AAX3466.md), [TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8 (Preserve network and write authority during SDK adoption)](../../.tickets/done/TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8.md) |
| [TKT-01M2NQT535ETD2WACPHJ93PA48 (Resolve configuration provenance and legacy import semantics)](../../.tickets/done/TKT-01M2NQT535ETD2WACPHJ93PA48.md) | [TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K (Verify a published Terva SDK and supported host floor)](../../.tickets/done/TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K.md) |
| [TKT-01M2NQ0DEFWGV6VTRQTHSM80CT (Add validated host configuration and legacy precedence)](../../.tickets/draft/TKT-01M2NQ0DEFWGV6VTRQTHSM80CT.md) | [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md), [TKT-01M2NQT535ETD2WACPHJ93PA48 (Resolve configuration provenance and legacy import semantics)](../../.tickets/done/TKT-01M2NQT535ETD2WACPHJ93PA48.md) |
| [TKT-01M2NQ0DK05JAWQRF7TBGF1CY7 (Migrate Tavily configuration secrets without losing legacy settings)](../../.tickets/draft/TKT-01M2NQ0DK05JAWQRF7TBGF1CY7.md) | [TKT-01M2NQ0DEFWGV6VTRQTHSM80CT (Add validated host configuration and legacy precedence)](../../.tickets/draft/TKT-01M2NQ0DEFWGV6VTRQTHSM80CT.md) |
| [TKT-01M2NQT58ZQPSE95M1WEXPN8SF (Define release test matrix and secure platform test access)](../../.tickets/done/TKT-01M2NQT58ZQPSE95M1WEXPN8SF.md) | None |
| [TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A (Validate replacement installation, upgrades and release platforms)](../../.tickets/draft/TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A.md) | [TKT-01M2NQ0DANEHS615M5VXNGS7W5 (Validate SDK conformance against supported Terva hosts)](../../.tickets/tickets/TKT-01M2NQ0DANEHS615M5VXNGS7W5.md), [TKT-01M2NQ0DEFWGV6VTRQTHSM80CT (Add validated host configuration and legacy precedence)](../../.tickets/draft/TKT-01M2NQ0DEFWGV6VTRQTHSM80CT.md), [TKT-01M2NQ0DK05JAWQRF7TBGF1CY7 (Migrate Tavily configuration secrets without losing legacy settings)](../../.tickets/draft/TKT-01M2NQ0DK05JAWQRF7TBGF1CY7.md), [TKT-01M2NQT58ZQPSE95M1WEXPN8SF (Define release test matrix and secure platform test access)](../../.tickets/done/TKT-01M2NQT58ZQPSE95M1WEXPN8SF.md) |
| [TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7 (Publish and smoke-test the first terva-ext-web release)](../../.tickets/draft/TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7.md) | [TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A (Validate replacement installation, upgrades and release platforms)](../../.tickets/draft/TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A.md) |
| [TKT-01M2NQ0DZR9RMP1PRT1VDA7GQP (Coordinate legacy migration notice and authorized zot-web archival)](../../.tickets/draft/TKT-01M2NQ0DZR9RMP1PRT1VDA7GQP.md) | [TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7 (Publish and smoke-test the first terva-ext-web release)](../../.tickets/draft/TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7.md) |

## Follow-up features and investigations

[TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG (Improve Terva web tool discovery and presentation)](../../.tickets/draft/TKT-01M2NQ0CTGYDCTP37XTS4BQ0BG.md)

Recommended next features and optional investigations do not block the first
replacement release. Spikes may conclude with an evidence-backed defer/reject decision.

| Work | Prerequisites |
| --- | --- |
| [TKT-01M2NQ0E22NNZVMM8Y2K97J4G8 (Withdraw and restore search when its backend is unavailable)](../../.tickets/draft/TKT-01M2NQ0E22NNZVMM8Y2K97J4G8.md) | [TKT-01M2NQ0DEFWGV6VTRQTHSM80CT (Add validated host configuration and legacy precedence)](../../.tickets/draft/TKT-01M2NQ0DEFWGV6VTRQTHSM80CT.md) |
| [TKT-01M2NQ0E4XN889GZMRVT99R6V8 (Measure whether web tools need eager discovery hints)](../../.tickets/draft/TKT-01M2NQ0E4XN889GZMRVT99R6V8.md) | [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md) |
| [TKT-01M2NQ0E7SGTAC7T78V8P96068 (Add web request display hints and cache/backend status)](../../.tickets/draft/TKT-01M2NQ0E7SGTAC7T78V8P96068.md) | [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md), [TKT-01M2NQ0DEFWGV6VTRQTHSM80CT (Add validated host configuration and legacy precedence)](../../.tickets/draft/TKT-01M2NQ0DEFWGV6VTRQTHSM80CT.md) |
| [TKT-01M2NQ0EAS2HE1A5MT7Z3YY7X2 (Evaluate small extension-authored context guidance)](../../.tickets/draft/TKT-01M2NQ0EAS2HE1A5MT7Z3YY7X2.md) | [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md), [TKT-01M2NQ0DEFWGV6VTRQTHSM80CT (Add validated host configuration and legacy precedence)](../../.tickets/draft/TKT-01M2NQ0DEFWGV6VTRQTHSM80CT.md) |
| [TKT-01M2NQ0EDS22HXXZ68K3DK87XN (Assess host-brokered saves without weakening download semantics)](../../.tickets/draft/TKT-01M2NQ0EDS22HXXZ68K3DK87XN.md) | [TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R (Migrate protocol integration to the verified Terva SDK)](../../.tickets/done/TKT-01M2NQ0CZ2R2XCRHF3FPEPC14R.md), [TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8 (Preserve network and write authority during SDK adoption)](../../.tickets/done/TKT-01M2NQ0D6MHKE0QQ49C6TPKVP8.md) |
| [TKT-01M2NQ0EGSJJMK5CESDRDQV7MM (Decide whether private page caches need project isolation)](../../.tickets/draft/TKT-01M2NQ0EGSJJMK5CESDRDQV7MM.md) | [TKT-01M2NQ0DEFWGV6VTRQTHSM80CT (Add validated host configuration and legacy precedence)](../../.tickets/draft/TKT-01M2NQ0DEFWGV6VTRQTHSM80CT.md) |
| [TKT-01M2NQ0EKTNJX636M9X5JK6J0W (Report supported startup progress during source builds)](../../.tickets/draft/TKT-01M2NQ0EKTNJX636M9X5JK6J0W.md) | [TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K (Verify a published Terva SDK and supported host floor)](../../.tickets/done/TKT-01M2NQ0CW9Y7T06NN29JJ0PW1K.md) |
| [TKT-01M2NQ0EPWH96WEK8BDQ2HB6Y7 (Evaluate checksum-verified prebuilt fallback for source installs)](../../.tickets/draft/TKT-01M2NQ0EPWH96WEK8BDQ2HB6Y7.md) | [TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7 (Publish and smoke-test the first terva-ext-web release)](../../.tickets/draft/TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7.md) |

## Scope decisions retained from the assessment

- Ordered session identity is covered by SDK integration and the save-race ticket.
- Selective Sequential ordering is assessed with authority/concurrency; network reads remain concurrent.
- Mandatory configuration cache invalidation is in the config task; additional project isolation is a separate spike.
- Runtime secret brokering is conditional, not a default feature commitment. Config-secret migration comes first.
- Per-call cancellation and trust/provenance metadata remain documented SDK limits, not promised capabilities.
- Session list/read features stay deferred to index/memory extensions unless a concrete web use case emerges.
- Tool/event interception and user-message rewriting stay deferred; connector tunneling is skipped.
- Unrelated README feature ideas (new search backends, JS rendering, infobox extraction) are outside this split plan.
- Automated git-ticket CI provisioning is separate workflow maintenance, not a modernization release prerequisite.

Do not copy ticket status into this table. When ticket files move, update links
as part of the owning work; `git ticket search` and `.tickets/epics.md` provide
the current store view.
