#!/bin/sh
cd "$(dirname "$0")"
ROOT_DIR=$(git rev-parse --show-toplevel)

# bump version in VERSION file by simply take the first argument e.g. v1.0.0 then write it to VERSION file
# if no frist argument provided, it will show error message and exit
# if the first argument is provided, it will validate the version format by using regex pattern ^v[0-9]+\.[0-9]+\.[0-9]+$
# if the version format is invalid, it will show error message and exit
# if the version format is valid, it will write the version to VERSION file then commit the change to git repository with message "bump version to $1" and exit

echo "Bumping version to $1"

if [ -z "$1" ]; then
		echo "Please provide version number. e.g. v1.0.0"
		exit 1
fi

VERSION_REGEX="^v[0-9]+\.[0-9]+\.[0-9]+$"
if ! [[ "$1" =~ $VERSION_REGEX ]]; then
		echo "Invalid version format. Please use v0.0.0 v(major.minor.patch) format"
		exit 1
fi

# write version to VERSION file without newline
echo $1 > $ROOT_DIR/VERSION

# git add $ROOT_DIR/VERSION
# git commit -m "bump version to $1"
# exit 0
