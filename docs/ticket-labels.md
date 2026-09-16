# Ticket labels

Labels are namespaced strings stored in each ticket's `labels` list. The
allowed vocabulary lives in `.tickets/config.yml`; `git ticket config` prints
it. Extend that list deliberately when a new grouping is needed. An undeclared
label produces `label_unknown` and fails `just ticket-check`.

## Grouping conventions

| Namespace | Values | Meaning |
| --- | --- | --- |
| `initiative:` | `terva-modernization` | The work described by the split plan, including both epics and all their children. Useful when aggregating exports across repositories. |
| `scope:` | `core`, `follow-up` | Exactly one per modernization ticket: core migration/release work, or subsequent features and investigations. |
| `area:` | See below | One or more directly relevant work areas. Shared dependencies alone do not make a ticket belong to an area. |

Scope is a grouping, not an authorization or readiness signal. Use native
`status`, `type`, `priority`, parent and dependency fields for those concepts.
Do not add labels that duplicate them, such as `draft`, `high`, or `spike`.
Future modernization children need their own labels; labels are not inherited
from the parent epic. Epics carry their broad primary areas.

| Area label | Covers |
| --- | --- |
| `area:sdk` | SDK selection, protocol integration, and host contracts |
| `area:sessions` | Session/workspace identity, transitions, and isolation |
| `area:security` | Authority, permissions, credentials, integrity, and privacy |
| `area:config` | Settings, precedence, backend readiness, and config secrets |
| `area:validation` | Dedicated compatibility, conformance, and release verification work |
| `area:release` | Platform packaging, publication, installation migration, and legacy transition |
| `area:discovery` | Tool availability, withdrawal, and eager/lazy discovery |
| `area:presentation` | Display hints, status UI, and build progress |
| `area:context` | Extension-authored guidance and prompt contributions |
| `area:downloads` | Workspace saves and host-brokered file writes |
| `area:cache` | Cache invalidation, lifetime, and project scoping |
| `area:launcher` | Source builds, startup, and prebuilt fallback |
| `area:workflow` | Ticket-store maintenance and agent workflow, outside the product backlog |

The initial backlog has 20 modernization tickets: 11 core and 9 follow-up,
including one epic in each scope. All remain draft. Labels do not promote or
claim work, and grouping legacy archival under release does not authorize it.

## CLI examples

```sh
# Entire open modernization backlog, including drafts and epics.
git ticket list --label initiative:terva-modernization

# Core work, and security work across both scopes.
git ticket list --label scope:core
git ticket list --label area:security

# Investigations only: label and type filters combine.
git ticket list --label initiative:terva-modernization --type spike

# Include completed and archived work in machine-readable output.
git ticket list --all --label initiative:terva-modernization --json
```

**Repeated `--label` flags match any label (OR), not all labels.** For example,
this returns core work plus security work from either scope:

```sh
git ticket list --label scope:core --label area:security
```

For an intersection, filter the structured JSON labels with `jq`:

```sh
git ticket list --label initiative:terva-modernization --json |
  jq -r '.tickets[] | select(.labels | contains(["scope:core", "area:security"])) | [.id, .title] | @tsv'
```

For a compact grouping export, keep the stable ID, title and label array:

```sh
git ticket list --all --label initiative:terva-modernization --json |
  jq '[.tickets[] | {id, title, status, type, priority, labels}]'
```

Run `git ticket ready` to find startable work. A label filter reports a grouping;
it does not evaluate that grouping as a ready queue.
