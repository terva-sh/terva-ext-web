# terva-ext-web

This is an independent local fork of zot-web. The target is the Terva
extension ecosystem; stock zot compatibility is no longer a product goal.
The machine-wide rules and `/home/sothr/workspace/AGENTS.md` also apply.

## Start here

1. Inspect `git status`, the branch, and `git remote -v` before editing.
2. Read `docs/plans/terva-ext-web.md`. It records the exact fork point,
   inspected Terva revision, decisions, feature assessment, migration risks,
   and implementation batches. It is the current direction.
3. Read `README.md` for existing tool behavior, but treat its permanent zot
   compatibility and no-rename statements as historical policy superseded by
   the split plan. No implementation rename has happened yet.
4. Read `go.mod`, `extension.json`, `main.go`, `internal/proto/proto.go`,
   `conformance_test.go`, `justfile`, `run.sh`, and `.goreleaser.yaml`.
5. Consult the sibling Terva checkout's `docs/extensions.md`,
   `packages/agent/ext/ext.go`, and `packages/agent/extproto/extproto.go` for
   actual SDK contracts. It is at `../terva`; inspect its status and revision
   first and leave its unrelated work untouched. `../terva-ext-caldav` is a
   useful naming/configuration example. Do not blindly copy the older SDK
   pin from `../terva-extension-template`.

## First implementation batch

Start a topic branch and implement **identity and packaging** from the plan.
Use repository/binary/module name `terva-ext-web`; preserve manifest and
handshake name `web`, all six `web_*` tools, and `/web-cache`. Update internal
imports and linker symbols together. Include the bundled research skill in
release archives. Resolve the inherited release-owner discrepancy before any
publication. Preserve existing configuration; do not rename or delete users'
installed extensions or credential files as part of source changes.

Keep the SDK migration as a separate reviewable batch after identity and
packaging pass their checks. The current code still uses handwritten
`internal/proto`; SDK migration was not done before this fork. Carry forward
SSRF, resource-limit, path, and output-sanitization defenses. Fix the
double-read of cwd around downloads with a meaningful session-switch test.
Do not promise per-call cancellation or trust metadata that the SDK does not
provide. Verify a published SDK version before pinning and vendoring it.

Use the repository's `just` recipes and CI definitions for checks. Establish
the unchanged-source baseline before code changes. The setup session could
not run Go tests because `go` was absent from PATH; this is a validation gap,
not a passed baseline. Check the environment again in the new session.
Commit coherent changes and record implementation rationale in source control.
There is no project ticket store yet. If one is established, read its generated
workflow before writing to it; the workspace ledger tracks clone setup only.

## Repository boundaries

- Exact fork point: `c50d773c60250a7315c37c2aedbeb891cb12f6a4`;
  tree: `6c368f89b3e622c7d4515ac4e02d47caed3a95c7`.
- The initial docs commit follows that point. Do not rewrite the inherited
  history or reuse the inherited `cut/*` tags for new releases.
- No remote exists yet. The source repository is `../zot-web`; do not set it
  as a push destination. Remote publication remains a later batch.
- Do not archive zot-web during modernization. Validate the replacement
  release and migration first; remote archival needs ticket authorization.
- Do not enable both old and new installations together: they share `web`
  identity, tool names, command names, configuration scope, and secret scope.
