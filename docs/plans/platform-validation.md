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

GitHub labels are planned selections, not observed successful jobs. Assert
GOOS/GOARCH and actual runner architecture in the jobs. Windows supports the
manifest's Bash launcher via Git Bash/MSYS, not a PowerShell replacement;
verify native Terva spawning that launcher, including paths with spaces.
No emulation is required or counted as native host integration. If a runner
label or architecture is unavailable, report failure and resolve access before
claiming that target validated; do not silently substitute a cross-build.

## Automation and remaining prerequisites

Read-only inspection found Forgejo CI and manual snapshot workflows using
read-only contents permissions and Go 1.25 containers. SDK migration updates
these to Go 1.27, explicitly approved by the user. The repository Actions
secret-name API returned no configured secrets. No secret values were read.
Current GoReleaser publication is disabled. No GitHub remote or mirror URL is
configured in this checkout; the user's release mirror provides GitHub runners,
but its exact destination/access still needs verification before publication.

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
jobs succeed. No permanent credential is needed for native test jobs. Mirror
transport permissions and any publication token availability remain a release
prerequisite; never record their values. Tag runs validate before publication;
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
distinguish protocol-harness success from actual host launch. Configuration and
release tickets supply the final migration fixtures and assertions.
