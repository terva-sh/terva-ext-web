# Untrusted web content

Fetched pages control their body, title, descriptions and image attributes.
Search results also contain external text. Any of it can include instructions
intended to redirect the model. Network and file guards limit tool effects;
they cannot make page content trustworthy.

## Output formatting

The extension flattens external values where they appear in its own metadata:

- `web_fetch` titles use `titleLine`, which collapses whitespace and caps the
  result at 300 runes. `/web-cache` uses the same title.
- `web_search` titles and snippets use `oneLine`.
- `web_images` flattens dimensions, alt text and captions.
- `web_links` flattens link text and resolves URLs before output.

This prevents a title or attribute containing a newline from forging another
metadata line or numbered result. Tests in `internal/fetch/sanitize_test.go`
and the search formatter check hostile inputs and normal output.

## No inline trust markers

The inherited proposal to wrap tool text in nonce-bearing trust tags was
rejected. Such tags would still be text for the model to interpret. Escaping
a closing tag would not prevent the model from following instructions inside
it. The supported SDK provides no dedicated result trust metadata.

Use a host-level trust mechanism if Terva adds one. Until then, fetched bodies,
snippets and captions remain untrusted. The bundled research skill tells the
agent to treat page instructions as source content, but that guidance is not
an enforcement boundary.

## Enforced limits

The fetcher validates destinations, redirects, ports and resource limits.
Workspace saves reject escapes, symlinks and Git metadata paths, and require
explicit overwrite. A call captures its workspace before fetching.

Terva controls tool approval. All six tools use `network-read`; the two writer
tools also request `ask` in the manifest. Yolo allows those requests, while an
explicit user deny still applies. The [authority contract](authority-contract.md)
records the tested policy. Approval does not disable application guards.
