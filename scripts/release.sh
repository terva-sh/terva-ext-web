#!/usr/bin/env bash
#
# The release-cut flow for zot-web: curate main into the public
# `release` branch and stage it for github.com/terva-sh/zot-web.
# Modeled on terva's flow (the engineering record lives there:
# terva docs/plans/release-process.md), minus everything zot-web
# doesn't need — there are NO public versions, tags, or binary
# releases. The public history IS the product: day-to-day commits
# grouped into a curated feat/fix narrative.
#
# Verbs:
#   cut          prepare a curation worktree (first cut: orphan root —
#                the public history starts at the curated commits)
#   verify       prove the curated tree is publishable (identity,
#                scrub, toolchain)
#   publish      push `release` to origin (backup) and the staging
#                gate; drop the internal cut/N range marker
#   status       show the in-progress cut's state
#   worklist     print the curation worklist
#   abort        tear down an in-progress cut, restore branches
#
# The staging gate is a local clone of github.com/terva-sh/zot-web
# (probe below; ZOT_WEB_MIRROR_DIR overrides). publish lands the
# branch there; going live is an explicit push from inside the clone.
#
# This file is EXCLUDED from the public tree: it embeds the internal
# scrub blocklist.

set -euo pipefail

# ---- fixed policy: what never ships publicly ----

EXCLUDES=(
  ".forgejo"
  ".goreleaser.yaml"
  ".claude"
  "docs/plans"
  "scripts/release.sh"
  "release.just"
)

BLOCKLIST=(
  "local.sothr.com"
  "container.local"
  "warricksothr"
  "Sothr-Mirrors"
  "ssh://git@"
)

# Per-machine: override with ZOT_WEB_MIRROR_DIR, else probe the
# conventional clone locations (macOS Workspace/, Linux workspace/).
MIRROR_URL_DEFAULT="${ZOT_WEB_MIRROR_DIR:-}"
if [ -z "$MIRROR_URL_DEFAULT" ]; then
  for _cand in "$HOME/Workspace/github.com/terva-sh/zot-web" \
               "$HOME/workspace/github.com/terva-sh/zot-web"; do
    if [ -d "$_cand" ]; then MIRROR_URL_DEFAULT="$_cand"; break; fi
  done
fi
: "${MIRROR_URL_DEFAULT:=$HOME/workspace/github.com/terva-sh/zot-web}"

MAIN_BRANCH="main"

# ---- shared plumbing ----

ROOT=$(git rev-parse --show-toplevel)
GIT_COMMON=$(git rev-parse --path-format=absolute --git-common-dir)
STATE_DIR="$GIT_COMMON/zotweb-cut"
STATE_FILE="$STATE_DIR/state"
WORKLIST="$STATE_DIR/worklist.md"
WT="$(dirname "$ROOT")/zot-web-release-cut"

msg()  { printf '\033[1m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[33mwarn:\033[0m %s\n' "$*" >&2; }
die()  { printf '\033[31merror:\033[0m %s\n' "$*" >&2; exit 1; }

state_get() { sed -n "s/^$1=//p" "$STATE_FILE"; }
state_set() {
  local key=$1 val=$2 tmp
  tmp=$(mktemp)
  { [ -f "$STATE_FILE" ] && grep -v "^$key=" "$STATE_FILE" || true; } >"$tmp"
  printf '%s=%s\n' "$key" "$val" >>"$tmp"
  mv "$tmp" "$STATE_FILE"
}

require_state() { [ -f "$STATE_FILE" ] || die "no cut in progress (run: just release-cut)"; }

require_clean_main() {
  [ -z "$(git -C "$ROOT" status --porcelain --untracked-files=no)" ] \
    || die "working tree has uncommitted changes; commit or stash first"
}

worktree_tree() { git -C "$WT" rev-parse 'HEAD^{tree}'; }

