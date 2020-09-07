#!/bin/sh

set -e;

if ! command -v git > /dev/null; then
    echo "could not find git"
    exit 1
fi

if ! git describe --tags > /dev/null; then
    echo "no git tags exist"
    exit 1
fi

VERSION="$(git describe --tags | sed -E 's@-[[:digit:]]+-g[a-f0-9]{7}$@@')"

echo "Found latest version: ${VERSION}"

go build -ldflags="-X main.version=${VERSION}"
