#!/bin/bash
set -euo pipefail

manager=${MANAGER:-yarn}

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
YARN="${EXTERNAL_ROOT}/build/bin/yarn.sh"

if [ "${manager}" = "yarn" ]; then
  BUILD="${YARN}"
else
  BUILD="pnpm"
fi

DEST_DIR="${EXTERNAL_ROOT}/frontend"

# ensure consistent yarn/pnpm versioning
cd "${REPO_ROOT}"

if ! command -v yarn *> /dev/null; then
  echo "Installing yarn..."
  tools/install-yarn.sh
fi

if ! command -v pnpm *> /dev/null; then
  echo "Installing pnpm..."
  tools/install-pnpm.sh
fi

# default to build, can pass in start as second argument if dev is desired
action="${2:-build}"

ln -sf "${REPO_ROOT}" "${DEST_DIR}"

cd "${REPO_ROOT}/frontend"
"${YARN}" --frozen-lockfile install

# Link deps from core repo.
cd node_modules
NODE_MODULES_DIR=$(pwd)
for package in "${LINKED_PACKAGES[@]}"; do
  cd "${package}"
  "${BUILD}" --global link
  cd "${NODE_MODULES_DIR}"
done

# Ensure yarn/pnpm in destination directory
cd "${EXTERNAL_ROOT}"

if ! command -v yarn *> /dev/null; then
  echo "Installing yarn..."
  "${REPO_ROOT}"/tools/install-yarn.sh
fi

if ! command -v pnpm *> /dev/null; then
  echo "Installing pnpm..."
  "${REPO_ROOT}"/tools/install-pnpm.sh
fi

# Use linked deps in consuming repo.
cd "${DEST_DIR}"
for package in "${LINKED_PACKAGES[@]}"; do
  "${BUILD}" link --global "${package}"
done

if [[ -f "pnpm-lock.yaml" ]]; then
  echo "Found pnpm lockfile..."
  pnpm install --frozen-lockfile
elif [[ -f "yarn.lock" ]]; then
  echo "Found yarn lockfile..."
  "${YARN}" install --frozen-lockfile
else
  echo "No lockfile. Generating one..."
  "${YARN}" install
fi
"${BUILD}" "${action}"