fingerprint_worktree() {
  # Stage everything into a scratch index and hash the tree — the
  # curator's commits must reproduce exactly this tree.
  local wt_gitdir idx
  wt_gitdir=$(git -C "$WT" rev-parse --absolute-git-dir)
  idx="$wt_gitdir/cut-index"
  rm -f "$idx"
  GIT_INDEX_FILE="$idx" git -C "$WT" add -A
  GIT_INDEX_FILE="$idx" git -C "$WT" write-tree
  rm -f "$idx"
}

# Newest cut/N marker that is an ancestor of the given commit — the
# worklist range start. Numeric N, no versions anywhere.
prev_cut_marker() {
  local sha=$1 t
  for t in $(git -C "$ROOT" tag -l 'cut/[0-9]*' | sort -t/ -k2 -rn); do
    if git -C "$ROOT" merge-base --is-ancestor "$t" "$sha" 2>/dev/null; then
      printf '%s' "$t"
      return 0
    fi
  done
  return 1
}

next_cut_marker() {
  local n
  n=$(git -C "$ROOT" tag -l 'cut/[0-9]*' | sed 's#cut/##' | sort -rn | head -1)
  printf 'cut/%d' "$(( ${n:-0} + 1 ))"
}

# The cut/N marker (if any) pointing exactly at the given commit. publish
# creates the marker before the gate push, so a partial publish leaves it
# behind; reusing it on a re-run keeps the number stable instead of minting a
# spurious cut/N+1.
marker_on() {
  local sha=$1 t
  for t in $(git -C "$ROOT" tag -l 'cut/[0-9]*'); do
    if [ "$(git -C "$ROOT" rev-parse "$t^{commit}" 2>/dev/null)" = "$sha" ]; then
      printf '%s' "$t"
      return 0
    fi
  done
  return 1
}

# Where the gate push will land: an existing 'mirror' remote wins, else the
# probed default (the same path ensure_mirror_remote would add).
mirror_dir() {
  git -C "$ROOT" remote get-url mirror 2>/dev/null || printf '%s' "$MIRROR_URL_DEFAULT"
}

# Fail before any side effect if the staging gate has a dirty working tree:
# receive.denyCurrentBranch=updateInstead rejects the push otherwise, which
# would strand a half-done publish (marker + origin pushed, gate not). No-op
# when the mirror is a real remote — the gateless case has no working tree.
require_clean_gate() {
  local gate
  gate=$(mirror_dir)
  git -C "$gate" rev-parse --is-inside-work-tree >/dev/null 2>&1 || return 0
  [ -z "$(git -C "$gate" status --porcelain)" ] \
    || die "staging gate at $gate has uncommitted changes — clean it (git -C \"$gate\" status) and re-run publish"
}

ensure_mirror_remote() {
  if ! git -C "$ROOT" remote get-url mirror >/dev/null 2>&1; then
    git -C "$ROOT" remote add mirror "$MIRROR_URL_DEFAULT"
  fi
  msg "mirror -> $(git -C "$ROOT" remote get-url mirror)"
}

# ---- verbs ----

