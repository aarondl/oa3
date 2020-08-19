#!/bin/sh

if ! command -v git > /dev/null; then
    echo "could not find git"
    exit 1
fi

if ! git describe --tags > /dev/null; then
    echo "no tags exist"
    exit 1
fi

VERSION="$(git describe --tags | sed --regexp-extended 's@-[[:digit:]]+-g[a-f0-9]{7}$@@')"

go build -ldflags="-X main.version=${VERSION}"
