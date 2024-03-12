#!/bin/bash
set -euo pipefail

REPO_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"

EXTERNAL_ROOT="${1}"
YARN=yarn

DEST_DIR="${EXTERNAL_ROOT}/frontend"

# ensure consistent yarn versioning
cd "${REPO_ROOT}"
tools/install-yarn.sh

# default to build, can pass in start as second argument if dev is desired
action="${2:-build}"

ln -sf "${REPO_ROOT}" "${DEST_DIR}"

cd "${REPO_ROOT}/frontend"
{
  "${YARN}" install --immutable
} || {
  echo "${REPO_ROOT}/frontend/yarn.lock would be modified by install. Please run yarn install in the frontend directory and commit the changes."
  exit 1
}

# # Use linked deps in consuming repo.
cd "${DEST_DIR}"

if [[ -f "yarn.lock" ]]; then
  echo "Found lockfile..."
  {
    "${YARN}" install --immutable
  } || {
    echo "${DEST_DIR}/yarn.lock would be modified by install. Please run yarn install in the frontend directory and commit the changes."
    exit 1
  }
else
  echo "No lockfile. Generating one..."
  "${YARN}" install
fi
"${YARN}" "${action}"
