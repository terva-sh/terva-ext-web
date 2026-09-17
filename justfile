# terva-ext-web dev tasks. Run `just` to list.
set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

# Maintainer-only release-cut targets (release-cut/-verify/-publish/…).
# Optional import: the public tree ships without release.just and this
# justfile still works there.
import? 'release.just'

# Default SearXNG instance for `just configure-searxng` (override by passing a URL).
SEARXNG_URL := "http://127.0.0.1:11984"

# Target host. Stock zot support is no longer a product goal.
HOST := "terva"

default:
    @just --list

# Build the extension binary (loaded by run.sh / copied in by `just install`).
# Offline build against vendor/ — mirrors what run.sh does on first launch.
build:
    go build -mod=vendor -o terva-ext-web .
    @echo "built ./terva-ext-web"

# Refresh the committed vendor/ tree after changing dependencies.
#
# We vendor so run.sh's build-on-first-launch is a fast OFFLINE compile: Terva
# blocks its whole startup until the extension sends `hello`, and a network
# module download there would stall (or, if it hangs, freeze Terva). Re-evaluate
# this approach if vendor/ grows large (currently ~6 MB / a handful of deps) —
# at some point committing prebuilt per-platform binaries (goreleaser) becomes
# the better trade-off than carrying a big vendor tree in the repo.
vendor:
    go mod tidy
    go mod vendor
    @echo "vendor/ refreshed — commit it alongside go.mod/go.sum"

# Install through Terva without removing existing installations or credentials.
install host=HOST: build
    {{host}} ext install "$PWD"

# Set the search backend through Terva's declared configuration API.
# For a private non-loopback instance, also configure allow_local_hosts in Terva.
configure-searxng url=SEARXNG_URL host=HOST:
    {{quote(host)}} ext config web set search_backend=searxng {{quote("searxng_url=" + url)}}

# Vet + gofmt check. gofmt walks the filesystem, so exclude the vendored
# third-party tree (go vet ./... already skips vendor/ in module mode).
lint:
    go vet ./...
    @test -z "$(gofmt -l $(find . -name '*.go' -not -path './vendor/*') | tee /dev/stderr)" || { echo "gofmt issues (run \`just fmt\`)"; exit 1; }

# Format sources (excluding the vendored tree).
fmt:
    gofmt -w $(find . -name '*.go' -not -path './vendor/*')

# Run tests.
test *ARGS:
    go test ./... {{ARGS}}

# Protocol conformance: build ./terva-ext-web and drive it over stdio as the supported
# Terva host wire profile. Tagged out of the default `test` run
# because it shells out to `go build`.
conformance:
    go test -tags conformance -run Conformance -v .

# Everything the Forgejo CI gate runs (.forgejo/workflows/ci.yml mirrors this).
ci: lint
    go test -race ./...
    just conformance
    just host-contract
    go mod vendor
    git diff --exit-code -- go.mod go.sum vendor/

# Build and load into a one-off host session for manual testing.
try DIR="." host=HOST: build
    {{host}} --ext "$PWD" --cwd "{{DIR}}"

# Print the version string the binary would report, built from source.
version:
    @go run -mod=vendor . --version

# Remove build output.
clean:
    rm -f terva-ext-web terva-ext-web.exe

# Validate ticket content and detect pending repairs without changing files.
ticket-check:
    git ticket check --fix --dry-run --strict

# Test published host policy separately from the extension's offline module.
# Requires network/module cache on first run; never reads installed credentials.
host-contract:
    cd tests/host-contract && go test -race -count=1 ./...
