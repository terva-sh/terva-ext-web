# Proposed split: zot-web → terva-ext-web

Status: identity/packaging merged in [PR #1](https://git.local.sothr.com/terva-sh/terva-ext-web/pulls/1)
as `2ac463c` (implementation `6e1d43c`); ticket workflow merged in
[PR #2](https://git.local.sothr.com/terva-sh/terva-ext-web/pulls/2) as `80d3ec6`.
The outstanding work is tracked in [the modernization backlog](modernization-backlog.md)
and `.tickets/`; [the critical-path review](modernization-critical-path.md)
records prerequisites and readiness gaps. Those tickets own current work status and dependencies; the
assessment below preserves direction and original evidence.

Origin is the server-verified public repository
`ssh://git@git.local.sothr.com:2222/terva-sh/terva-ext-web.git`.
The API and `git ls-remote` confirmed it was empty before the initial push.
Remote main began at `b6efdecbd7439543a7d39dc8a77e609f887002f9`, before
identity implementation. No binary release has been published and no
repository has been archived. This document was moved from the zot-web
worktree after cloning the exact fork point; it is not part of that tree.

## Local setup completed

The independent working clone is at
`/home/sothr/workspace/git.local.sothr.com/terva-sh/terva-ext-web`, adjacent
to the primary zot-web checkout. It started on `main` at the commit and tree
recorded below, with complete ancestry and its own object database (no
alternates or shallow history). `git fsck --full` passed before adding
orientation documents. Historical tags `cut/1` through `cut/6` are retained;
these are inherited zot-web release records, not new terva-ext-web releases.

At initial local setup, no remote was configured. The temporary clone source was removed to
prevent accidental publication back to zot-web. The Forgejo repository was subsequently
created and its server-advertised clone URL and ownership verified as recorded
above. The source checkout is unchanged.

Local setup is tracked in the workspace ledger under
TKT-01M2NK39WCKXF4JAQV3FTCV887 (Create local terva-ext-web fork and startup directions).
The root `AGENTS.md` directs implementation agents. Local setup, identity and
packaging, remote creation, and ticket workflow are complete. The identity
batch and passing baseline are recorded in [identity-packaging.md](identity-packaging.md).
SDK migration and the remaining batches are draft work in the backlog index.

## Recommendation and evidence

Create terva-ext-web with zot-web's history through the exact commit below.
Make identity, SDK, and feature changes in subsequent commits in the new
repository. Archive zot-web only after the replacement has a verified
installation and release, with an archive notice pointing to it.

The compatibility commitment being retired is support for stock zot and its
older extension protocol. Terva remains the target; newer Terva protocol
features still need their documented version or capability checks.

Inspected source:

- zot-web commit: `c50d773c60250a7315c37c2aedbeb891cb12f6a4`.
- Exact source tree: `6c368f89b3e622c7d4515ac4e02d47caed3a95c7`.
- Manifest version: `0.3.1`; manifest/handshake name: `web`.
- This checkout was clean before this proposal. The primary checkout's
  `main` and this worktree started at the same commit.
- Local Terva source: `7f754b9bb7c6754284dc4f3cc5fa6a96525de49d`, described
  by Git as `v0.137.0-1-g7f754b9b`; extension protocol **6**.
- The locally available `v0.137.0` tag also contains `WithDisplay`,
  `Essential`, `Sequential`, `OnConfig`, `OnCompaction`, and `HostToolCall`.
  Remote branch freshness and public module availability were not checked.
- Sibling `terva-ext-caldav` uses module/binary naming `terva-ext-caldav`,
  manifest name `caldav`, and `terva.sh/terva v0.131.8`. The older template
  pins `v0.108.0`; copying its pin would miss newer APIs.

There is a discrepancy with the premise that this extension already imports
the Terva SDK. At this fork point, `go.mod` has no Terva dependency and
`main.go` imports the handwritten `internal/proto`. That package implements
a subset of protocol 2. Its conformance profiles simulate zot protocol 1
and early Terva protocol 2, not current Terva. SDK adoption is work still
to do in this tree. If another branch contains that migration, port it
*after* the recorded fork point rather than changing the baseline silently.

Evidence paths in zot-web: `go.mod`, `internal/proto/proto.go`,
`conformance_test.go`, `main.go`, `internal/config/config.go`,
`extension.json`, `.goreleaser.yaml`, and `run.sh`.
Evidence paths in Terva at the revision above: `docs/extensions.md`,
`packages/agent/ext/ext.go`, `packages/agent/extproto/extproto.go`, and
`packages/agent/extensions/`. These identify the implementation examined,
not a promise that every local change is available in a published module.

## Split and migration sequence

1. Record the source commit and tree above in the split's implementation
   ticket. Create an independent repository with full ancestry through that
   commit, with its initial default-branch tip at that commit. Avoid a
   shallow export or source-only copy, which loses ancestry. Avoid a mirror
   push of every ref, which could carry unrelated worktree branches.
2. Verify the destination path is unused and the destination remote is empty
   before creation/publication. The documented workspace layout and existing
   source checkout place siblings under
   `/home/sothr/workspace/git.local.sothr.com/terva-sh/`; the proposed new
   child is `terva-ext-web`. The existing origin identifies the forge owner
   as `terva-sh`. Verify creation and visibility settings at execution time.
3. Prove destination `HEAD` equals the source commit and `HEAD^{tree}` equals
   the recorded tree, and run `git fsck`. Record which historical release tags
   are retained. Git history does not copy release assets, issues, CI secrets,
   branch protections, or mirrors; inventory those separately if needed.
4. In a new commit, rename the repository-facing identity to `terva-ext-web`.
   Following the sibling convention, use Go module `terva-ext-web` and update
   all internal imports and linker symbol paths together. Keep manifest and
   handshake name `web`, the six `web_*` tool names, and `/web-cache`.
   Keeping these preserves permission-rule and tool-call identities.
5. Update the launcher, binary name, default User-Agent, just recipes,
   release scripts, `.gitignore`, CI, archive names, install instructions,
   and release documentation. Supersede historical statements that zot
   support and zot naming are permanent; retain historical decision records
   as history. Choose the first new release version explicitly; `0.4.0` is
   a reasonable continuation after the inherited `0.3.1`.
6. Pin and vendor a verified published Terva SDK version, migrate the protocol
   integration, and validate a new installation and an upgrade with existing
   settings. Prefer checking `v0.137.0` first, since the inspected local tag
   contains the proposed APIs. Do not ship a local `replace` directive.
7. Publish and exercise release archives before directing existing users to
   migrate. Add a migration/archive notice to zot-web afterward. Archiving
   is a remote settings change and needs explicit authorization in the
   implementation ticket under the machine-wide repository rules.

Retaining `web` means the old and new extensions must not be enabled together:
their tool, command, configuration, and secret identities overlap. Migration
must identify installations through `terva ext list` and host-reported
directories, disable the old installation, and enable the new one. Do not
guess data paths from the repository basename. Existing docs disagree about
whether data is keyed by `web` or `zot-web`; the host's resolved path wins.
Preserve old settings and installation until the new install passes a smoke
test. Rollback disables the new install and re-enables the old one.

At assessment time, GoReleaser named `warricksothr/zot-web` and archives
omitted `skills/`. The identity batch resolved both: the destination is
`terva-sh/terva-ext-web` and archives include the research skill. Binary
publication stays disabled pending the release-validation tickets.

## Feature assessment

These priorities reflect value to a web-tools extension. A higher protocol
number does not by itself make a feature useful.

| Surface in inspected Terva | Current web extension | Recommendation |
| --- | --- | --- |
| Shared `ext` SDK and `extproto` types | Handwritten framing, handshake, dispatch, and result types | **First:** replace `internal/proto`; retain fetch/search logic and behavior tests. Gain shared framing recovery and future API updates. |
| Manifest config, `Config()`, `OnConfig` | JSON files and env overrides, initialized once | **First:** guided backend, URL, limits, User-Agent and allowlist settings; build a validated immutable runtime configuration and swap it atomically on updates. |
| Secret config fields; v6 runtime secret broker | Tavily key may live in extension `config.json` | **First:** use a `secret` config field for a user-supplied Tavily key. Use the v6 broker only for credentials acquired at runtime. Config secrets and broker secrets are separate APIs. |
| `data_secrets` declaration | Absent; legacy config can contain a key | **First:** establish secret-free data before declaring `false`. Never label a migrated directory clean while a legacy credential file remains there. |
| Authority classes and permission rules | All six tools are `network-read`, including writers | **First:** retain network gating and ask defaults; assess combined network/write effects. The wire supplies one authority string, not a set of effects. Changing writers to workspace-mutation alone can lose network gating. |
| Ordered lifecycle events, `OnSession`, `Host().Emits`, project data helpers | Tracks session cwd; process-global page cache | **First:** preserve session behavior and capture one workspace identity per operation. Later decide whether private-page caches should be scoped per project/configuration. |
| v4 tool withdrawal and restore | Search stays visible when unconfigured | **Next:** withdraw only unavailable search at supported session boundaries; keep usable fetch tools. Runtime readiness checks remain necessary, including after config changes. |
| `Essential()` with host cap and lazy tools | No visibility hints | **Next:** measure whether search/fetch need eager visibility; avoid pinning all six tools. Existing tool discovery may suffice. |
| `WithDisplay`, context cards, status segments, panels/widgets | Plain tool output and `/web-cache` text | **Next:** add query/URL subjects and a cache/backend status display. Keep text command output for clients without panels. Display hints do not change model output or permissions. |
| Static context and v3 `RefreshContext`; compaction events | Bundled research skill only | **Optional:** small extension-authored guidance about backend availability and paging. Avoid copying fetched page bodies into system-prompt contributions. Keep frequent cache changes out of the prompt. |
| v3 `HostToolCall` | Writes files directly through workspace guards | **Investigate:** broker text writes through host policy if overwrite, path, size and raw-byte semantics can be preserved. Arbitrary binary downloads are not automatically compatible with a text write tool. |
| v3 session list/read APIs | No session index | **Defer:** past-session search belongs in the existing index/memory extensions unless a concrete web research use case requires it. |
| Tool/event interception and user-message hooks | Not used | **Defer:** a fetch extension has no need to rewrite unrelated tools or user messages. |
| v5 experimental connector tunneling | No chat connector role | **Skip:** web retrieval is an extension role. `connproto`/`connsdk` target chat transports; they are not substitutes for `extproto`/`ext`. |
| Declarative bundles and bootstrap progress | Research skill; offline source launcher | **Next:** include skill in archives and emit supported startup progress for builds. Evaluate checksum-verified prebuilt fallback separately. |
| SDK tool ordering (`Sequential`) | Concurrent tool goroutines | **Selective:** keep network reads concurrent. Evaluate ordering for stateful commands or writes; serializing every request would unnecessarily slow retrieval. |

Config migration must define precedence explicitly. Proposed order during a
bounded transition: new `TERVA_EXT_WEB_*` env overrides, legacy `ZOT_WEB_*`
overrides, explicitly set host values, legacy file values, then defaults.
Keep the provider-standard `TAVILY_API_KEY` variable. The host sends resolved
defaults, so do not infer that a field was explicitly configured merely from
its presence: prove the import can preserve legacy choices before switching
to host config. Never include credential values in reports, tickets, or
migration logs. Do not delete or rewrite legacy credential files as an
incidental rename operation. At-rest encryption of config secrets depends on
host setup; a masked `secret` field alone does not prove encryption.

Changing the allowlist, backend, or User-Agent needs an explicit cache policy:
old entries must not bypass tightened egress policy or return a snapshot from
the wrong configuration. Validate new settings before replacing a working
configuration, while ensuring rejected settings are visible to the user.

## Correctness gaps and limits

The current save handlers call `e.CWD()` before the fetch and again afterward.
A session switch during the network request can therefore validate in one
workspace and save in another. Capture cwd once for both checks and writes,
and test the switch while a fetch is blocked. The SDK migration alone does
not correct this application-level race.

Do not advertise cancellation as an SDK migration benefit yet. At the inspected
revision, `ext.ToolHandler` accepts only JSON arguments and returns a result;
it has no call context. The extension currently creates background-derived
timeouts. Host queue cancellation does not undo an already accepted invocation.
Cooperative per-call cancellation requires separate protocol/SDK work or a
clearly defined extension-level mechanism.

Likewise, the inspected SDK result surface has text/image content blocks, with
no dedicated trust/provenance field. Keep the existing SSRF, redirect,
resource-limit, output-sanitization, and workspace-write protections. The
decision in `untrusted-web-content.md` remains applicable; a repository rename
does not establish a new model trust boundary.

Choose the minimum supported Terva version after confirming the SDK release.
Protocol 2 is the correctness floor for ordered session identity; protocol 3
is needed for host tool requests; protocol 6 becomes a hard floor if runtime
brokered secrets become required. Optional events should use the advertised
event set, and optional display/visibility features should degrade as their
API specifies. Removing stock zot testing does not require rejecting every
older Terva host.

## Reviewable implementation batches and validation

1. **Exact repository split:** commit/tree/ancestry proofs, destination metadata
   inventory, unchanged-source baseline tests. This batch changes no source.
2. **Identity and packaging:** build/version/hello agree, module imports and
   linker stamps resolve, source launcher and each supported release archive
   load, archive includes the research skill, release destination is correct.
3. **SDK and correctness:** race tests plus real subprocess conformance for
   the supported Terva floor and current host; assert handshake, all six tools,
   text and image results, malformed/oversized-frame recovery, concurrent calls,
   session-switch saves, and shutdown. Keep stdout JSON-only. Retire zot-only
   profiles after the new support contract is documented.
4. **Config and credentials:** exercise first setup, legacy precedence,
   malformed settings, missing backend credentials, config updates during
   in-flight calls, cache invalidation, and logs free of credentials. Verify
   migration with synthetic fixtures, never copied real credentials.
5. **Visibility and presentation:** confirm absent-search withdrawal and later
   restoration at supported boundaries, lazy tool discovery, display fallback,
   and no unexpected prompt growth from status/cache changes.
6. **Release and archive:** new-install and upgrade smoke tests, rollback proof,
   migration notice, then ticket-authorized archival.

The existing gates are race tests, tagged conformance, vet/format checks,
vendor synchronization and release snapshots. The initial assessment could
not run Go tests because Go was absent from PATH. The identity batch later
located the installed toolchain and passed unchanged-source and post-change
checks; see its evidence record. Current-host SDK compatibility and all-target
runtime validation remain outstanding tickets.

Alternatives considered: renaming zot-web in place obscures the stable legacy
endpoint; retaining dual-host compatibility keeps the maintenance constraint
the split is intended to remove; rewriting the fetcher adds risk without
helping protocol adoption. A history-preserving split followed by small SDK
and product changes provides a clear rollback boundary for each decision.
