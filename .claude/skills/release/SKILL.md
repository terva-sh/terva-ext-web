---
name: release
description: Cut and publish zot-web's curated public history to github.com/terva-sh/zot-web. Use when asked to cut a release, publish zot-web, or push to the public mirror. There are no versions, tags, or binaries publicly — the curated branch is the release.
---

# Releasing zot-web

You are driving the flow recorded in `docs/plans/release-process.md`
(THIS repo's engineering record — read it if anything here
surprises you; it also records where zot-web deliberately differs
from terva's flow, which it is modeled on). No versions, tags, or
binaries ship publicly: the `release` branch's curated feat/fix
history IS the product. Tooling: `scripts/release.sh` behind
`just release-*`.

## Flow

1. Preflight: clean tree, on `main`, `just ci` green.
2. `just release-cut` — curation worktree at `../zot-web-release-cut`
   plus a worklist (`just release-worklist`). The first cut stages the
   WHOLE tree on an orphan root: tell the project's story in a handful
   of commits.
3. Curate in the worktree: stage one theme at a time, hand-written
   lowercase `feat:`/`fix:` messages, never referencing internal SHAs,
   hostnames, or branches. The staged diff must be partitioned, not
   changed — verify enforces byte-identity.
4. `just release-verify` — identity, blocklist scrub, gofmt/vet/race.
5. `just release-publish` — pushes `release` to origin (backup) and
   the staging gate (a local clone of github.com/terva-sh/zot-web,
   probed under `~/workspace/github.com/terva-sh/`;
   `ZOT_WEB_MIRROR_DIR` overrides), and drops the internal `cut/N`
   range marker. If publish warns the mirror is the real GitHub
   remote, the gate is missing on this machine — set it up like
   terva's (clone + `receive.denyCurrentBranch updateInstead`).
   Publish is re-runnable: a dirty gate fails fast before any push,
   and the `cut/N` marker is reused (not re-minted), so clean the gate
   and run it again to resume a partial publish.
6. **Test from the gate** (`go test ./...`, run the extension), then
   go live from inside it — publish prints the exact command. Ask the
   maintainer before the go-live push unless they already said to.