cmd_cut() {
  [ ! -d "$STATE_DIR" ] || die "a cut is already in progress (just release-status / release-abort)"
  [ ! -e "$WT" ] || die "stale worktree at $WT — remove it or run release-abort"
  require_clean_main
  [ "$(git -C "$ROOT" rev-parse --abbrev-ref HEAD)" = "$MAIN_BRANCH" ] \
    || die "cut from $MAIN_BRANCH (currently on $(git -C "$ROOT" rev-parse --abbrev-ref HEAD))"

  git -C "$ROOT" fetch origin --quiet || warn "could not fetch origin; cutting from the local $MAIN_BRANCH"
  if [ "$(git -C "$ROOT" rev-parse "$MAIN_BRANCH")" != "$(git -C "$ROOT" rev-parse "origin/$MAIN_BRANCH" 2>/dev/null || echo unknown)" ]; then
    warn "$MAIN_BRANCH differs from origin/$MAIN_BRANCH — cutting the LOCAL branch"
  fi

  local cut_sha prev=""
  cut_sha=$(git -C "$ROOT" rev-parse "$MAIN_BRANCH")
  prev=$(prev_cut_marker "$cut_sha" || true)

  mkdir -p "$STATE_DIR"
  local release_prev=none
  if git -C "$ROOT" show-ref --verify -q refs/heads/release; then
    release_prev=$(git -C "$ROOT" rev-parse refs/heads/release)
    git -C "$ROOT" worktree add "$WT" release
  elif git -C "$ROOT" show-ref --verify -q refs/remotes/origin/release; then
    # Second machine: continue the published history, never re-root.
    msg "creating local release branch from origin/release"
    git -C "$ROOT" worktree add -b release "$WT" origin/release
  else
    # First cut: the public history has no upstream to continue — it
    # starts at the curated commits, on an orphan root.
    git -C "$ROOT" worktree add --detach "$WT" "$MAIN_BRANCH"
    git -C "$WT" switch --orphan release
  fi

  msg "building the candidate public tree"
  find "$WT" -mindepth 1 -maxdepth 1 ! -name .git -exec rm -rf {} +
  git -C "$ROOT" archive --format=tar "$cut_sha" | tar -x -C "$WT"
  local p
  for p in "${EXCLUDES[@]}"; do rm -rf "${WT:?}/$p"; done

  local candidate_tree
  candidate_tree=$(fingerprint_worktree)

  state_set cut_sha "$cut_sha"
  state_set prev "${prev:-}"
  state_set release_prev_sha "$release_prev"
  state_set candidate_tree "$candidate_tree"

  {
    echo "# Curation worklist — zot-web"
    echo
    echo "Worktree: $WT (branch: release)"
    echo "This file lives in the git dir, never in the tree — it cannot ship."
    echo
    echo "Stage one theme at a time, commit with a hand-written message"
    echo "(feat:/fix: style; NEVER reference internal SHAs, hosts, or"
    echo "branches). Repeat until 'git status' is clean, then"
    echo "'just release-verify'."
    echo
    if [ -n "${prev:-}" ]; then
      echo "Work landed since $prev ($(git -C "$ROOT" rev-parse --short "$prev")..$(git -C "$ROOT" rev-parse --short "$cut_sha")):"
    else
      echo "First cut — the staged diff is the WHOLE tree. Tell the"
      echo "project's story in a handful of commits. Internal history"
      echo "for reference:"
    fi
    echo
    git -C "$ROOT" log --format='- [ ] %h %s' --reverse ${prev:+$prev..}"$cut_sha"
  } >"$WORKLIST"

  msg "candidate ready: $WT"
  msg "worklist:        $WORKLIST"
  msg "next: author the curated commits in the worktree, then 'just release-verify'"
}

cmd_verify() {
  require_state
  [ -d "$WT" ] || die "worktree missing at $WT"

  msg "verify: curation state"
  [ -z "$(git -C "$WT" status --porcelain)" ] \
    || die "worktree not clean — commit or drop the remaining changes ($WT)"
  git -C "$WT" rev-parse -q --verify HEAD >/dev/null \
    || die "no commits on release yet — author the curated commits first"

  local candidate_tree head_tree
  candidate_tree=$(state_get candidate_tree)
  head_tree=$(worktree_tree)
  [ "$head_tree" = "$candidate_tree" ] \
    || die "released tree ($head_tree) differs from the prepared candidate ($candidate_tree) — curation must partition the diff, not change it"

  msg "verify: scrub (blocklist + exclusions)"
  local pat p
  for pat in "${BLOCKLIST[@]}"; do
    if git -C "$WT" grep -nIF "$pat" -- . >/dev/null 2>&1; then
      git -C "$WT" grep -nIF "$pat" -- . | head -20 >&2
      die "blocklisted string '$pat' in the public tree"
    fi
  done
  for p in "${EXCLUDES[@]}"; do
    [ ! -e "$WT/$p" ] || die "excluded path present in the public tree: $p"
  done

  msg "verify: toolchain (gofmt, vet, race tests)"
  local fmt
  fmt=$(cd "$WT" && gofmt -l . | grep -v '^vendor/' || true)
  [ -z "$fmt" ] || die "gofmt issues in the public tree: $fmt"
  (cd "$WT" && go vet ./...)
  (cd "$WT" && go test -race ./...)

  state_set verified_tree "$(worktree_tree)"
  msg "verified — 'just release-publish' when ready"
}

