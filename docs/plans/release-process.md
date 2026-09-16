# terva-ext-web release process

The first new release version is 0.4.0. Use `just release-verify` to validate
GoReleaser configuration and `just release-snapshot` for local five-target
archives and SHA-256 checksums. Snapshots use 0.4.0-next independently of the
inherited cut/* tags. Archives carry matching binary/manifest versions, the
launcher, README, license, and bundled research skill.

The planned Forgejo owner/repository is **terva-sh/terva-ext-web**. The inherited
warricksothr/zot-web destination is retired. No remote exists yet; verify the
actual server-advertised clone URL, owner, and empty destination before any
publication. GoReleaser publication and tag-triggered release CI are disabled.
The manual release workflow produces snapshots only and needs no write token.

A later release batch must validate the supported Terva host, new installation,
settings-preserving upgrade and rollback, and each supported platform's launch.
Then configure the verified new remote and enable publication deliberately.
Do not reuse cut/* tags, create an orphan public history, publish to zot-web,
or archive the old repository as part of this batch.

See [identity and packaging](identity-packaging.md) for decisions and checks.

---

## Historical zot-web release process (superseded)

The following is retained solely as an inherited decision record. Its paths,
commands, permanent naming policy, and release-cut flow do not apply to this
fork. The old script has been replaced with local snapshot/verify commands.

# Release process: a curated public history

The engineering record for how zot-web ships publicly. The flow is
terva's release-cut, trimmed (terva's `docs/plans/release-process.md`
carries the deep rationale — staging-gate philosophy, curation
doctrine, the scrub model); this file records where zot-web differs
and everything needed to run a cut without leaving this repo. The
tooling is `scripts/release.sh` behind the `release-*` just targets
(`release.just`, an optional import the public justfile works
without).

## Why this shape

- **main's history is private.** Day-to-day commits and CI reference
  internal hosts. The public history is a separate curated narrative
  on the `release` branch — main SHAs never appear publicly.
- **There is nothing to version publicly.** No public tags, binaries,
  or release objects: the curated feat/fix history IS the product.
  Internal binaries still build on the Forgejo (goreleaser, `v*`
  tags) — that machinery never ships.
- **The name stays zot-web.** Wire compatibility with zot's extension
  protocol is the point; rebranding would work against it. The public
  home is github.com/terva-sh/zot-web and the module path matches —
  code identity = public identity, renamed once on main rather than
  rewritten at cut time (Go import paths can't be marker-rewritten).

## The flow

```
just release-cut        curation worktree at ../zot-web-release-cut
(author curated commits in the worktree)
just release-verify     byte-identity, scrub, gofmt/vet/race
just release-publish    push release → origin (backup) + staging gate
(test in the gate, then go live from inside it)
git -C <gate> push origin refs/heads/release
```

Curation conventions are terva's: one theme per commit, hand-written
lowercase `feat:`/`fix:` messages, never an internal SHA, hostname,
or branch name. The staged diff must be partitioned, not changed —
verify enforces byte-identity against the prepared candidate.

## Versioning

There are no *public* versions (no tags, release objects, or binaries —
see above), but the extension artifact still carries one: terva/zot
require a `version` in `extension.json` and surface it in `ext list`,
and the same string rides the protocol hello, the default User-Agent,
and `--version` (`internal/version.Version`, the committed default a
source build reports). Keep them in lockstep and bump when a cut changes
behavior:

- **Before `release-cut`**, set both `extension.json` `"version"` and
  `internal/version.Version` to the same value — minor for features,
  patch for fixes.
- `TestManifestVersionMatchesCode` pins the two equal, so a drift fails
  `release-verify`'s race-test pass (and ordinary `go test`). They had
  silently disagreed (`0.1.0` vs `0.2.0`) before this guard.

(Internal goreleaser builds override `Version` from the `v*` git tag via
ldflags; that path is internal-only and unrelated to the public cut.)

## Cut ranges without versions

Each publish drops a lightweight `cut/N` tag (N increments) on the
main commit it cut from; the next cut's worklist is `cut/N..HEAD`.
The markers are internal — pushed to origin, never the mirror. The
first cut had no marker and staged the whole tree on an **orphan
root**: zot-web is original work with no upstream history to
continue, so the public history simply starts at the curated story.

## What never ships

`EXCLUDES` in scripts/release.sh removes: `.forgejo/` (internal
runners and action mirrors), `.goreleaser.yaml` (internal release
conventions), `.claude/` (working notes/skills), `docs/plans/` (this
file — references internal infra), `scripts/release.sh` +
`release.just` (embed the blocklist). The `BLOCKLIST` (same strings
as terva's) fails verify on a hit anywhere in the tree.

## Staging gate

A local clone of github.com/terva-sh/zot-web; publish lands `release`
there, going live is an explicit push from inside it. Probed at
`~/Workspace/github.com/terva-sh/zot-web` (macOS) or
`~/workspace/github.com/terva-sh/zot-web` (Linux);
`ZOT_WEB_MIRROR_DIR` overrides. One-time setup per machine:

```
git clone git@github.com:terva-sh/zot-web.git ~/workspace/github.com/terva-sh/zot-web
git -C ~/workspace/github.com/terva-sh/zot-web config receive.denyCurrentBranch updateInstead
```

Test from the gate before pushing (`go test ./...`, run the binary).
A staged cut that fails gets fixed on main and re-cut — with no
public versions there's nothing to skip or renumber; the branch just
moves when it's actually ready. Publish warns if the `mirror` remote
points at real GitHub (gateless: every publish goes live instantly).

## Recovery

`just release-status` shows the in-progress cut; `just release-abort`
tears it down (restores or deletes the local release branch — nothing
external exists before the gate push). The gate can always be
hard-reset to `origin/release`; the source of truth for an unpushed
cut is re-cuttable from main at any time.

`release-publish` is safe to re-run. It preflights the gate working
tree (a dirty gate fails *before* any push, since
`denyCurrentBranch=updateInstead` would reject it and strand a
half-done publish) and reuses the existing `cut/N` marker rather than
minting a fresh one — so if a push dies partway, fix the cause and run
it again to resume.

## Cut record

Each entry is a curated cut and its go-live. The `cut/N` markers are
internal tags on the main commit each was cut from; there are no public
versions. The public mirror's only branch is `release`, so it is the
default branch by construction.

- **cut/1** — 2026-06-12, orphan root, four commits: the extension
  skeleton (proto/config/manifest/run.sh), hardened fetch, SearXNG
  search, vendored deps. Live.
- **cut/2** — 2026-06-13: terva host integration (read-only /
  ask-before-write hints, bundled research skill), loopback-by-default
  SSRF allowlist, save-path preflight + fetch-hint polish, and config
  reading from `data_dir` with an install-dir fallback. Live.
- **cut/3** — 2026-06-14: the MIT license. Live.
- **cut/4** — 2026-06-19: terva extension protocol v2 (host detection
  via `terva_version`, network-read authority on every tool, the
  `session_start` subscription with live-cwd saves, no `min_protocol`)
  and the zot/terva conformance harness. Live.
- **cut/5** — 2026-06-19: version bumped to 0.3.0 (manifest = code,
  pinned equal by a test). Live.
- **cut/6** — 2026-06-21: flatten untrusted page-derived values (page
  `<title>`, search result title, `<img>` width/height) out of the
  tools' own result scaffolding so a newline can't forge a metadata,
  result, or listing line; version 0.3.1. (The decision to *not* add an
  inline untrusted-content marker is recorded in
  `untrusted-web-content.md` — internal, never shipped.) Live.

Each cut was verified (byte-identity to the candidate, scrub, race
tests) and gate-tested (build, `go test ./...`, `--version` →
zot-web 0.3.1) before go-live.
