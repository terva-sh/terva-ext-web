# Platform validation policy

The user selected Linux amd64 as the primary development platform on
2026-09-16. Other archive targets are validated on GitHub runners at release
tag time, with manual dispatch for rehearsals. Their absence on the local
Forgejo runner does not block SDK/config development. Cross-compilation is
build evidence only; a tagged candidate must pass its native runtime checks
before its assets are announced as validated.

| Archive | Runtime environment | When / operator |
| --- | --- | --- |
| linux/amd64 | Local Linux x86_64 and Forgejo docker; GitHub ubuntu-24.04 | Development and tagged candidate / maintainer + CI |
| linux/arm64 | GitHub ubuntu-24.04-arm | Tagged candidate / GitHub Actions |
| darwin/amd64 | GitHub macos-15-intel | Tagged candidate / GitHub Actions |
| darwin/arm64 | GitHub macos-15 | Tagged candidate / GitHub Actions |
| windows/amd64 | GitHub windows-2025, Git Bash | Tagged candidate / GitHub Actions |

All five GitHub runner labels were available in the 2026-09-16 manual
rehearsal. Jobs assert GOOS/GOARCH and actual runner architecture. Windows
supports the manifest's explicit `bash ./run.sh` launch via Git Bash/MSYS;
verify native Terva spawning that launcher, including paths with spaces.
No emulation is required or counted as native host integration. If a runner
label or architecture is unavailable, report failure and resolve access before
claiming that target validated; do not silently substitute a cross-build.

## Automation and remaining prerequisites

Read-only inspection found Forgejo CI and manual snapshot workflows using
read-only contents permissions and Go 1.25 containers. SDK migration updates
these to Go 1.27, explicitly approved by the user. The repository Actions
secret-name API returned no configured secrets. No secret values were read.
Current GoReleaser publication is disabled. The GitHub mirror destination and
transport access are verified below; native validation runs there without
changing repository settings.

References inspected in sibling checkouts:
- `terva-ext-index/.github/workflows/release.yml`: native runner matrix,
  version-tag/manual triggers, github.server_url guard, artifact upload before
  a separate publishing job.
- `terva-ext-jmap-mail/.github/workflows/release.yml`: setup-go from go.mod,
  offline vendored builds, tag/manifest/code version agreement.
- `terva-ext-caldav/docs/release-process.md`: source pin verification. Its
  separate-history/source-only release policy does not replace this fork's
  history or archive policy.

Use GitHub-only job guards because Forgejo also discovers `.github` workflows.
Validation needs only contents:read and artifact upload; a separate publication
job may use the ephemeral GITHUB_TOKEN with contents:write after all validation
jobs succeed. No permanent credential is needed for native test jobs. Any
publication token availability remains a release prerequisite; never record
its value. Tag runs validate before publication;
manual dispatch validates without publishing. Implement this in the release
validation batch, preserving the existing five-target archives and skills.

## Reproducible report

Each job uploads a report and logs containing:

- source candidate commit and public mirror commit if different, tag, workflow
  run URL, runner image, native OS/architecture, Go version;
- source/archive mode, exact Terva version and checksum (initial supported floor
  and current baseline v0.137.0), extension binary/archive SHA-256, manifest and
  hello version, tested archive file listing;
- separate temporary source, extracted archive, legacy-install, legacy-data,
  replacement-install/data, workspace A/B and TERVA_HOME fixture directories;
- source offline build and launcher, archive launch without Go, JSON-only hello,
  six tools, command, fixture network text/image calls, session-switch saves,
  shutdown, and actual Terva host integration outcomes;
- synthetic legacy settings/credentials only, upgrade preservation and rollback
  results, skipped checks with reasons, explicit pass/fail per check.

Never read installed user credentials or enable old/new installations together.
Checksums bind evidence to the exact candidate; changing source invalidates that
candidate's report. Go 1.27+, Bash and a C compiler are needed for source/race
checks; archives require Bash and the native executable, not Go. Reports must
distinguish subprocess protocol checks from actual host launch. Configuration and
release tickets supply the final migration fixtures and assertions.

## GitHub mirror and executable validation

The user supplied `https://github.com/terva-sh/terva-ext-web` and SSH URL
`git@github.com:terva-sh/terva-ext-web.git`. Read-only GitHub metadata and
`git ls-remote` verified a public empty repository with authenticated push
access on 2026-09-16. Added that URL as local remote `mirror`; origin remains
Forgejo. No remote settings were changed.

