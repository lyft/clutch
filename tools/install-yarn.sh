#!/bin/bash

YARN_VERSION="1.22.5"
ROOT_DEST_DIR="$PWD/frontend"
DEST_DIR="$ROOT_DEST_DIR/.yarn/releases"
DEST_FILE="${DEST_DIR}/yarn-${YARN_VERSION}.js"
WRAPPER_DEST_DIR="${PWD}/build/bin/"
WRAPPER_DEST_FILE="${WRAPPER_DEST_DIR}/yarn.sh"

if [[ ! -d "${ROOT_DEST_DIR}" ]]; then
  echo "Could not find frontend directory. Ensure you're running this script from the root of your project."
  exit 1
fi

if [[ ! -f "${DEST_FILE}" ]]; then
  echo "Downloading yarn v${YARN_VERSION} to build environment..."
  mkdir -p "${DEST_DIR}"
  curl -sSL -o "${DEST_FILE}" \
    "https://github.com/yarnpkg/yarn/releases/download/v${YARN_VERSION}/yarn-${YARN_VERSION}.js"
fi

# Install a wrapper script in build/ that executes yarn if it doesn't exist already.
WRAPPER_SCRIPT="#!/bin/bash\nnode \"${DEST_FILE}\" \"\$@\"\n"
if [[ ! -f "${WRAPPER_DEST_FILE}" || $(< "${WRAPPER_DEST_FILE}") != $(printf "%b" "${WRAPPER_SCRIPT}") ]]; then
  mkdir -p "${WRAPPER_DEST_DIR}"
  printf "%b" "${WRAPPER_SCRIPT}" > "${WRAPPER_DEST_FILE}"
  chmod +x "${WRAPPER_DEST_FILE}"
fi
