# Identity and packaging — 2026-09-16

This batch follows fork `c50d773c60250a7315c37c2aedbeb891cb12f6a4` and its
recorded tree `6c368f89b3e622c7d4515ac4e02d47caed3a95c7`, on topic branch
`t3code/implement-identity-packaging`. The starting worktree was clean and had
no remotes. The sibling Terva checkout remained at
`7f754b9bb7c6754284dc4f3cc5fa6a96525de49d`; its unrelated untracked ticket was
left untouched.

The unchanged-source baseline passed `just ci` (vet, format, race tests,
subprocess conformance for both inherited profiles, vendor synchronization)
and a full five-target GoReleaser snapshot. Go was absent from PATH but the
installed mise Go 1.27.1 toolchain was available and used explicitly. This
closes the setup session's missing-toolchain validation gap; it does not
claim testing with Go 1.25 or a current Terva host.

## Decisions

- Repository/module/binary and default User-Agent become `terva-ext-web`.
  Internal imports and linker symbols move together. Version 0.4.0 is the
  first new release version, distinguishing the fork from inherited 0.3.1.
- Keep `web`, the six `web_*` tools, `/web-cache`, manifest permissions, and
  all configuration inputs. Renaming configuration variables now would mix
  identity work with the later precedence/credential migration.
- Ship the research skill and license in every archive. Stamp an archive-only
  manifest to match the linked binary version, including 0.4.0-next snapshots.
  Snapshot versions are explicit because inherited cut/* tags are not semver
  releases and GoReleaser otherwise derives 0.0.1-next.
- Retire the inherited orphan-history release-cut script rather than merely
  rename it: it could publish to zot-web, guess local mirror paths, and reuse
  inherited cut/* tags. Local verify/snapshot commands replace it. Historical
  release decisions remain in the release-process document as history.
- The planned release owner is `terva-sh`, not the inherited GoReleaser owner
  `warricksothr`; the repository is `terva-ext-web`. Publication is disabled in
  GoReleaser and CI pending actual destination verification in the later batch.
  No remote, tag, release, repository setting, or installed extension changed.
- Delegate installation to Terva rather than automatically remove and recopy
  existing installations. The old helper could destroy installation-local
  state and mask an identity collision. Migration must use resolved host paths
  and preserve the old installation/settings for rollback.
- Retain the handwritten protocol and legacy conformance profiles. SDK adoption,
  the known download cwd race and its session-switch test, configuration/secret
  migration, and current-host support remain separately reviewable work.

## Validation limits

All five archives are cross-built and inspected. The Linux amd64 archive can
be executed here without Go on PATH; other OS/architecture runtime checks and
real Terva install/upgrade/rollback remain prerequisites for publication.
Windows retains the inherited Bash launcher requirement, now selecting the
.exe binary under Git Bash/MSYS. No Windows-native launch support is claimed.

## Final checks

- `just ci` passed after the rename: formatting, vet, race tests, both inherited
  subprocess conformance profiles, and a clean regenerated vendor tree.
- `just build` and `run.sh --version` reported `terva-ext-web 0.4.0`; the
  launcher rebuilt after source changes. Shell syntax checks passed for the
  launcher and both packaging scripts; `just release-verify` passed.
- `just release-snapshot` produced all five linux/darwin amd64/arm64 and
  windows amd64 archives. Each SHA-256 matched checksums.txt. Each archive
  contained exactly the binary, README, LICENSE, run.sh, extension.json, and
  skills/web-research/SKILL.md. Manifest versions were 0.4.0-next; each binary
  had the renamed Go module identity and embedded snapshot version.
- The extracted Linux amd64 archive ran with PATH limited to /usr/bin:/bin
  (no Go toolchain). Its version output and hello both reported 0.4.0-next;
  the hello name was web, all six tool names and web-cache were registered,
  stdout was JSON-only, stderr was empty, and shutdown was acknowledged.
- Archive inspection caught GoReleaser treating `dst: extension.json` as a
  directory; using `dst: .` with `strip_parent: true` places the stamped
  manifest at the archive root. The corrected archives passed reinspection.

The source launcher and archive checks used temporary directories and no real
configuration or credentials. No installed extension was changed. Cross-build
success is not a claim that every platform was executed locally.
