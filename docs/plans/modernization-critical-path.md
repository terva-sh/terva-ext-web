# Modernization critical path

The SDK, session-safe saves, host configuration, credential handling, authority
checks and native source/archive validation are implemented. Native reports
cover the recorded candidate, not subsequent source or packaging changes.
The user selected Terva-only configuration cleanup and a documentation pass
before installation validation. The ticket store holds live status.

## Release order

1. Complete the selected cleanup and validate its source and archives.
2. [Validate installation, upgrades and release platforms](../../.tickets/draft/TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A.md).
   Use the actual supported Terva CLI for isolated source/archive installs,
   manual settings migration and rollback. Check skill discovery, preserve old
   settings and never enable both `web` installations together.
3. [Publish and smoke-test the first release](../../.tickets/draft/TKT-01M2NQ0DVK4ZKJ8QCET28XNTV7.md).
   Prepare the concrete versioned candidate, get release approval, and require
   native checks before publication. Test downloaded assets afterward.
4. [Coordinate legacy notice and archival](../../.tickets/draft/TKT-01M2NQ0DZR9RMP1PRT1VDA7GQP.md).
   A migration notice follows the verified replacement. Remote archival needs
   explicit ticket authorization in the old repository.

## Inputs already established

| Input | Evidence |
| --- | --- |
| Published SDK, host version and toolchain | [SDK verification](sdk-verification.md) |
| Configuration precedence and host form behavior | [Configuration](config-provenance.md) |
| Credential ownership and preservation | [Credential migration](credential-migration.md) |
| Host approval behavior and workspace saves | [Authority contract](authority-contract.md) |
| Wire behavior and published host driver | [SDK conformance](sdk-conformance.md) |
| Native runners and archive validation | [Platform validation](platform-validation.md) |

Linux amd64 is the development target. GitHub supplies native Linux arm64,
macOS amd64/arm64 and Windows amd64 runners. Manual rehearsals validate without
publishing; release tags must validate their own archives.

The original grooming separated SDK, configuration and platform preparation
so missing access could be resolved before release. That work is complete.
Installation/upgrade/rollback evidence and publication approval remain open.

Follow-up discovery, presentation, context, host-write, cache-isolation,
prebuilt-fallback and ddgr work stays outside the release path unless selected.
See [the backlog](modernization-backlog.md). New findings need tickets and
recorded rationale, not an untracked extension of the current batch.
