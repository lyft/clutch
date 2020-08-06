#!/bin/bash

VERSION="0.0.0-$(git log -1 --format=%cd --date=format:'%Y%m%d%H%M%S')-$(git rev-parse --short=12 HEAD)"
PACKAGE=$1

if yarn info "@clutch-sh/${PACKAGE}" | grep -q "$VERSION"; then
  echo "Found existing version of ${PACKAGE}@${VERSION}"
  exit 0
fi

(
  cd "$PWD" && yarn publish --new-version="0.0.0-beta.$(git log -1 --format=%cd --date=format:'%Y%m%d%H%M%S')" --access public --no-git-tag-version $1
)