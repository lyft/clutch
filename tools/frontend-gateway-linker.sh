#!/bin/bash
set -euo pipefail

REPO_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"
# Packages should be added to this list if there can only be one of them present when using Clutch as a submodule.
LINKED_PACKAGES=(
  "esbuild"
  "react"
  "react-dom"
  "react-router"
  "react-router-dom"
  "recharts"
  "@emotion/react"
  "@emotion/styled"
  "@mui/styles"
  "@mui/material"
  "@types/enzyme"
  "@types/jest"
  "@types/mocha"
  "@types/node"
  "@types/react"
  "@types/react-dom"
  "typescript"
)

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
"${YARN}" install --immutable

# Link deps from core repo.
# cd node_modules
# NODE_MODULES_DIR=$(pwd)
# for package in "${LINKED_PACKAGES[@]}"; do
#   cd "${package}"
#   "${YARN}" link
#   cd "${NODE_MODULES_DIR}"
# done

# Ensure yarn in destination directory
cd "${EXTERNAL_ROOT}"
"${REPO_ROOT}"/tools/install-yarn.sh

# # Use linked deps in consuming repo.
cd "${DEST_DIR}"
# for package in "${LINKED_PACKAGES[@]}"; do
#   "${YARN}" link "${package}"
# done

if [[ -f "yarn.lock" ]]; then
  echo "Found lockfile..."
  "${YARN}" install --immutable
else
  echo "No lockfile. Generating one..."
  "${YARN}" install
fi
"${YARN}" "${action}"
