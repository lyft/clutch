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
shift 2

ln -sf "${REPO_ROOT}" "${DEST_DIR}"

cd "${REPO_ROOT}/frontend"
if [[ -f "yarn.lock" ]]; then
  echo -e "\nFound lockfile in ${REPO_ROOT}/frontend, running install...\n"
  {
    "${YARN}" install --immutable
  } || {
    echo -e "\n${REPO_ROOT}/frontend/yarn.lock would be modified by install. Please run yarn install in the frontend directory and commit the changes.\n"
    exit 1
  }
else
  echo -e "\nNo lockfile found in ${REPO_ROOT}/frontend - Generating one...\n"
  "${YARN}" install
fi


cd "${EXTERNAL_ROOT}"
"${REPO_ROOT}"/tools/install-yarn.sh

cd "${DEST_DIR}"

if [[ -f "yarn.lock" ]]; then
  echo -e "\nFound lockfile in ${DEST_DIR}, running install...\n"
  {
    "${YARN}" install --immutable
  } || {
    echo -e "\n${DEST_DIR}/yarn.lock would be modified by install. Please run yarn install in the frontend directory and commit the changes.\n"
    exit 1
  }
else
  echo -e "\nNo lockfile found in ${DEST_DIR} - Generating one...\n"
  "${YARN}" install
fi

ROOT_NODE_MODULES_DIR="../clutch/frontend/node_modules"
cd "${DEST_DIR}/node_modules"
NODE_MODULES_DIR=$(pwd)
echo -e "\nLinking packages...\n"
for package in "${LINKED_PACKAGES[@]}"; do
  echo "${package}: Linking..."
  BASE=$(echo "${package}" | cut -d "/" -f 1)
  SUB=$(echo "${package}" | cut -d "/" -f 2)

  if [[ "${BASE}" != "${SUB}" ]]; then
    rm -rf "${BASE:?}/${SUB}"
    ln -s -f -F "../${ROOT_NODE_MODULES_DIR}/${BASE}/${SUB}" "${BASE}/${SUB}"
    if [ -d "${BASE}/${SUB}/bin" ]; then
      cd "${NODE_MODULES_DIR}/.bin"
      for binary in "${BASE}/${SUB}/bin"/*; do
        ln -s -f "../${BASE}/${SUB}/bin/$(basename "${binary}")" .
      done
      cd "${NODE_MODULES_DIR}"
    fi
  else 
    rm -rf "${BASE}"
    ln -s -f -F "${ROOT_NODE_MODULES_DIR}/${BASE}" "${BASE}"
    if [ -d "${BASE}/bin" ]; then
      cd "${NODE_MODULES_DIR}/.bin"
      for binary in "${BASE}/bin"/*; do
        ln -s -f "../${BASE}/bin/$(basename "${binary}")" .
      done
      cd "${NODE_MODULES_DIR}"
    fi
  fi
done

if [[ "${action}" != "install" ]]; then
  echo -e "\nRunning yarn ${action} ${@}...\n"
  "${YARN}" "${action}" "${@}"
else 
  echo -e "\nCompleted install and linking!"
fi
