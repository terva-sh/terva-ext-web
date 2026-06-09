# zot-web dev tasks. Run `just` to list.
set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

# Default SearXNG instance for `just configure-searxng` (override by passing a URL).
SEARXNG_URL := "http://127.0.0.1:11984"

default:
    @just --list

# Build the extension binary (name matches extension.json "exec").
build:
    go build -o zot-web .
    @echo "built ./zot-web"

# Build, then (re)install into $ZOT_HOME so the latest binary is loaded.
install: build
    #!/usr/bin/env bash
    set -euo pipefail
    name="$(basename "$PWD")"
    zot ext remove "$name" -y || true                 # -y skips the confirm; matches the dir basename
    zot ext install "$PWD"
    # `zot ext install` copies git-aware and skips .gitignore'd files — which
    # includes the built ./zot-web binary (extension.json's exec target). Copy it
    # in explicitly so the installed extension can actually run.
    line="$(zot ext list | grep -E "/${name}\$" || true)"
    [[ -n "$line" ]] || { echo "install: could not find installed dir in 'zot ext list'" >&2; exit 1; }
    dir="/${line#*/}"
    cp -f zot-web "$dir/zot-web"
    echo "copied binary -> $dir/zot-web"
    zot ext list

# Point the installed extension at a SearXNG backend (default: SEARXNG_URL).
configure-searxng url=SEARXNG_URL:
    #!/usr/bin/env bash
    set -euo pipefail
    # Resolve the installed data dir from `zot ext list` — portable across OSes
    # and a custom $ZOT_HOME. The dir is the last column; the path may contain
    # spaces, so take everything from the first '/'. URL defaults to SEARXNG_URL;
    # pass one to override (a bare host:port gets an http:// prefix).
    url="{{url}}"
    [[ "$url" == *://* ]] || url="http://$url"
    name="$(basename "$PWD")"
    line="$(zot ext list | grep -E "/${name}\$" || true)"
    [[ -n "$line" ]] || { echo "extension not installed; run \`just install\` first" >&2; exit 1; }
    dir="/${line#*/}"
    cat > "$dir/config.json" <<JSON
    {
      "search_backend": "searxng",
      "searxng_url": "$url"
    }
    JSON
    echo "configured searxng -> $url"
    echo "wrote $dir/config.json"

# Vet + gofmt check.
lint:
    go vet ./...
    @test -z "$(gofmt -l . | tee /dev/stderr)" || { echo "gofmt issues (run \`just fmt\`)"; exit 1; }

# Format sources.
fmt:
    gofmt -w .

# Run tests.
test *ARGS:
    go test ./... {{ARGS}}

# Build and load into a one-off zot session for manual testing.
try DIR=".": build
    zot --ext "$PWD" --cwd "{{DIR}}"

# Remove build output.
clean:
    rm -f zot-web
