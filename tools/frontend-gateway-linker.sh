#!/bin/bash
set -euo pipefail

REPO_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"
LINKED_PACKAGES=("react" "react-dom" "react-router" "react-router-dom")

EXTERNAL_ROOT="${1}"
YARN="${EXTERNAL_ROOT}/build/bin/yarn.sh"

DEST_DIR="${EXTERNAL_ROOT}/frontend"

# ensure consistent yarn versioning
cd "${REPO_ROOT}"
tools/install-yarn.sh

# default to build, can pass in start as second argument if dev is desired
action="${2:-build}"

ln -sf "${REPO_ROOT}" "${DEST_DIR}"

cd "${REPO_ROOT}/frontend"
"${YARN}" --frozen-lockfile install

# Link deps from core repo.
cd node_modules
for package in "${LINKED_PACKAGES[@]}"; do
  cd "${package}"
  "${YARN}" link
  cd ..
done

# Ensure yarn in destination directory
cd "${EXTERNAL_ROOT}"
"${REPO_ROOT}"/tools/install-yarn.sh

# Use linked deps in consuming repo.
cd "${DEST_DIR}"
for package in "${LINKED_PACKAGES[@]}"; do
  "${YARN}" link "${package}"
done

"${YARN}" --frozen-lockfile install
"${YARN}" "${action}"
