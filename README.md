# zot-web

A [zot](https://github.com/patriceckhart/zot) extension that gives the agent
web access through three LLM-callable tools:

- **`web_search(query, count?)`** — ranked results (title, URL, snippet).
- **`web_fetch(url, max_chars?, offset?)`** — a page's main content as Markdown,
  led by a metadata block. Image URLs are replaced with compact `[image:N]`
  placeholders to save tokens; `offset` pages through long documents.
- **`web_images(url)`** — resolve the `[image:N]` placeholders from a fetched
  page back to their URLs (plus dimensions, caption, and source page). Served
  from cache when warm; fetches on a cold cache, so it also works standalone.

Single static Go binary, no external runtime services. It implements the zot
extension wire protocol directly (no dependency on the zot module).

> **Status: v0.** Search (Tavily + SearXNG), the SSRF-guarded fetcher, and
> article extraction (readability via the maintained
> `codeberg.org/readeck/go-readability` fork → `html-to-markdown` with GFM
> tables, image indexing, and a heuristic tag-stripper fallback) are all
> functional. Design rationale lives in the zot repo at
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

## Manual install & the `run.sh` launcher

```bash
# From a git URL — zot shallow-clones the repo (it does NOT build Go sources):
zot ext install https://git.local.sothr.com/warricksothr/zot-web.git
# From a local checkout:
zot ext install /path/to/zot-web
# Or, for one session straight from the working copy:
zot --ext /path/to/zot-web
```

`extension.json` points `exec` at **`./run.sh`**, a launcher that compiles the
binary on first launch (and after any source change) and then execs it. zot
never builds Go sources itself (`language` is informational), so this is what
makes a git-URL install work without committing a platform-specific binary — at
the cost of needing a **Go 1.25+ toolchain on `PATH`** on the host and a
one-time build before the extension responds. The build is **offline**
(`go build -mod=vendor` against the committed `vendor/` tree), so the first
launch needs no network and can't hang on a module fetch — which matters because
zot blocks its startup until the extension sends its `hello`. Build chatter goes
to stderr (zot captures it to `$ZOT_HOME/logs/ext-web.log`); the compiled
`./zot-web` is gitignored.

`just install` still builds locally and copies the binary into the install dir,
pre-seeding it so the first launch skips the build. The install dir is named
after the source folder's basename (here, `zot-web`), not the manifest `name`
(`web`); `zot --ext` runs from the working copy directly.

### Dependencies are vendored

`vendor/` is committed so the first-launch build is fast and offline (see
above). After changing dependencies, refresh it with `just vendor` (runs
`go mod tidy` + `go mod vendor`) and commit the result alongside
`go.mod`/`go.sum`. **Re-evaluate this approach if `vendor/` grows large**
(currently ~6 MB across a handful of modules): past some point, shipping
prebuilt per-platform binaries (goreleaser) beats carrying a big vendor tree.

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

### `web_fetch` output

The output leads with a small metadata block so the model can tell a short page
from a truncated dense one:

```text
# Artificial intelligence
https://en.wikipedia.org/wiki/Artificial_intelligence
Content-Type: text/html; charset=UTF-8
Chars: 0-500 of 397898
Images: 17 (shown as [image:N]; resolve with web_images)

**Artificial intelligence** (AI) is the capability of …

…[397398 more chars; continue with offset=500]
```

- A `Final-URL:` line appears only when redirects landed somewhere other than
  the requested URL.
- `Chars: start-end of total` reports the returned window against the full
  rendered length. When `end < total`, the trailing hint gives the exact
  `offset` to pass to the next `web_fetch` call to keep reading — the page is
  already cached, so continuation costs no extra network request.
- Relative links and image sources are resolved against the **final** URL after
  redirects, so an `http→https` redirect doesn't leave stale links in the body.

Binary responses (images, PDFs, octet-streams) are not dumped as raw bytes —
`web_fetch` returns a one-line summary like
`[image/png content, 40075 bytes — not rendered as text]` instead. Textual
types (`text/*`, JSON, XML, SVG) pass through normally.

readability drops `<table>` elements from article content, so data tables
(e.g. large sortable Wikipedia tables) are recovered separately and appended
under a `## Tables` heading, rendered leniently (cell text flattened, images
dropped, ragged rows padded). Each table is capped at 50 rows with a
truncation note; the dropped rows are not stored, so they are not reachable via
`offset`.

### Images and the page cache

By default `web_fetch` strips image URLs out of its Markdown, leaving a short
`[image:N: alt]` handle where each image was. This keeps long CDN URLs out of
the model's context. To get the actual links, the model calls `web_images(url)`,
which returns each handle's URL plus dimensions, the nearest `<figcaption>`
caption, and the enclosing source-page link (e.g. a Wikimedia `File:` page).

**The placeholder contract:** `[image:N]` in `web_fetch` maps to `[image:N]` in
`web_images` for the same URL. Ids are assigned in document order and are stable
for a cached page; identical image URLs are de-duplicated to a single id.

Every fetched page is cached (in memory, per the TTL/size settings above), so
`web_images` normally costs no network request. If called for a URL that was
never fetched (or whose cache entry expired), it transparently fetches and
renders the page first — it does **not** error, so it is safe to call directly.
Set `fetch_inline_images: true` to restore inline image URLs and disable the
indexing (and the `web_images` workflow).

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

- [x] Replace the heuristic HTML extractor with readability +
      `JohannesKaufmann/html-to-markdown` (heuristic kept as a fallback).
      _(Uses the maintained `codeberg.org/readeck/go-readability/v2` fork —
      `go-shiori/go-readability` is now deprecated.)_
- [x] GFM table rendering + image indexing (`[image:N]` + `web_images`) with an
      in-memory page cache.
- [x] Recover data tables that readability strips, rendered leniently under a
      `## Tables` section (row-capped). Tables land at the end, not inline.
- [ ] Infobox / vertical key-value tables → cleaner key/value lists (irregular
      tables still degrade to spaced blocks today).
- [ ] More search backends (Brave, Serper, Exa) behind the same interface.
- [ ] Optional JS rendering fallback (e.g. Jina Reader) — deferred for now.
- [ ] Optional prebuilt per-platform binaries (goreleaser) so the `run.sh`
      launcher can skip the on-host build and drop the Go-toolchain requirement.
