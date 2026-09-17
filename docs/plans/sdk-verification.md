# Published Terva SDK verification, 2026-09-16

This record selected `terva.sh/terva v0.137.0` for the SDK migration. This is a
published module, retrieved into an empty temporary module cache through
`https://proxy.golang.org` with `sum.golang.org` verification. The proxy's
`@latest` query also resolved to v0.137.0 on this date. No local replace,
repository go.mod edit, credential file, or sibling worktree change was used.

| Evidence | Value |
| --- | --- |
| Origin | https://github.com/terva-sh/terva |
| Tag | refs/tags/v0.137.0 |
| Origin commit | a4231906ee8365179998e4f6cb581cd3d5a10a45 |
| Module checksum | h1:H9niHYCahK8NeOngWQnbGTcXb77YOjNCvwepbKpy4I8= |
| go.mod checksum | h1:JPb85Jz4dgUpiNYJf2d+TgIWdkgXp7HNI+JvX1HvbOc= |
| Published module Go requirement | 1.27.0 |
| Verification toolchain | go1.27.1 linux/amd64 |
| SDK wire revision | 6 |
| Public version timestamp | 2026-09-15T05:12:43Z |

## Supported host and capabilities

Terva v0.137.0 is the initial supported host and SDK version. It is both the selected published SDK and the current public
version; floor/current tests therefore use the same version initially, and
must not be reported as two distinct version results. At the time of this probe, real-host conformance had not run. The later
validation records are linked below.

The correctness requirement is protocol 2 for ordered session identity;
the implementation uses `RequireProtocol(2)`. Earlier Terva versions that speak
protocol 2 may work, but are outside the initial verified support baseline.
Do not equate the SDK's protocol 6 constant with a need to reject all pre-v6
hosts: only runtime secret-broker operations would introduce that requirement.
The local CLI inspected during this probe reported Terva 0.135.1. Later tests
use the published v0.137.0 host as a separate fixture.

| Published API | Contract / migration decision |
| --- | --- |
| ext.New, Tool, Command, Run | Replace handwritten framing/dispatch; preserve web identity, tool/command names and stdout JSON |
| ext.ToolResult, TextResult, TextErrorResult, ImageBytes | Text/image results exist; image content encodes bytes, not a trust/provenance attribute |
| OnSession, Host().CWD, Session | Session fields update on the ordered reader path; snapshot Host().CWD once per download handler, preserving the save-race fix |
| Host().Emits, optional event subscriptions | Use advertised events for optional behavior; do not infer all features from a version string |
| Config(), OnConfig, Config.Has | Resolved values and defaults are present; Has is only map membership. The separate import/provenance decision remains necessary |
| WithAuthority | One authority string; retain network-read and manifest ask defaults for writers until the security review proves an alternative |
| Sequential | Optional serial lane; network reads remain concurrent by default |
| Essential, WithDisplay | Optional presentation/discovery hints; do not create hard protocol floors for ignored display metadata |
| OnCompaction and RefreshContext | APIs exist; optional context changes must follow documented boundary/capability rules |
| HostToolCall | Request/response requires protocol 3; not needed for the initial direct guarded saves |
| Withdraw/Restore | Visibility support is detected at protocol 4, using session-scoped handles and runtime readiness fallback |
| Secret operations | Runtime broker requires protocol 6; config-secret values are already delivered separately and do not justify introducing the broker |
| ToolHandler | JSON arguments to ToolResult, with no per-call context. SDK adoption does not supply cooperative tool cancellation |
| ReadFrame | SDK uses bounded framing with oversized/malformed-frame recovery; exercise it in conformance rather than retaining the old scanner |
| bootstrap | Launcher frame before hello; older hosts can reject it. This is not an SDK operation and support cannot be handshaken before sending it |

Evidence is in the downloaded v0.137.0 module, principally
`packages/agent/ext/ext.go`, `packages/agent/extproto/extproto.go`, and
`docs/extensions.md`. The separately inspected local Terva checkout remains at
`7f754b9bb7c6754284dc4f3cc5fa6a96525de49d`; its untracked ticket was untouched.

## Build and vendor findings

An isolated module importing the actual SDK built successfully after
`go mod tidy`, `go mod verify` (all modules verified), and `go mod vendor`.
The probe required the SDK plus these transitive modules:

- filippo.io/age v1.3.1
- filippo.io/hpke v0.4.0
- golang.org/x/crypto v0.54.0
- golang.org/x/sys v0.47.0

The probe's vendor files totalled 735,453 bytes. This is not a prediction of
this extension's final vendor size: Go's module graph can also raise existing
x/image, x/text and other selected versions. Inspect and record the real
combined-module changes in the migration commit.

The SDK probe completed hello (including min_protocol=2), tool registration,
ready and shutdown with JSON-only stdout on Linux amd64. Offline vendored,
CGO-disabled builds passed for linux amd64/arm64, darwin amd64/arm64 and windows
amd64. Only Linux amd64 was executed. A probe exchanged frames with a test process;
it is not a full Terva-host installation or extension behavior test.

The later SDK batch raised go.mod, launcher instructions and CI to Go 1.27.
Published host-driver tests and five-target native validation now supplement
this initial probe. See [SDK conformance](sdk-conformance.md) and
[platform validation](platform-validation.md) for the tested candidates.

## Reproducing the public retrieval

Run in a newly created temporary directory, not the repository. Replace
`/absolute/path/to/go1.27/bin` with the verified toolchain directory.

```sh
export PATH="/absolute/path/to/go1.27/bin:$PATH"
mkdir -p module-cache build-cache
export GOMODCACHE="$PWD/module-cache" GOCACHE="$PWD/build-cache"
export GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org
export GOPRIVATE= GONOPROXY= GONOSUMDB= GOFLAGS= GOENV=off GOTOOLCHAIN=local
go mod download -json terva.sh/terva@v0.137.0
go list -m -json terva.sh/terva@latest
```

The first command must report the checksums above. The second is time-sensitive:
if a newer public version appears, keep the pinned evidence and evaluate any
upgrade deliberately. Reproduce the small SDK probe by requiring this version,
registering one text tool with `WithAuthority("network-read")`, subscribing with
OnSession/OnConfig, calling RequireProtocol(2), and running its standard stdin/
stdout loop. Build with `-mod=vendor` and GOPROXY=off for the offline checks.

Alternatives: copying the older extension-template pin misses the assessed
APIs; using a local replace does not prove publication; widening the supported
host range without versioned runtime tests creates an unsupported promise.
The published v0.137.0 baseline has the required APIs and a verifiable origin,
so those alternatives are unnecessary.
