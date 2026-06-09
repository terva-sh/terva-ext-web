# zot-web

A [zot](https://github.com/patriceckhart/zot) extension that gives the agent
web access through three LLM-callable tools:

- **`web_search(query, count?)`** — ranked results (title, URL, snippet).
- **`web_fetch(url, max_chars?)`** — a page's main content as Markdown. Image
  URLs are replaced with compact `[image:N]` placeholders to save tokens.
- **`web_images(url)`** — resolve the `[image:N]` placeholders from a previously
  fetched page back to their URLs. Served from cache, so it costs no network.

Single static Go binary, no runtime dependencies. It implements the zot
extension wire protocol directly (no dependency on the zot module).

> **Status: v0.** Search (Tavily + SearXNG), the SSRF-guarded fetcher, and
> article extraction (`go-readability` → `html-to-markdown` with GFM tables,
> image indexing, and a heuristic tag-stripper fallback) are all functional.
> Design rationale lives in the zot repo at
> `docs/plans/web-tools-extension-research.md`.

## Quick start (`just`)

With [`just`](https://github.com/casey/just) installed, the whole flow is two
commands — repeatable on any machine:

```bash
just install                 # build, then (re)install into $ZOT_HOME/extensions/
just configure-searxng       # point it at the default local SearXNG (127.0.0.1:11984)
# or target a specific instance:
just configure-searxng https://searx.example/
```

`just configure-searxng` writes `config.json` into the installed extension's
data dir, resolving that dir from `zot ext list` so it works regardless of OS
(macOS, Linux) or a custom `$ZOT_HOME`. A bare `host:port` is accepted and gets
an `http://` prefix. The default instance is the `SEARXNG_URL` variable at the
top of the `justfile`.

`just install` removes and recopies the install dir, but **preserves an
existing `config.json`** across the reinstall — so you only need
`configure-searxng` once (or to change instances).

See `just --list` for the rest (`try`, `lint`, `test`, …).

## Manual build & install

```bash
go build -o zot-web .                  # exec name must match extension.json
zot ext install /path/to/zot-web       # copies the dir into $ZOT_HOME/extensions/<dir-basename>/
cp zot-web "$ZOT_HOME/extensions/zot-web/zot-web"   # see caveat below
# or, for one session straight from the working copy:
zot --ext /path/to/zot-web
```

The directory must contain `extension.json` (pointing at the `./zot-web`
binary you built) — already included here. Two gotchas:

- `zot` does **not** build Go extensions for you (`language` is informational),
  so build first.
- When the source is a git repo, `zot ext install` copies **git-aware** and
  skips `.gitignore`d files — and `./zot-web` is gitignored (it's a build
  artifact). So the binary is *not* copied; copy it in manually as shown, or
  just use `just install`, which does this for you.

The install dir is named after the source folder's basename (here, `zot-web`),
not the manifest `name` (`web`). `zot --ext` runs from the working copy
directly, so it needs no copy step.

## Configure

Settings come from `config.json` in the extension's data dir
(`$ZOT_HOME/extensions/zot-web/config.json`), with environment variables taking
precedence. `just configure-searxng` (above) writes this file for you; to do it
by hand, start from the default Tavily backend:

```bash
export TAVILY_API_KEY=tvly-...
```

Or switch to a self-hosted SearXNG instance (no key, private):

```jsonc
// $ZOT_HOME/extensions/zot-web/config.json
{
  "search_backend": "searxng",
  "searxng_url": "http://127.0.0.1:11984",
  "allow_local_hosts": ["localhost", "intranet.example", "10.0.0.0/24"]
}
```

> SearXNG must have `json` listed under `search.formats` in its `settings.yml`,
> otherwise its API returns `403`.

### All settings

| config.json key | env override | default | meaning |
|---|---|---|---|
| `search_backend` | `ZOT_WEB_SEARCH_BACKEND` | `tavily` | `tavily` or `searxng` |
| `tavily_api_key` | `TAVILY_API_KEY` | — | Tavily bearer token |
| `searxng_url` | `ZOT_WEB_SEARXNG_URL` | — | SearXNG base URL |
| `fetch_max_bytes` | `ZOT_WEB_FETCH_MAX_BYTES` | `2097152` | response body cap |
| `fetch_timeout_sec` | `ZOT_WEB_FETCH_TIMEOUT_SEC` | `25` | per-fetch timeout |
| `fetch_inline_images` | `ZOT_WEB_FETCH_INLINE_IMAGES` | `false` | keep image URLs inline instead of `[image:N]` placeholders |
| `fetch_cache_ttl_sec` | `ZOT_WEB_FETCH_CACHE_TTL_SEC` | `600` | how long a rendered page stays cached (`0` = no expiry) |
| `fetch_cache_max_entries` | `ZOT_WEB_FETCH_CACHE_MAX_ENTRIES` | `32` | max cached pages, LRU-evicted (`0` = caching off) |
| `allow_local_hosts` | `ZOT_WEB_ALLOW_LOCAL_HOSTS` (comma-sep) | — | SSRF escape hatch (see below) |

### Images and the page cache

By default `web_fetch` strips image URLs out of its Markdown, leaving a short
`[image:N: alt]` handle where each image was. This keeps long CDN URLs out of
the model's context. To get the actual links, the model calls
`web_images(url)`, which returns the `N → URL` mapping. Because every fetched
page is cached (in memory, per the TTL/size settings above), that follow-up
call normally costs no network request. Set `fetch_inline_images: true` to
restore inline image URLs and disable the indexing.

## Security: SSRF protection + the local allowlist

Because the model chooses the URL, `web_fetch` is the main attack surface
(prompt-injected pages can try to make it hit internal services). By default it:

- allows `http`/`https` only;
- resolves the host and **refuses private/reserved/loopback/link-local
  addresses** — including the cloud metadata address `169.254.169.254`;
- dials the validated IP directly (closing the DNS-rebinding gap) and re-checks
  on every redirect; caps redirects, time, and response size.

To deliberately reach local services, add them to **`allow_local_hosts`**. Each
entry is one of:

- a **hostname** — matched against the request host (e.g. `localhost`,
  `grafana.internal`);
- an **IP** — matched against the resolved address (e.g. `127.0.0.1`);
- a **CIDR** — matched against the resolved address (e.g. `192.168.1.0/24`).

This is a precise escape hatch, not an "allow all local" switch: only the
targets you list are exempted.

## Roadmap

- [x] Replace the heuristic HTML extractor with `go-shiori/go-readability` +
      `JohannesKaufmann/html-to-markdown` (heuristic kept as a fallback).
- [x] GFM table rendering + image indexing (`[image:N]` + `web_images`) with an
      in-memory page cache.
- [ ] Infobox / vertical key-value tables → cleaner key/value lists (irregular
      tables still degrade to spaced blocks today).
- [ ] More search backends (Brave, Serper, Exa) behind the same interface.
- [ ] Optional JS rendering fallback (e.g. Jina Reader) — deferred for now.
- [ ] Release binaries (goreleaser) so `zot ext install <git-url>` needs no
      local build.
