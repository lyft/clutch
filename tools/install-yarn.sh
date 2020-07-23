#!/bin/bash

YARN_VERSION="1.22.4"
ROOT_DEST_DIR="$PWD/frontend"
DEST_DIR="$ROOT_DEST_DIR/.yarn/releases"
DEST_FILE="${DEST_DIR}/yarn-${YARN_VERSION}.js"
WRAPPER_DEST="${PWD}/build/bin/yarn.sh"

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
if [[ ! -f "${WRAPPER_DEST}" || $(< "${WRAPPER_DEST}") != $(printf "%b" "${WRAPPER_SCRIPT}") ]]; then
  printf "%b" "${WRAPPER_SCRIPT}" > "${WRAPPER_DEST}"
  chmod +x "${WRAPPER_DEST}"
fi
