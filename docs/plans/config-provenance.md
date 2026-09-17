# Terva configuration

The extension reads settings supplied by Terva v0.137.0 and current environment
overrides. On 2026-09-16, the user requested removal of legacy configuration
support before the first release. This supersedes the earlier 0.4.x transition
period. Standalone `config.json`, `ZOT_WEB_*` and `configuration_source` no
longer affect the extension.

## Host contract

The published host's `ResolveExtensionConfig` sends stored values or non-nil
manifest defaults for declared keys. `extdriver` supplies that map at startup
and on `config_update`. The SDK does not identify which source supplied a value.

The manifest omits defaults. The application applies them after resolving
host settings, preserving explicit false, zero and empty values. A regression
test prevents manifest defaults from being added accidentally. The extension
does not inspect host configuration or authentication files.

The host omits a secret it cannot decrypt. The extension cannot distinguish
that from an unset secret, so both leave search unconfigured unless
`TAVILY_API_KEY` supplies a key.

## Precedence and validation

For nonsecret fields, `TERVA_EXT_WEB_*` overrides the host value, which overrides
the application default. Presence matters. An empty optional string clears its
value; an empty required backend, number or boolean fails validation.
`TAVILY_API_KEY` overrides the host secret, including an explicit empty value.

The host allowlist is JSON array text. It replaces the default loopback list;
`[]` removes those exceptions. `TERVA_EXT_WEB_ALLOW_LOCAL_HOSTS` appends a
comma-separated list. Entries must be hostnames, IPs or CIDRs.

JSON null, invalid types and values outside the documented limits are rejected.
The resolver does not clamp limits or fall back after an invalid override.
Errors name the field but omit its value. The [README settings table](../../README.md#settings)
records types, defaults and limits.

## Live updates

Each accepted host update creates a runtime with its own validated settings,
provider and cache. Calls already in progress retain their runtime. Later calls
use the new runtime and cannot reuse pages fetched under an old allowlist,
backend or User-Agent.

A rejected update retains working settings and sends a field-only notification.
Invalid startup configuration blocks all network tools. Missing search
credentials block search while fetch tools remain available. Environment
changes require a process restart.

Blank nonsecret form input means unset. Blank secret form input keeps the saved
secret; clearing it requires the host's clear operation. Do not infer richer
form behavior from the SDK map.

## Migration and checks

Operators must enter old values in Terva and rename environment overrides before
using the replacement. Existing files remain untouched for manual reference
and rollback with the old installation. The extension does not import keys or
write migration markers. See [credential migration](credential-migration.md).

Tests cover host and environment precedence, numeric bounds, allowlist
replacement and appending, explicit zero/false/empty values, safe error output,
retired inputs, live updates and in-flight cache ownership. Subprocess tests
prove that valid and malformed old files are ignored without changing them.
All credential fixtures are synthetic.
