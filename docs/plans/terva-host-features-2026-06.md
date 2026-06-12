# June 2026 terva host features — what zot-web adopted

> Record of how zot-web took up the terva host features that landed in
> June 2026 (terva's `feat/landscape-tier-a`: approval modes,
> permission rules, hooks, MCP, extension env sanitization, bundle
> contributions). **zot-web stays a zot extension** — every change here
> is additive and invisible to stock zot, because a tool that runs on
> both hosts grows both ecosystems.

## Done

### Network-read authority annotation (code)

`internal/proto` gained a `WithAuthority()` tool option (with a
`NetworkRead()` shorthand); all six tools register with
`"authority": "network-read"`, since each reaches the network. terva
gates on it — only `local-read` is auto-allowable — so the web tools
are prompted in `--approval workspace`/`auto-edit` and refused in
`--approval plan`. The field is emitted only when set, so a
non-annotated tool's `register_tool` frame is byte-for-byte what zot
has always received; stock zot ignores the field entirely. Pinned by a
proto test.

### Bundle skill (data)

`skills/web-research/SKILL.md` ships in the repo and is copied in by
`zot ext install` (git-aware). terva discovers extension-bundled
skills automatically; zot does not load them, so it's a no-op there.

### Suggested permission rules (manifest data)

`extension.json` gained a restrict-only `permissions` block: `ask`
before `web_fetch_raw` and `web_fetch_image` (they write files). terva
honors this even in yolo; the user overrides it with an `allow` in
their own config. zot ignores the unknown manifest key. The four
readers carry no rule and follow the host's approval mode.

## Layer ordering (settled on the terva side)

The terva permission layers are evaluated **user → project →
extension**, first match wins. The user is sovereign: an explicit user
rule beats any restrict-only suggestion, so zot-web's manifest
`ask`-defaults are genuinely overridable (this is why shipping them is
safe rather than presumptuous). Project and extension layers may only
tighten, never `allow`.

## Not adopted (deliberately)

- **Naming sweep ($ZOT_HOME → terva).** Skipped on purpose: zot-web is
  and stays a zot extension. Docs lead with zot naming; the
  "Host integration" README section notes the terva extras. Both
  `$ZOT_HOME` and `$TERVA_HOME` resolve through terva's envcompat
  layer regardless.
- **Hooks / MCP from the manifest.** Those run additional programs and
  stay an explicit user-config decision on the host; not something a
  bundle should contribute.

## Compat note (no action)

terva spawns extensions from a sanitized environment (strips `LD_*`,
`DYLD_*`, `PYTHONPATH`, `NODE_OPTIONS`, …). `TAVILY_API_KEY`,
`SEARXNG_*`, `PATH`, `HOME` pass through, so config resolution and
`run.sh` are unaffected.
