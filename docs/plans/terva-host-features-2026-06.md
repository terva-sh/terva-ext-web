# June 2026 host integration record

This file records work inherited from zot-web. It does not define support for
stock zot in terva-ext-web. The fork replaced the handwritten protocol package
with the published Terva SDK; see [SDK verification](sdk-verification.md).

The inherited integration added `network-read` authority on all six tools,
the bundled web research skill, and manifest `ask` rules for raw and image
writers. Those behaviors remain, now covered by SDK conformance and published
host-driver tests.

The original permission account said yolo honored manifest `ask`. Testing
Terva v0.137.0 corrected that claim: yolo allows the tools unless a user deny
applies. Plan mode denies network tools even with a user allow. Current behavior
is recorded in [the authority contract](authority-contract.md).

The old dual-host naming policy, conditional zot wire compatibility and
release-cut workflow are retired. The extension uses Terva host configuration
and `TERVA_EXT_WEB_*` overrides. Historical `ZOT_WEB_*` aliases and standalone
configuration files are no longer read. The repository history preserves the
original implementation and decisions.
