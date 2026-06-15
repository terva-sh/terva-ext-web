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

## First-cut record

2026-06-12: orphan root, four commits — the extension skeleton
(proto/config/manifest/run.sh), hardened fetch, SearXNG search,
vendored deps. Verified (scrub + race tests), gate-tested (build,
tests, `--version` → zot-web 0.2.0), staged. Go-live: _pending — fill
in date + confirm the default branch flipped to `release`._
