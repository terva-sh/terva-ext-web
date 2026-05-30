# zot-web

A [zot](https://github.com/patriceckhart/zot) extension that gives the agent
web access through two LLM-callable tools:

- **`web_search(query, count?)`** — ranked results (title, URL, snippet).
- **`web_fetch(url, max_chars?)`** — a page's main content as text.

Single static Go binary, no runtime dependencies. It implements the zot
extension wire protocol directly (no dependency on the zot module).

> **Status: v0 scaffold.** Search (Tavily + SearXNG) and the SSRF-guarded
> fetcher are functional; HTML extraction is a heuristic placeholder pending
> `go-readability` + `html-to-markdown` (see Roadmap). Design rationale lives in
> the zot repo at `docs/plans/web-tools-extension-research.md`.

## Build

```bash
go build -o zot-web .
```

## Install into zot

```bash
zot ext install /path/to/zot-web      # copies into $ZOT_HOME/extensions/web/
# or, for one session straight from the working copy:
zot --ext /path/to/zot-web
```

The directory must contain `extension.json` (pointing at the `./zot-web`
binary you built) — already included here.

## Configure

Settings come from `config.json` in the extension's data dir
(`$ZOT_HOME/extensions/web/config.json`), with environment variables taking
precedence. Minimum to get search working with the default Tavily backend:

```bash
export TAVILY_API_KEY=tvly-...
```

Or switch to a self-hosted SearXNG instance (no key, private):

```jsonc
// $ZOT_HOME/extensions/web/config.json
{
  "search_backend": "searxng",
  "searxng_url": "http://localhost:8888",
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
| `allow_local_hosts` | `ZOT_WEB_ALLOW_LOCAL_HOSTS` (comma-sep) | — | SSRF escape hatch (see below) |

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

- [ ] Replace the heuristic HTML extractor with `go-shiori/go-readability` +
      `JohannesKaufmann/html-to-markdown`.
- [ ] More search backends (Brave, Serper, Exa) behind the same interface.
- [ ] Optional result caching in the data dir.
- [ ] Optional JS rendering fallback (e.g. Jina Reader) — deferred for now.
- [ ] Release binaries (goreleaser) so `zot ext install <git-url>` needs no
      local build.
