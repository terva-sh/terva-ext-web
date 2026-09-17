# Release process

The planned first terva-ext-web release is 0.4.0. No binary release has been
published. Forgejo is the source repository; GitHub mirrors the same history
and runs native validation. The inherited zot-web release-cut process is
retired. Do not create orphan history or reuse its `cut/*` tags.

| Purpose | Repository |
| --- | --- |
| Source, tickets and review | `ssh://git@git.local.sothr.com:2222/terva-sh/terva-ext-web.git` |
| Native runners and release mirror | `git@github.com:terva-sh/terva-ext-web.git` |

## Local checks and archives

```sh
just ci
just ticket-check
just release-verify
just release-snapshot
```

Snapshots use version `0.4.0-next`. GoReleaser builds five targets and stamps
matching binary and manifest versions. Archives include the launcher, README,
license and research skill. SHA-256 checksums are in `dist/checksums.txt`.
Source manifest and `internal/version.Version` must agree; a test checks them.

## Native validation

After Forgejo review and checks, push the candidate to GitHub. For a rehearsal:

```sh
gh workflow run release-validation.yml --repo terva-sh/terva-ext-web --ref main
```

The workflow also runs on `v*` tags. It packages the candidate once, then tests
source and extracted archives on native Linux amd64/arm64, macOS amd64/arm64
and Windows amd64 runners. Windows requires Git Bash. Reports include the
candidate commit, archive and binary checksums, runner images and results.
See [platform validation](platform-validation.md) for the commands and evidence.

The 2026-09-16 runs passed all five targets, first at `7fde307` and again
after Terva-only configuration cleanup at `b162e83`. Each report applies only
to its recorded candidate and archives. Later source, dependency, launcher,
manifest or packaging changes require new evidence. The eventual tagged release must test its exact versioned archives.

## Publication requirements

GoReleaser publication is disabled. Both manual and tag workflows validate
without publishing. A separate publication change must depend on successful
native jobs and use the approved destination. Do not copy the old Gitea
publication settings into a GitHub release job without reviewing them.

Before publication, complete isolated full-CLI source/archive installation,
manual settings migration and rollback checks. Use synthetic credentials and
host-reported paths. Confirm skill discovery and that the two `web`
installations are never enabled together. The user must approve the concrete
release candidate and publication action through the release ticket.

Only after the replacement is published and smoke-tested should the old
repository receive a migration notice. Archiving it requires explicit ticket
authorization. See [the critical path](modernization-critical-path.md).
