# terva-ext-web

Terva web tools, forked from zot-web with its history intact. Stock zot is not
supported. Machine-wide rules and `/home/sothr/workspace/AGENTS.md` also apply.

## Start here

1. Inspect `git status`, the branch and `git remote -v` before editing.
2. Read `docs/plans/terva-ext-web.md` for direction and
   `docs/plans/modernization-critical-path.md` for remaining release work.
3. Read `README.md` for current behavior, setup and configuration.
4. Read `go.mod`, `extension.json`, `main.go`, `runtime.go`, `conformance_test.go`,
   `justfile`, `run.sh` and `.goreleaser.yaml` before changing integration.
5. Check actual contracts in the pinned `terva.sh/terva v0.137.0` SDK and host
   source. `docs/plans/sdk-verification.md` records the verified version.
   Inspect any sibling checkout's status and revision before using it. Do not
   copy another extension's SDK pin without verifying the published version.

## Implementation rules

Use a topic branch. Repository, module and binary are `terva-ext-web`.
Manifest and handshake name are `web`. Preserve all six `web_*` tools and
`/web-cache`, and include the research skill in release archives.

Configuration comes from Terva host settings, `TERVA_EXT_WEB_*` overrides and
`TAVILY_API_KEY`. Do not restore standalone configuration files, `ZOT_WEB_*`
aliases or `configuration_source`. Their removal was user-approved; see
`docs/plans/config-provenance.md`. Existing installation and credential files
must remain untouched by source changes.

Preserve SSRF, resource-limit, path and output-sanitization checks. Each save
captures one workspace before fetching. Keep the session-switch regression.
The SDK has no per-call cancellation context or dedicated result trust metadata;
do not claim either. Keep SDK, configuration and release changes reviewable.

Run checks appropriate to the change through `just`. `just ci` covers vet,
formatting, race tests, subprocess conformance, the published host driver and
vendor consistency. Recheck the toolchain and establish a baseline before code
changes. Native validation covers five targets on GitHub; its reports certify
specific commits and archives. Source, dependency, manifest, launcher or
packaging changes need new native evidence before release.

## Tickets

Track project work in `.tickets/`. The workspace ledger covers infrastructure
and clone setup. Read the generated workflow below before writing tickets.

- At session start, inspect `git ticket ready`, draft tickets and in-progress
  tickets. Read the selected ticket and check `git ticket files PATH` for prior
  decisions about files you will change.
- Record the approach, alternatives, progress and validation in the ticket.
  Commit ticket changes as work proceeds. Include the ID and title in commits
  primarily concerning one ticket.
- Follow `docs/ticket-labels.md` and the vocabulary in `.tickets/config.yml`.
  Status, type, priority and dependencies belong in native fields.
- Before handoff, run `git ticket check` and `just ticket-check`, commit intended
  changes and leave unrelated work untouched.
- Regenerate the block with `git ticket instructions --write`; do not hand-edit
  it. Keep project-specific instructions above it.
- After cloning, run `git ticket install-merge-driver`. The committed
  `.gitattributes` does not install local Git configuration.

`just ticket-check` requires git-ticket locally. Go CI does not provision it.

## Repository and release boundaries

The fork point is `c50d773c60250a7315c37c2aedbeb891cb12f6a4`, tree
`6c368f89b3e622c7d4515ac4e02d47caed3a95c7`. Preserve inherited history and
leave the old `cut/*` tags as historical records.

Origin is `ssh://git@git.local.sothr.com:2222/terva-sh/terva-ext-web.git`.
The GitHub mirror is `git@github.com:terva-sh/terva-ext-web.git`.
Never push this fork to zot-web. Binary publication remains disabled pending
installation, migration and rollback validation and release approval.

Do not archive zot-web without ticket authorization. Do not enable old and new
installations together: they share tool, command, configuration and secret
identities. Resolve paths through the host rather than repository basenames.

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