`.github/workflows/release-validation.yml` runs only on GitHub version tags
or manual dispatch. Linux builds all archives once; five native jobs check
source and extracted archives. Manual runs use snapshot versions, tag runs
use the tag version, and both skip publication. Reports bind source commit,
archive/binary checksums, runner and Go/host versions to completed checks.
Archive tests verify launch without Go and paths containing spaces, exact
version and research skill inclusion, subprocess wire and actual host driver.
No permanent credential or write permission is needed by this workflow.

Local commands (Go 1.27+, Python 3.13, Bash and C compiler):

```sh
just release-snapshot
python3 scripts/validate-archive.py --source linux amd64
python3 scripts/validate-archive.py dist linux amd64
```

Source tests build a race-instrumented extension; archive tests execute the
exact static release executable via WEB_CONFORMANCE_BINARY and the extracted
launcher via WEB_CONFORMANCE_INSTALL. Reports are ignored build artifacts and
uploaded by CI. Full CLI install/upgrade/rollback remains release validation.

## Rehearsal findings and fixes

The first manual run, [35146845707](https://github.com/terva-sh/terva-ext-web/actions/runs/35146845707),
passed source and archive checks on the four Linux/macOS targets and exposed
Windows failures. Native validation found issues that cross-compilation and
Linux tests could not establish:

- `filepath.IsAbs` does not reject Windows root-relative paths. Save preflight
  and writes now use `filepath.IsLocal`, with regression tests for rooted,
  drive-relative, UNC, reserved-device and alternate-stream paths. Git metadata
  guards also reject case and trailing-dot/space aliases.
- Timestamp-based LRU access ranks can tie on Windows. A lock-protected access
  sequence now determines eviction order; TTL and displayed storage times still
  use timestamps. Sleeps would hide the production bug, so tests remain fast.
- Terva's native Windows driver cannot execute a shebang script directly. The
  manifest now specifies `exec: bash` with `args: ["./run.sh"]` on every platform.
  This preserves one launcher and requires Git Bash on Windows as documented.

A new regression initially assumed `COM1.txt` remained a reserved device name.
Go 1.27's Windows implementation explicitly queries the OS for names with
extensions because modern Windows permits them. The test now uses `COM1`,
which is reserved; the implementation retains Go's native locality semantics.

These findings are tracked in TKT-01M2NYSRTKARHFGRF6XFEMM8A6 (Fix native Windows
save paths and cache eviction ordering) and TKT-01M2NZD0DP9MTFSJ62F8BZV96V
(Invoke Bash explicitly for native Windows host launch).

## Passing native evidence, 2026-09-16

[Manual run 35148398553](https://github.com/terva-sh/terva-ext-web/actions/runs/35148398553)
passed packaging and all five native jobs at candidate
`7fde307f42609844d87881101a13c14a4d53dcc3`. The ten downloaded source/archive
reports are preserved in [native-2026-09-16.json](../validation/native-2026-09-16.json),
including archive/binary SHA-256, runner image, native architecture, Go version,
Terva v0.137.0 driver version, file lists and completed checks.

| Target | Source vet/race/conformance/host | Archive checksum/launcher/conformance/host |
| --- | --- | --- |
| Linux amd64 | Passed | Passed |
| Linux arm64 | Passed | Passed |
| macOS amd64 | Passed | Passed |
| macOS arm64 | Passed | Passed |
| Windows amd64 (Git Bash) | Passed | Passed |

Archives use snapshot version `0.4.0-next`. All five verified the bundled
research skill, matching binary/manifest version, launch without Go from a path
containing spaces, protocol conformance and actual published Terva driver
loading. No release tag, binary release publication or remote settings change
was performed. The workflow has read-only contents permission; publication
remains disabled for both manual and tag runs.

These reports certify this candidate and these archive checksums. Documentation
and ticket completion commits do not change the tested implementation. Any
source, dependency, manifest, launcher or packaging change needs fresh native
evidence; the eventual release tag must run the workflow again for its exact
versioned archives. Full CLI install, legacy settings upgrade and rollback are
still tracked by TKT-01M2NQ0DQ10D7KQC2H7W1XVX0A (Validate replacement installation,
upgrades and release platforms); native driver loading does not replace those
checks.

## Terva-only configuration revalidation

[Run 35150980130](https://github.com/terva-sh/terva-ext-web/actions/runs/35150980130)
passed packaging and all five native source/archive jobs for cleanup candidate
`b162e831b076d85084f89c701293fdfbb436f129`. The [ten reports](../validation/native-terva-only-2026-09-16.json)
record the new archive and binary hashes. This run covers removal of legacy
configuration, the host-configured test fixtures, ignored-file regressions and
the rewritten README and research skill in the archives.

The old reports remain as evidence for their own candidate. Full CLI installation,
manual migration and rollback still require the release-validation ticket.
