#!/bin/sh
# Stamp the archive manifest without changing the source manifest.
set -eu
version=${1:?release version required}
case "$version" in ''|*[!0-9A-Za-z.+-]*) echo 'invalid release version' >&2; exit 1 ;; esac
mkdir -p .release-package
sed "s/\"version\": \"[^\"]*\"/\"version\": \"$version\"/" extension.json > .release-package/extension.json
