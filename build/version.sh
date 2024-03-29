#!/bin/bash
set -e -x

#if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
#    DIRTY="_dirty"
#fi

COMMIT=$(git rev-parse --short HEAD)
GIT_TAG=$(git tag -l --contains HEAD | head -n 1)

if [[ -z "$DIRTY" && -n "$GIT_TAG" ]]; then
    VERSION=$GIT_TAG
else
    VERSION="${COMMIT}${DIRTY}"
fi

echo $VERSION
