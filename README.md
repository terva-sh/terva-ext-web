# terva-ext-web

Web search, page retrieval and image downloads for Terva. The extension uses
the published Terva v0.137.0 SDK and keeps the `web` identity, six `web_*` tools
and `/web-cache` command. Stock zot is not supported.

Terva v0.137.0 is the tested host version. Source builds require Go 1.27+ and
Bash. Archives require Bash, or Git Bash on Windows, and no Go toolchain.
Native source and archive tests pass on Linux amd64/arm64, macOS amd64/arm64
and Windows amd64. See the [validation evidence](docs/plans/platform-validation.md).
The planned first release is 0.4.0. No binary release has been published.

## Install

The source repository is [terva-sh/terva-ext-web](https://git.local.sothr.com/terva-sh/terva-ext-web).
[GitHub](https://github.com/terva-sh/terva-ext-web) mirrors its history and runs
native release validation.

With Go and Bash on PATH, use Terva to install a local checkout:

```sh
terva ext install /path/to/terva-ext-web
```

For development, [just](https://github.com/casey/just) provides these commands:

```sh
just build
terva --ext /path/to/terva-ext-web
just install
```

The manifest runs `bash ./run.sh`. On a source installation, the launcher builds
from `vendor/` on first use and after source changes. It needs no module downloads.
Build messages go to stderr; stdout carries the extension protocol.

`just release-snapshot` builds archives containing the executable, manifest,
launcher, README, license and research skill. Extract one into a new directory
and try it with `terva --ext /path/to/extracted-directory`. Full CLI installation,
upgrade and rollback validation remains open before publication.

## Tools

| Tool | Use |
| --- | --- |
| `web_search` | Search with Tavily or SearXNG. Returns titles, URLs, snippets and publication dates when available. |
| `web_fetch` | Read a page as Markdown. Use `offset` to continue through long pages. |
| `web_images` | List image URLs, captions and dimensions. Resolves the `[image:N]` handles in fetched pages. |
| `web_links` | List deduplicated hyperlinks with absolute URLs and anchor text. |
| `web_fetch_image` | View an image, resize it, or save it to the workspace. |
| `web_fetch_raw` | Save the original response bytes to a workspace file. |

Search accepts `query`, `count`, `freshness`, `include_domains`,
`exclude_domains` and `depth`. Freshness is `day`, `week`, `month` or `year`.
Tavily supports domain filters and advanced depth directly. SearXNG uses a
`site:` query hint and post-filtering for domains.

## Configure

Use Terva's extension configuration form or `terva ext config web`. For a
local SearXNG instance:

```sh
terva ext config web set search_backend=searxng searxng_url=http://127.0.0.1:11984
```

`just configure-searxng URL` runs that host command. It updates only the backend
and URL. For a LAN or VPN instance, also add the destination to
`allow_local_hosts`. SearXNG must enable `json` under `search.formats` in its
`settings.yml`, or it returns HTTP 403.

For Tavily, set `search_backend` to `tavily` and enter `tavily_api_key` in the
host's secret field. You can instead
supply `TAVILY_API_KEY` in the process environment.

### Configuration precedence

Settings come from `TERVA_EXT_WEB_*` environment variables, then explicit Terva
host values, then application defaults. `TAVILY_API_KEY` overrides the Tavily
secret field. The extension no longer reads standalone `config.json` files,
`ZOT_WEB_*` variables or `configuration_source`.

The manifest omits defaults so explicit `false`, `0` and empty values reach the
application. Blank nonsecret form input means unset. Empty optional strings
clear a value; empty numeric values fail validation.

The host's `allow_local_hosts` field accepts JSON array text, for example
`["localhost", "127.0.0.1", "::1", "search.internal"]`. It replaces the default
list. `TERVA_EXT_WEB_ALLOW_LOCAL_HOSTS` appends comma-separated entries. Use
`[]` in the host field, with no environment additions, to block loopback too.

Valid host updates apply to later calls and create a fresh cache. In-flight
calls finish with their captured settings and cache. Invalid updates report the
field name without its value and retain the last working configuration. Invalid
startup settings block every network tool; missing search credentials block
search only. Environment changes require a restart.

### Credentials

The Tavily key comes from the host secret field or `TAVILY_API_KEY`. A missing
or undecryptable host key stays unconfigured. An explicitly empty
`TAVILY_API_KEY` clears the effective credential.

Blank secret form input keeps the saved key. Use Terva's secret/config clear
operation to remove it. A masked input does not establish encryption at rest;
follow Terva's `terva secret init` and backup instructions for that setup.
This extension uses resolved configuration secrets, not the runtime secret broker.

The extension does not rotate credentials or delete old copies. Its manifest
keeps `data_secrets: true` because a preserved data directory may contain a key.
Provider errors omit raw response bodies, and configured key echoes are redacted
from search results and errors.

### Settings

| setting key | preferred env override | default | meaning |
|---|---|---|---|
| `search_backend` | `TERVA_EXT_WEB_SEARCH_BACKEND` | `tavily` | `tavily` or `searxng` |
| `tavily_api_key` | `TAVILY_API_KEY` | unset | Tavily bearer token |
| `searxng_url` | `TERVA_EXT_WEB_SEARXNG_URL` | unset | SearXNG base URL |
| `fetch_max_bytes` | `TERVA_EXT_WEB_FETCH_MAX_BYTES` | `2097152` | response body cap (maximum `33554432`) |
| `fetch_image_max_bytes` | `TERVA_EXT_WEB_FETCH_IMAGE_MAX_BYTES` | `5242880` | max encoded size of a `web_fetch_image` result after resize (maximum `20971520`) |
| `fetch_timeout_sec` | `TERVA_EXT_WEB_FETCH_TIMEOUT_SEC` | `25` | per-fetch timeout (maximum `60`) |
| `fetch_inline_images` | `TERVA_EXT_WEB_FETCH_INLINE_IMAGES` | `false` | keep image URLs inline instead of `[image:N]` placeholders |
| `fetch_cache_ttl_sec` | `TERVA_EXT_WEB_FETCH_CACHE_TTL_SEC` | `600` | how long a rendered page stays cached (`0` = no expiry; maximum `3600`) |
| `fetch_cache_max_entries` | `TERVA_EXT_WEB_FETCH_CACHE_MAX_ENTRIES` | `32` | max cached pages, LRU-evicted (`0` = caching off; maximum `128`) |
| `fetch_cache_max_bytes` | `TERVA_EXT_WEB_FETCH_CACHE_MAX_BYTES` | `67108864` | total bytes the page cache may retain, LRU-evicted (`0` = no byte bound; maximum `268435456`) |
| `user_agent` | `TERVA_EXT_WEB_USER_AGENT` | `terva-ext-web/<version>` | User-Agent for every fetch; `browser` expands to a common desktop-browser UA |
| `allow_local_hosts` | `TERVA_EXT_WEB_ALLOW_LOCAL_HOSTS` (comma-separated) | `localhost, 127.0.0.1, ::1` | private-address exceptions; the host field replaces the default, the env var appends |

### User-Agent

The default is `terva-ext-web/<version>`. `user_agent` or
`TERVA_EXT_WEB_USER_AGENT` overrides it. The per-call `user_agent` on
`web_fetch`, `web_fetch_raw` or `web_fetch_image` takes precedence and forces a
fresh request instead of using the cached page. The value `browser`, in any
case, expands to the bundled desktop Chrome User-Agent string.

The extension does not consult `robots.txt` for individual on-demand requests.
A future bulk crawler would need a separate policy and implementation.

## Page output and caching

`web_fetch` returns a title, URL, content type and character window before the
page body. Redirected requests also include `Final-URL`. The continuation hint
gives the next `offset`; while the page remains cached, each window uses the
same snapshot. Relative links resolve against the final URL.

HTML article text uses readability and Markdown conversion, with a text
fallback for pages without an extractable article. Data tables omitted by
readability appear under `Tables`, capped at 50 rows per table. Dropped rows
are unavailable through `offset`.

RSS and Atom feeds render up to 100 entries with title, date, link and summary.
Text encodings use the response charset, HTML metadata or content sniffing.
PDFs use their text layer and page markers. There is no OCR. Unreadable PDFs
and other binary bodies return a summary suggesting a raw download.
`web_fetch_raw` preserves the received bytes, subject to download limits.

Image URLs become `[image:N]` handles by default. `web_images` returns the
matching URLs, dimensions, captions and source-page links. IDs remain stable
within a cached page and identical URLs share an ID. Discovery covers lazy-load
attributes, `srcset`, picture sources, image links, social metadata and a
whole-page fallback. Set `fetch_inline_images: true` to retain inline URLs
instead of indexed handles.

`web_images`, `web_links` and `web_fetch_raw` use the page cache when possible
and fetch on a miss. The cache retains Markdown, gzip-compressed raw bodies,
links and images. It limits entries and retained bytes with LRU eviction. A
single page larger than the byte budget remains as the only entry. A page holds
at most 5,000 links and 2,000 images.

The cache belongs to the extension process and can serve multiple sessions.
It is not isolated by project. `/web-cache` lists entries; `/web-cache clear`
evicts them but does not cancel an in-flight fetch, which can populate the
cache after the clear. Project isolation remains a tracked follow-up.

## Image viewing and workspace saves

`web_fetch_image` accepts PNG, JPEG, GIF and WebP. `max_dimension` resizes the
longest edge without upscaling. PNG, JPEG and GIF retain their format; resized
WebP becomes PNG. Set `save_path` for a workspace file and `inject: false` for
a save without an image content block.

Saves reject nonlocal paths, escapes, symlinks and Git metadata paths. Parent
directories are created as needed. Existing files require `overwrite: true`.
Each call captures its workspace before the request, so a session switch cannot
redirect the save. Raw downloads use the same path policy.

Encoded image output is limited by `fetch_image_max_bytes`. Images are checked
against a 40-million-pixel limit before full decoding, and at most three
images decode or resize concurrently. The original download may exceed the
encoded output limit so it can be resized.

All tool responses must also fit the host's 4 MiB message limit, including JSON
escaping and base64. Use a smaller `max_dimension` or a save-only request when
an image exceeds that limit. Oversized image injection fails before a file is
written. Save-only requests retain their configured download limits.

## Network and content safety

Fetches allow HTTP and HTTPS only. The fetcher rejects private and special-use
addresses unless allowlisted, dials the validated IP to prevent DNS rebinding,
and repeats the checks on redirects. It limits redirects, time and response
size, and blocks known non-web service ports even on public hosts.

The default allowlist is `localhost`, `127.0.0.1` and `::1`. Other private,
link-local, documentation, benchmarking, CGNAT and multicast ranges remain
blocked. This includes cloud metadata at `169.254.169.254`.

Allowlist entries can be hostnames, IPs or CIDRs. A hostname entry trusts that
name's DNS results, including otherwise blocked addresses. Choose narrow
entries for private services. These checks also apply to SearXNG queries.

Fetched text remains untrusted input. Titles, snippets and image attributes
are flattened where they appear in generated metadata, but this cannot prevent
a model from following malicious page instructions. See the
[content-safety decision](docs/plans/untrusted-web-content.md).

## Terva permissions and research skill

All six tools declare `network-read`. The manifest requests `ask` for
`web_fetch_raw` and `web_fetch_image`, including image calls without a save.
The tested host policy is:

| Approval mode | Default behavior |
| --- | --- |
| `plan` | Denies the tools. |
| `ask`, `auto-edit`, `workspace` | Prompts for approval. |
| `yolo` | Allows them, including the manifest's writer `ask` rules. |

Explicit user allow/deny rules apply outside plan mode. A user deny also applies
in yolo. The extension cannot enforce a prompt against the selected host mode.
Its SSRF, file-path and resource checks still apply after host approval. See the
[host policy tests and decision](docs/plans/authority-contract.md).

The bundled [web research skill](skills/web-research/SKILL.md) describes how to
search, read pages and cite sources. It ships with source installs and archives.
The SDK has no per-call cancellation context or dedicated result trust metadata.

## Migrating an existing installation

Do not enable zot-web and terva-ext-web together. They share tool, command,
configuration and secret identities. Use `terva ext list` and host-reported
paths to identify them. Preserve the old installation and settings, disable it,
then test the replacement. To roll back, disable the replacement before enabling
the old installation. Do not infer data paths from the repository name.

This cleanup removes the previously planned 0.4.x legacy configuration period.
Before using the replacement, enter your settings in Terva and rename
`ZOT_WEB_*` overrides to `TERVA_EXT_WEB_*`. Keep `TAVILY_API_KEY` unchanged.
The `configuration_source` field is retired. Old standalone files are ignored
and preserved; there is no automatic import. Rollback uses the old installation,
not a legacy mode in this one. See [credential migration](docs/plans/credential-migration.md).

## Development and release

```sh
just ci
just ticket-check
just release-verify
just release-snapshot
```

`just ci` runs vet, formatting, race tests, subprocess conformance, the published
host driver and vendor consistency checks. `tests/host-contract` is a separate
module and downloads its dependencies on first use. After dependency changes,
run `just vendor` and commit `go.mod`, `go.sum` and `vendor/` together.

Work is tracked with `git ticket`. The [backlog](docs/plans/modernization-backlog.md)
and [release process](docs/plans/release-process.md) link the remaining work.
SDK migration, session-safe saves, host configuration and native platform tests
are complete. Installation, upgrade, rollback and first publication remain.

## License

[MIT](LICENSE). Copyright 2026 Drew Short.
