# Configuration provenance and legacy transition

Decision for published Terva v0.137.0. This specifies the configuration and
secret batches; SDK adoption alone does not change configuration behavior.

## Available contract

In the published module, `packages/agent/build/extmcp.go`:
`ResolveExtensionConfig` selects a stored value, otherwise a non-nil manifest
default, and delivers only schema-declared keys. `extdriver/extconfig.go`
delivers that map on hello and config_update. `ext.Config.Has` tests map
membership; it carries no source/provenance. An explicit value equal to a
manifest default is indistinguishable from an unset field with that default.

**Omit manifest defaults for all migrated fields.** Keep defaults in the
application and describe them in form help. With this schema, a delivered
nonsecret key represents a stored host value, including false, zero or an
explicit empty value. Test the manifest has no defaults so a later schema edit
cannot silently invalidate precedence. Do not compare values against defaults
or read the host's config/auth files to guess provenance.

Exception: the resolver omits a secret it cannot decrypt. Absent and
undecryptable secrets are indistinguishable to the extension. Do not solve this
by using Config.Has or reading encrypted storage directly.

## Transition mode and precedence

Add an optional `configuration_source` select: `legacy` or `host`, no manifest
default; application default `legacy` preserves installed configurations.

- `legacy`: nonsecret host values override legacy file values per field. The
  Tavily credential remains legacy/env-owned; a host secret is not used until
  the operator explicitly selects `host`. Explain this in setup help.
- `host`: do not read the legacy configuration file at all. Host values override
  application defaults. A missing/undecryptable Tavily key stays missing unless
  supplied by an environment override; never resurrect a legacy credential.
- The mode itself comes from the host, with new then legacy environment
  overrides if supported. Invalid mode rejects the candidate settings.

This is an explicit opt-in import: use the host configuration form to enter
legacy choices, then select `host` and verify retrieval. The extension never
copies credentials into host storage, writes migration markers, or rewrites or
deletes old files. Repeating the form settings is idempotent. Rollback selects
`legacy` or disables the replacement and re-enables the old installation; old
files remain available. No automatic credential import is promised.

For scalar fields select the first *present* environment variable in the new
`TERVA_EXT_WEB_` namespace, then the corresponding `ZOT_WEB_` variable, then
host, then legacy (only in legacy mode), then application defaults. Empty new
environment values are explicit: optional strings can clear, while required
backend/numeric/bool values reject instead of falling through to legacy.
TAVILY_API_KEY remains the provider credential override, above host/legacy;
do not add competing renamed credential variables.

| Field | Host type | Application default / value rules |
| --- | --- | --- |
| search_backend | select | tavily; only tavily or searxng |
| searxng_url | string | empty means unconfigured; configured URL must be valid http(s) |
| tavily_api_key | secret | empty means unconfigured; host only in host mode; TAVILY_API_KEY first |
| user_agent | string | empty selects terva-ext-web/version; browser retains its existing meaning |
| allow_local_hosts | string | JSON array text, not comma list; default localhost, 127.0.0.1, ::1 |
| fetch_max_bytes | int | 2097152; positive, at most 33554432 |
| fetch_image_max_bytes | int | 5242880; positive, at most 20971520 |
| fetch_timeout_sec | int | 25; positive, at most 60 |
| fetch_inline_images | bool | false; explicit false overrides legacy true |
| fetch_cache_ttl_sec | int | 600; 0 means no expiry, maximum 3600 |
| fetch_cache_max_entries | int | 32; 0 disables cache, maximum 128 |
| fetch_cache_max_bytes | int | 67108864; 0 retains historical unbounded-byte meaning with entry cap; maximum 268435456 |

Host allowlist text is decoded to an array; [] explicitly removes loopback.
Legacy JSON arrays remain arrays. Both REPLACE the lower-priority list.
Environment allowlists retain documented additive behavior: choose the new
variable if present, otherwise legacy, split its comma list and append to the
resolved list. Never append both namespaces. Empty env adds nothing; use host
[] for no local access. Validate entries and explain this exception clearly.

Missing file uses defaults; unreadable or malformed first existing legacy file
fails visibly rather than falling through to another file. Legacy file order is
host-reported DataDir then ExtensionDir; never infer paths from a repo name.
JSON null is invalid for a present migrated field, not missing. Out-of-range
new/host values reject rather than silently clamp. Any intentional change from
legacy clamping must be documented and covered in migration tests.

## Runtime and credential behavior

Build a validated immutable runtime (config, provider and fresh fetch/cache
client) for each accepted config update. Atomically replace the runtime so
in-flight calls finish with the configuration they captured; later calls use
new policy and cannot reuse an earlier allowlist/backend/User-Agent cache.
Rejected updates retain the previous working runtime and send a generic,
field-naming notification without values. Initial invalid settings block every
network tool, not only search. Missing search credentials only block search;
fetch tools remain available. Never echo raw parse errors containing secrets,
URLs with userinfo, configuration maps or credential values.

Keep data_secrets conservative while old files may contain a key. A masked
secret field is not proof of encryption; encryption depends on host setup.
Runtime secret broker APIs are unnecessary for a user-provided key and would
unnecessarily impose protocol 6. This decision needs no new host API.

## Synthetic implementation fixtures

Use temporary directories and fake credentials, never copied user config.
Each case must assert effective values and safe diagnostics, not print keys.

| Input | Expected result |
| --- | --- |
| no host map/file/env | application defaults; search unconfigured |
| legacy backend searxng, host key absent | legacy backend preserved |
| same legacy, host backend tavily (equal to default) | explicit tavily wins |
| legacy true/TTL 900, host false/0 | false and no-expiry TTL preserved |
| host null, malformed bool/int, limit overflow | candidate rejected, no partial swap |
| legacy UA custom, explicit host empty | default UA; no legacy resurrection |
| empty host form submission | host omits nonsecret key; legacy/default fallback (not explicit empty) |
| host allowlist [] text, legacy private host | empty list, loopback blocked |
| host list plus both env namespaces | append only selected new namespace |
| invalid new env with valid legacy env | reject, do not silently fall through |
| malformed/unreadable data file and valid install file | report data-file failure in legacy mode |
| same files in host mode | no legacy read; valid host config accepted |
| legacy fake key plus saved host fake key, legacy mode | legacy key remains active pending opt-in |
| switch host mode, host fake key present | host key used; repeated switch unchanged |
| host mode, host key absent/undecryptable, legacy key exists | missing credential, never use legacy key |
| host secret form left blank | host retains saved secret; not a clear operation |
| provider env explicitly empty | search unconfigured, no host/legacy fallback |
| changed allowlist/backend/UA during blocked fetch | old call finishes its snapshot; new call uses fresh client/cache |
| failed update during valid operation | old runtime stays active; notification names field, no values |
| rollback to legacy or old installation | preserved file usable; never run both web identities together |

A form cannot express an empty nonsecret value (blank means unset), and blank
secret input means keep. Use the explicit host mode for dropping legacy
fallback, [] for an empty allowlist, and documented host clear operations for
stored credentials. Do not promise form semantics the host does not support.
