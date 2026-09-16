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
   the split plan. The identity and packaging batch is recorded in
   `docs/plans/identity-packaging.md`.
4. Read `go.mod`, `extension.json`, `main.go`, `internal/proto/proto.go`,
   `conformance_test.go`, `justfile`, `run.sh`, and `.goreleaser.yaml`.
5. Consult the sibling Terva checkout's `docs/extensions.md`,
   `packages/agent/ext/ext.go`, and `packages/agent/extproto/extproto.go` for
   actual SDK contracts. It is at `../terva`; inspect its status and revision
   first and leave its unrelated work untouched. `../terva-ext-caldav` is a
   useful naming/configuration example. Do not blindly copy the older SDK
   pin from `../terva-extension-template`.

## First implementation batch

The **identity and packaging** batch is implemented; see its decision record.
Continue on a topic branch for subsequent batches.
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
## Project ticket workflow

Track codebase work in this repository's `.tickets/` store. The workspace
ledger is for workspace infrastructure and clone setup, not this project's
implementation. Read the generated Tickets block below before writing tickets;
run `git ticket instructions` for the full rationale.

- At session start, inspect `git ticket ready`, `git ticket list --status draft`,
  and `git ticket list --status in-progress`. Read the selected ticket before
  editing, and use `git ticket files PATH` to find recorded work on a file.
- Keep `docs/plans/terva-ext-web.md` as the modernization direction. Link the
  relevant plan and source paths from implementation tickets; record decisions,
  alternatives, progress, and validation in the ticket as the work happens.
  Keep SDK/correctness, configuration/credentials, and release work separate.
- Commit ticket changes as you go. Include the ticket ID and title in commits
  primarily concerning one ticket. Before handing off or ending a session,
  run `git ticket check` (or the stricter `just ticket-check`), commit all
  intended ticket changes, and leave unrelated changes untouched.
- Regenerate the block below with `git ticket instructions --write`; do not
  hand-edit it. Keep repository-specific additions in this section so they
  survive regeneration.
- After cloning, run `git ticket install-merge-driver` to configure the local
  driver used by the committed `.gitattributes`. Git configuration is local to
  each clone and is not installed by checking out the attributes file.

`just ticket-check` requires the installed `git-ticket` tool. It is a local
handoff gate; the Go CI pipeline does not yet provision or run git-ticket.

## Repository boundaries

- Exact fork point: `c50d773c60250a7315c37c2aedbeb891cb12f6a4`;
  tree: `6c368f89b3e622c7d4515ac4e02d47caed3a95c7`.
- The initial docs commit follows that point. Do not rewrite the inherited
  history or reuse the inherited `cut/*` tags for new releases.
- Origin is `ssh://git@git.local.sothr.com:2222/terva-sh/terva-ext-web.git`.
  The source repository is `../zot-web`; never use it as a push destination.
  Binary release publication remains a later batch.
- Do not archive zot-web during modernization. Validate the replacement
  release and migration first; remote archival needs ticket authorization.
- Do not enable both old and new installations together: they share `web`
  identity, tool names, command names, configuration scope, and secret scope.

<!-- git-ticket:begin -->

## Tickets

Work is tracked as Markdown tickets in `.tickets/`, managed with `git ticket`.
`git ticket help` lists every command.

`git ticket instructions` prints the long form of this block, carrying the same
rules with the reason for each one and the failure it prevents. Read it when a
rule here surprises you, when you are about to work around one, or when you want
the parts this summary leaves out.

Name yourself on every command that writes, as `agent:tool/session`. Without it
the write is attributed to the first actor in `config.yml`, usually a person, and
your claim then tells other agents that a human is holding the ticket.

```sh
git ticket note TKT-01ABCD "..." --actor agent:yourtool/session-3
```

### Finding work

`git ticket ready` is the queue: open, unblocked, every dependency closed.

`git ticket list --status draft` is the rest of the backlog and is usually the
larger half, because everything filed lands in `draft` and stays there until a
person promotes it. Read both before reporting that there is nothing to pick up.

