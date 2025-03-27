#!/bin/sh
ROOT_DIR=$(git rev-parse --show-toplevel)

if command -v pre-commit &> /dev/null
then
    printf "pre-commit is already installed\n"
    pre-commit install --config $ROOT_DIR/.local/githook/pre-commit-config.yaml
else
    brew install pre-commit
    pre-commit install --config $ROOT_DIR/.local/githook/pre-commit-config.yaml
    printf "pre-commit command install done\n"
fi
