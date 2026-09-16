#!/usr/bin/env bash
# Local packaging only. Publication is a separate, unconfigured batch.
set -euo pipefail
cd "$(dirname "$0")/.."
case "${1:-}" in
  verify) goreleaser check ;;
  snapshot) goreleaser release --snapshot --clean ;;
  *) echo 'Usage: scripts/release.sh {verify|snapshot}. Publication and legacy cut/* flows are disabled; see docs/plans/release-process.md.' >&2; exit 1 ;;
esac