Do not promote a draft yourself. Name the ones that look startable, say what
makes each startable, and let the person you are working with choose.

`git ticket show ID` reads one ticket. It prints the newest note in full and
replaces older ones with a line naming how many there are; `--list` and
`--show N` on the `note` command print those.

`git ticket search QUERY` matches a substring of the title, the body, and the
references. Add `--regex` to search by pattern instead.

Any unique ID prefix works, down to four characters. Copy one from a listing
rather than shortening it yourself, because a ULID opens with a timestamp and
tickets filed together are identical that far in.

### Doing the work

A draft cannot be claimed. If you were asked to pick up something still in
`draft`, that request is the promotion, so run `git ticket status ID ready`
first. A ticket already on the queue needs no such step.

Then `git ticket claim ID` and `git ticket status ID in-progress`.

Read the code before you plan, then write the approach with
`git ticket plan ID "..."`. It replaces rather than appends, so revising it
leaves one plan instead of a stack.

While you work, `git ticket note ID "..."` records what the next person will
need and does not have, and `git ticket ac ID --check N` ticks an acceptance
criterion, counting checkbox lines from one.

Leave unticked any criterion you could not satisfy, and say in a note what
stopped you. Nothing reports an empty box, so an honest one costs nothing, while
a tick you did not earn costs the next reader their trust in every other box.

Finish with `git ticket summary ID "..."` saying where it landed, then
`git ticket status ID done` and `git ticket release ID`.

If you cannot proceed, `git ticket status ID blocked --reason "..."`. The reason
is required.

`note` appends. `plan` and `summary` replace, so correct a note by adding another
that says which one it supersedes.

### Filing new work

`git ticket create --title "..." --type bug --priority high` files a ticket.
Types are task, bug, chore, spike, and epic. Add `--parent` to file it under an
epic.

Run `git ticket config` before you invent a label. It prints what this store
permits, and a label outside that set is a warning that fails
`check --strict`. `git ticket schema` prints the types, priorities, statuses and
error codes every store shares.

Write prose longer than a line to a file and pass the file: `--description-file`
on `create` and `update`, and `--file` on `plan`, `note`, `comment`, and
`summary`. A path of `-` reads stdin.

Inside that prose write subheadings as `###`, because a line opening with `## `
starts a new section and everything below it lands somewhere you did not intend.
The exception runs the other way. If the subheading names a section the format
owns, `### Acceptance criteria` is prose no command can reach. Use `--ac` and
`--dod`, or write `## `, which opens the real section. `check` reports the
mistake as `section_heading_demoted`.

A ticket you file lands in `draft` and stays there. File it, say that you filed
it, and go back to what you were doing.

Record structure rather than describing it in prose:
`git ticket link ID --depends-on OTHER` says this waits on that, and
`git ticket deps ID` walks the chain.

Keep the title under 72 characters. Over that `check` warns, and over 120 the
write is refused.

When you mention a ticket in prose, put its title beside the ID the first time,
then the bare ID is enough. A ULID tells a reader nothing on its own.

### When check fails

`git ticket check --fix --dry-run --strict` plans every repair, prints what it
would do, and writes nothing. Run `git ticket check --fix` and commit what it
changed.

If this project runs the check in CI, it reports the repair and does not commit
it for you.

### Driving it from a script

`--json` on any command gives one envelope on stdout with a stable error `code`
to switch on. Read the vocabulary from `git ticket schema` rather than
hard-coding it.

A write answers with `mutation-result`, whose `ticket` is an `{id, revision}`
stub rather than the ticket, so read the body back with `show --json` when you
need it. Body sections come back camelCase, as `implementationPlan` and
`acceptanceCriteria`.

Every write takes `--if-revision R` and refuses if the ticket moved since you
read it. Pass it whenever you read, decide, and then write.

Text that opens with a dash goes after a bare `--`.

<!-- git-ticket:end -->
