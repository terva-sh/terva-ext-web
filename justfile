# zot-web dev tasks. Run `just` to list.
set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

default:
    @just --list

# Build the extension binary (name matches extension.json "exec").
build:
    go build -o zot-web .
    @echo "built ./zot-web"

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
