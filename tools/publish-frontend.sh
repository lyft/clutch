#!/bin/bash

VERSION="1.0.0-beta.$(git log -1 --format=%cd --date=format:'%Y%m%d%H%M%S')"
PACKAGE=$1

if yarn info "@clutch-sh/${PACKAGE}" | grep -q "$VERSION"; then
  echo "Found existing version of ${PACKAGE}@${VERSION}"
  exit 0
fi

yarn publish --new-version="${VERSION}" --access public --no-git-tag-version