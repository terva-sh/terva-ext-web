# Credential migration evidence

The default-free tavily_api_key secret field is consumed only in explicit
host configuration mode. Legacy mode retains old file/environment ownership;
TAVILY_API_KEY overrides either mode, including explicit empty. This avoids
interpreting a missing (possibly undecryptable) host field as permission to
resurrect a legacy credential. Import is operator entry through the host form
followed by host-mode opt-in; retry is idempotent and rollback preserves files.

No credential file is changed or deleted and no key rotated. data_secrets is
true because old config files may remain. Host encryption is optional setup:
published v0.137.0 docs/extensions.md documents `terva secret init`, while the
host form seals only when encryption is configured. A secret field masks UI
input; it does not itself prove encrypted storage. Runtime broker is for
runtime-acquired credentials and is not needed here; min_protocol remains 2.

Synthetic tests cover legacy/host duplicate values, repeated opt-in, env
precedence and explicit empty, omitted/undecryptable host-field behavior,
rollback and unchanged legacy fixture bytes. Subprocess tests capture stdout
and stderr for configured and malformed secret values without calling a real
provider. Local fake-provider tests echo a synthetic credential in errors and
result fields; neither may expose it. No real credential files were inspected.

Provider-response error snippets were removed because a provider can echo its
authorization input. Status-based guidance remains, and configured-key echoes
are redacted from successful result fields and printable transport errors.
Underlying error identity is preserved for SSRF/timeout inspection. This is a
credential-output safeguard, not a claim of generic result trust metadata.