cmd_publish() {
  require_state
  local verified_tree cut_sha
  verified_tree=$(state_get verified_tree)
  cut_sha=$(state_get cut_sha)
  [ -n "$verified_tree" ] || die "not verified — run 'just release-verify' first"
  [ -d "$WT" ] || die "worktree missing at $WT"
  [ -z "$(git -C "$WT" status --porcelain)" ] || die "worktree not clean — re-verify"
  [ "$(worktree_tree)" = "$verified_tree" ] || die "tree changed since verify — run 'just release-verify' again"

  # Fail fast before any side effect if the gate can't receive the push. With
  # the marker reused (not re-minted) below, a publish that dies past this point
  # is safe to re-run once the cause is fixed.
  require_clean_gate

  local marker
  marker=$(marker_on "$cut_sha" || next_cut_marker)
  git -C "$ROOT" rev-parse -q --verify "refs/tags/$marker" >/dev/null 2>&1 \
    || git -C "$ROOT" tag "$marker" "$cut_sha"

  msg "pushing the internal channel (origin: release backup + $marker range marker)"
  git -C "$ROOT" push origin "refs/heads/release:refs/heads/release" "refs/tags/$marker"

  msg "pushing the staging gate (mirror)"
  ensure_mirror_remote
  git -C "$ROOT" push mirror "refs/heads/release:refs/heads/release"

  msg "cleaning up"
  git -C "$ROOT" worktree remove "$WT"
  rm -rf "$STATE_DIR"

  local mirror_url
  mirror_url=$(git -C "$ROOT" remote get-url mirror)
  case "$mirror_url" in
  git@github.com:* | https://github.com/* | ssh://git@github.com/*)
    warn "mirror is the REAL GitHub remote — this publish went live with NO staging gate"
    ;;
  *)
    msg "staged in the gate at $mirror_url — inspect and test it, then GO LIVE from inside it:"
    msg "  git -C $mirror_url push origin refs/heads/release"
    ;;
  esac
}

cmd_status() {
  [ -f "$STATE_FILE" ] || { msg "no cut in progress"; return 0; }
  cat "$STATE_FILE"
  if [ -d "$WT" ]; then
    echo
    git -C "$WT" log --oneline -5 2>/dev/null || echo "(no commits on release yet)"
    echo
    git -C "$WT" status --short | head -20
  fi
}

cmd_worklist() {
  require_state
  cat "$WORKLIST"
}

cmd_abort() {
  require_state
  [ ! -e "$WT" ] || git -C "$ROOT" worktree remove --force "$WT"
  local prev
  prev=$(state_get release_prev_sha)
  if [ "$prev" = "none" ]; then
    git -C "$ROOT" branch -D release 2>/dev/null || true
  else
    git -C "$ROOT" branch -f release "$prev"
  fi
  rm -rf "$STATE_DIR"
  msg "aborted — branches restored, worktree and state removed"
}

usage() { sed -n '2,27p' "$0" | sed 's/^# \{0,1\}//'; }

cmd="${1:-}"
shift || true
case "$cmd" in
cut) cmd_cut ;;
verify) cmd_verify ;;
publish) cmd_publish ;;
status) cmd_status ;;
worklist) cmd_worklist ;;
abort) cmd_abort ;;
-h | --help | help | "") usage ;;
*) die "unknown verb: $cmd (cut | verify | publish | status | worklist | abort)" ;;
esac
