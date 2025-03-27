#!/bin/sh
cd "$(dirname "$0")"
ROOT_DIR=$(git rev-parse --show-toplevel)

if ! echo $PATH | grep -q "$(go env GOBIN)" ; then
    printf 'directory GOBIN DOES NOT exists in PATH\n'
    printf 'please add this to .zshrc\n\n'
    printf '  $ export PATH=$PATH:'"$(go env GOBIN)\n\n"
    PATH="$PATH:$(go env GOBIN)"
fi


if ! ./setup-dev-tool.sh; then
    printf 'setup dev tool fail\n'
fi

if ! ./install-deps.sh; then
    printf 'install deps fail\n'
fi

if ! ./setup-pre-commit.sh; then
    printf 'setup pre-commit fail\n'
fi

if [ ! -f "$ROOT_DIR/.env" ]; then
		cp "$ROOT_DIR/.env.template" "$ROOT_DIR/.env"
fi
