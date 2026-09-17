# Credential migration

Terva's `tavily_api_key` secret field supplies the Tavily key.
`TAVILY_API_KEY` overrides it; an explicitly empty variable disables the key.
A missing or undecryptable host secret never falls back to an old file.

The user retired legacy configuration support on 2026-09-16. This replaces the
previous host-mode opt-in process. Before switching installations, enter the
settings and key through Terva, then verify search. Do not enable the old and
new extensions together because both use the `web` identity.

Standalone configuration files, `ZOT_WEB_*` variables and
`configuration_source` are ignored. The extension does not read, import,
rewrite or delete old credential files, and it does not rotate keys. Preserve
the old installation and settings until the replacement passes its smoke test.
Rollback disables the replacement before re-enabling the old installation.
There is no legacy mode in the replacement.

The manifest keeps `data_secrets: true` because old files may remain in its
data directory. A masked secret field does not establish encryption at rest.
Terva documents `terva secret init`; follow its backup and recovery instructions
before enabling encryption. This is host setup, not an extension side effect.
The extension needs no runtime secret broker and retains minimum protocol 2.

Tests cover host/environment precedence, explicit empty credentials, missing
host secrets and unchanged old files. Subprocess tests inspect stdout and
stderr for synthetic-key leaks. Fake providers echo fixture keys in response
fields and errors to exercise redaction. No installed credentials are used.

Provider errors omit raw response bodies because a provider may echo request
credentials. Status-based guidance remains. Result fields and printable
transport errors redact the configured key while preserving error identity for
SSRF and timeout handling. This does not make retrieved content trusted.
