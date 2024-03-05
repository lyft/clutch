#!/bin/bash
set -euo pipefail

REPO_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"

# Packages should be added to this list if there can only be one of them present when using Clutch as a submodule.
LINKED_PORTAL_PACKAGES=(
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
)

LINKED_FILE_PACKAGES=(
  "typescript"
)

COMBINED_PACKAGES=(
  "${LINKED_PORTAL_PACKAGES[@]}"
  "${LINKED_FILE_PACKAGES[@]}"
)

EXTERNAL_ROOT="${1}"
YARN=yarn

DEST_DIR="${EXTERNAL_ROOT}/frontend"
YALC_STORE_FOLDER="${DEST_DIR}/.yalc"

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

# Link deps from core repo.
cd node_modules
NODE_MODULES_DIR=$(pwd)
for package in "${COMBINED_PACKAGES[@]}"; do
  cd "${package}"
  yalc publish --no-scripts --push --store-folder="${YALC_STORE_FOLDER}" --quiet
  cd "${NODE_MODULES_DIR}"
done

# Ensure yarn in destination directory
cd "${EXTERNAL_ROOT}"
"${REPO_ROOT}"/tools/install-yarn.sh

# # Use linked deps in consuming repo.
cd "${DEST_DIR}"
echo "Linking & Setting resolutions..."
for package in "${LINKED_PORTAL_PACKAGES[@]}"; do
  yalc link "${package}" --pure --store-folder="${YALC_STORE_FOLDER}" --quiet
  npm pkg set resolutions.${package}="portal:.yalc/${package}"
done

for package in "${LINKED_FILE_PACKAGES[@]}"; do
  yalc link "${package}" --pure --store-folder="${YALC_STORE_FOLDER}" --quiet
  npm pkg set resolutions.${package}="file:.yalc/${package}"
done

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
