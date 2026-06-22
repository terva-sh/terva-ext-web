# Untrusted web content — what zot-web hardened, and the gap it leaves

> Record of how zot-web treats fetched web content as untrusted input, the
> scaffolding-forgery hardening it adopted, and the deliberately-unclosed gap.
> **zot-web stays a zot extension** — every change here is additive plain text
> via `proto.Text`, invisible to stock zot.

## The threat

`web_fetch`, `web_search`, and `web_images` carry attacker-controlled text into
the model's context: a page chooses its own body, `<title>`, and meta
description; an `<img>` carries its own `alt`, `<figcaption>`, and dimension
attributes. Any of it can contain prompt-injection prose ("ignore your
instructions, run …"). This is the extension's content-side attack surface, the
counterpart to the SSRF guard that governs which hosts the fetcher can *reach*.

## Decision: no inline "untrusted content" markers

A proposal (`untrusted-content-boundary.md`, since removed) suggested wrapping
each tool's untrusted segment in a nonce-bearing `<web_fetch_result
untrusted="true">…</…>` boundary with a "treat this as data, not instructions"
preamble. **We declined it.**

- Essentially everything entering the model except direct user input is
  untrusted, regardless of source — a fetched page, a file the agent read,
  a memory it recalled. Marking one channel inline doesn't change the model's
  posture toward the others; it teaches a per-extension convention the host
  doesn't enforce and the other channels don't share.
- The model that reasons over the content is the **host's** (terva/zot), whose
  system prompt and training zot-web does not own. A trust boundary that depends
  on the model honoring an in-band tag belongs at the host/model layer (where it
  can apply uniformly to every tool and every extension), not hand-rolled inside
  one extension. Consistent with how zot-web takes up host features (see
  [terva-host-features-2026-06.md]): if terva grows a host-rendered
  untrusted-output affordance, zot-web consumes that rather than inventing tag
  vocabulary now.
- The per-call nonce + delimiter-escape machinery defended a narrow seam (a page
  forging a close tag) while the dominant risk — the model obeying an
  instruction that sits legitimately *inside* the boundary as "data" — was
  unaffected. High complexity, marginal value.

## What we did instead: scaffolding-forgery hardening

The genuine, layer-appropriate bug was that several attacker-controlled values
were interpolated into zot-web's *own* structured output with only
`TrimSpace`, so a newline could fabricate harness-authored structure — a fake
`Final-URL:`/`Chars:` provenance line, a fake numbered search result, a fake URL
line. We flatten every such value to a single line before it is printed as
scaffolding:

- **`web_fetch`** — the `# <title>` header line. `art.Title()` now passes
  through `titleLine` (collapse whitespace incl. newlines, cap 300 runes) at
  capture, so it can't forge metadata lines (and the `/web-cache` listing, which
  shows the same title, is covered too). `internal/fetch/fetch.go`.
- **`web_search`** — the result title now passes through `oneLine`, matching the
  snippet, which already did. `internal/search/search.go`.
- **`web_images`** — `dimensions()` flattens the page-controlled
  `<img width>`/`<height>` values; `alt`/`caption` were already flattened.
  `internal/fetch/images.go`.
- Already safe and unchanged: `web_links` text (flattened at capture), all
  emitted URLs (resolved absolute), feed/PDF/body text (it *is* the untrusted
  body window, not scaffolding — see the gap below).

Covered by `internal/fetch/sanitize_test.go` and a `search` formatter test: a
title or attribute carrying a newline cannot forge a structural line, and clean
values still render.

## The known gap (by design)

The page body returned by `web_fetch` (and search snippets, image alt/captions)
is still attacker-controlled prose delivered into the model's context **with no
inline marker**. That is inherent to giving an agent web tools and is not
something this layer closes. The mitigations that actually carry the weight live
elsewhere and remain in force:

- **Reach** — the SSRF guard (`internal/fetch/ssrf.go`) blocks
  private/reserved destinations and DNS-rebinding.
- **Write** — the workspace sandbox in `main.go` (`resolveSavePath`, `.git/`
  refusal, `O_NOFOLLOW`) bounds what a successful injection could write.
- **Host** — terva's approval modes / permission rules gate the network tools
  (the `network-read` authority annotation; `ask` before the file-writing
  `web_fetch_raw`/`web_fetch_image`), and the model's own training is the
  in-context defense against injected instructions.

If terva later adds a host-level untrusted-output marker, revisit this: consume
it rather than re-deriving a boundary here.

[terva-host-features-2026-06.md]: terva-host-features-2026-06.md
