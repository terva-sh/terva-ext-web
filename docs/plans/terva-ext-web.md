# Terva modernization direction

terva-ext-web is a Terva-only fork of zot-web. It retains the original history
and web tool behavior while using the published Terva SDK, host configuration
and independent packaging. The [ticket store](../../.tickets/) records current
work; the [critical path](modernization-critical-path.md) lists release tasks.

## Fork and identity

The exact fork point is `c50d773c60250a7315c37c2aedbeb891cb12f6a4`, tree
`6c368f89b3e622c7d4515ac4e02d47caed3a95c7`. It used handwritten
`internal/proto`; SDK adoption was later work, not part of the inherited tree.
The original assessment inspected Terva revision
`7f754b9bb7c6754284dc4f3cc5fa6a96525de49d`. Published version verification
subsequently established v0.137.0; see [the verification record](sdk-verification.md).

Repository, binary, module and default User-Agent use `terva-ext-web`.
Manifest and handshake name remain `web`, with all six `web_*` tools and
`/web-cache`. Shared identities mean old and new installations must not run
together. Source changes do not rename installations or move credentials.

[PR #1](https://git.local.sothr.com/terva-sh/terva-ext-web/pulls/1) completed
identity and packaging as `2ac463c`, with implementation `6e1d43c`.
[PR #2](https://git.local.sothr.com/terva-sh/terva-ext-web/pulls/2) added the
ticket workflow as `80d3ec6`. [Identity and packaging](identity-packaging.md)
records the unchanged-source baseline and checks. The initial lack of Go on
PATH was a validation gap, later resolved with the installed toolchain.

Forgejo is the source repository. GitHub mirrors the same history and runs
native release validation. Keep inherited `cut/1` through `cut/6` as records;
new releases use version tags and archives. Do not rewrite the fork as an
orphan release branch or publish it over zot-web.

## Implemented decisions

The extension uses `terva.sh/terva v0.137.0`, Go 1.27 and committed vendor
sources. Minimum protocol 2 supplies ordered session identity. Both the tested
host floor and current baseline are v0.137.0. Do not infer optional capabilities
from a version or promise per-call cancellation or result trust metadata that
the SDK does not supply.

Every network operation captures a validated runtime. Host updates create a
new provider and cache; in-flight calls finish with their original settings.
Every save also captures its workspace before fetching. The session-switch
regression proves that a later session cannot redirect that write.

The user removed the planned legacy configuration period before release.
Only Terva host configuration, `TERVA_EXT_WEB_*` overrides and `TAVILY_API_KEY`
remain. Old files are ignored and preserved for manual migration or rollback
with the old installation. See [configuration](config-provenance.md) and
[credentials](credential-migration.md).

All tools keep `network-read` authority; raw and image tools also request
manifest `ask`. Yolo can allow them. Application destination, path and resource
checks apply after host approval. See [the authority contract](authority-contract.md).
Output sanitization limits forged metadata but does not make fetched text
trusted; see [untrusted content](untrusted-web-content.md).

SDK and native validation passed before this cleanup. Reports apply to their
recorded candidates and archive checksums. Changes require revalidation as
described in [platform validation](platform-validation.md).

## Remaining work

Complete full-CLI source/archive installation, manual settings migration and
rollback checks. Preserve the old installation until the replacement passes;
never enable both `web` identities together. Publish the approved release only
after native checks pass for its exact tagged archives. Add a migration notice
to zot-web afterward; archival requires separate ticket authorization.

Later tickets cover unavailable-search withdrawal, discovery hints, request
status, small context guidance, host-brokered saves, project cache isolation,
source-build progress, prebuilt fallback and ddgr search. Those are outside the
first release unless selected explicitly. Runtime readiness must remain correct
regardless of tool visibility. Host-brokered saves require proof that path,
overwrite and byte-preservation rules survive the change.

The fork keeps existing fetch/search implementations to avoid coupling SDK work
to a retrieval rewrite. A history-preserving split keeps the old installation
available for rollback and makes each implementation batch reviewable.
