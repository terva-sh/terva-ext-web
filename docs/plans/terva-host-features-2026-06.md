# Changes to consider for the June 2026 terva host features

> Notes record, written while terva's `feat/landscape-tier-a` branch
> landed approval modes, permission rules, tool-use hooks, an MCP
> client, and subprocess env sanitization (see terva's
> `docs/plans/harness-landscape-2026.md`). zot-web needs **no code
> change to keep working** — everything here is opportunity, not
> breakage. Ordered by value.

## 1. Read-only tool annotation (blocked on terva, highest value)

terva's new `plan` approval mode excludes all extension tools from
the registry because their side effects are unknown to the host. But
`web_search`, `web_fetch`, `web_images`, and `web_links` are exactly
the tools a planning/research session wants — they read the web and
touch nothing locally. (`web_fetch_image`/`web_fetch_raw` can write
workspace files, so their exclusion is correct.)

The fix is host-side: the extension `register_tool` frame (and the
permission classification behind it) would need an optional
`read_only: true` annotation — the same idea as MCP's `readOnlyHint`.
Once terva supports it, zot-web should annotate the four read-only
tools. Worth filing on the terva side; tracked there in the
harness-landscape doc's B8/protocol notes.

## 2. Recommended permission rules (README addition, doable now)

terva users can now pre-answer tool prompts with rules in
`$TERVA_HOME/config.json`. The README should suggest a starter set,
e.g.:

```json
{
  "permissions": [
    { "tool": "web_search", "decision": "allow" },
    { "tool": "web_fetch", "decision": "allow" },
    { "tool": "web_fetch_raw", "decision": "ask" },
    { "tool": "web_fetch_image", "args": "\"save_path\"", "decision": "ask" }
  ]
}
```

Rationale to document with it: the two fetch-and-save tools write
into the workspace, so they deserve a prompt in cautious setups; the
read-only four are safe to allow. Also worth one line: a pre-tool-use
hook (terva `docs/hooks.md`) can veto `web_*` calls by URL pattern —
e.g. block fetches of internal hostnames *before* they reach the
extension's own SSRF guard, as defense in depth.

## 3. Bundle-manifest contributions (when terva B8 lands)

terva is growing `extension.json` into a declarative bundle that can
also contribute skills, command templates, and suggested permission
rules (data, not code). When that ships, zot-web could bundle:

- a **web-research skill** (`SKILL.md`): how to chain `web_search` →
  `web_fetch` → `web_links`/`web_images`, paging with `offset`,
  citing URLs;
- the **recommended permission rules** from §2 as a manifest
  contribution. Note terva's trust rule: extension-supplied policy
  can never self-`allow` — ship the `ask` entries and let the README
  carry the `allow` suggestions for users to adopt.

## 4. Host env sanitization (compat note only)

terva now spawns extensions from a sanitized environment: loader and
interpreter injection vars (`LD_*`, `DYLD_*`, `PYTHONPATH`,
`NODE_OPTIONS`, …) are stripped. `TAVILY_API_KEY`,
`TERVA_WEB_*`-style overrides, `PATH`, and `HOME` pass through
untouched, so zot-web's config resolution is unaffected. Only worth a
line in the README's environment section. One caveat for the
`run.sh` launcher: it must not rely on any stripped variable
(it doesn't today).

## 5. Doc naming sweep ($ZOT_HOME → terva, cosmetic)

README/justfile still say `$ZOT_HOME`, `zot ext list`,
`logs/ext-web.log` under the zot name. Both hosts work (terva reads
legacy dirs through its envcompat layer), but the docs should
eventually lead with terva naming and mention zot compatibility,
matching the module identity (`github.com/terva-sh/zot-web`) and the
public mirror. Fold into the next README pass rather than a dedicated
change.
