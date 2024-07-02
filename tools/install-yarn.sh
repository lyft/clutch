#!/bin/bash

YARN_VERSION="4.3.1"
ROOT_DIR="$PWD"
ROOT_DEST_DIR="${ROOT_DIR}/frontend"

ARGS="$*"

if [[ $ARGS = *'CWD'* ]]; then
  echo "CWD arg detected. Utilizing working directory for install."

  ROOT_DEST_DIR="${ROOT_DIR}"
elif [[ ! -d "${ROOT_DEST_DIR}" ]]; then
  echo "Could not find frontend directory. Ensure you're running this script from the root of your project."
  exit 1
fi

if [[ $ARGS = *'bin='* ]]; then
  echo "bin arg detected. Utilizing specified directory for wrapper script."
  BIN_DIR=$(echo "${ARGS}" | sed 's/.*bin\=//' | cut -d ' ' -f1)
fi

DEST_DIR="${ROOT_DEST_DIR}/.yarn/releases"
DEST_FILE="${DEST_DIR}/yarn-${YARN_VERSION}.js"
YARN_VERSION_FILE=".yarn/releases/yarn-${YARN_VERSION}.js"
WRAPPER_DEST_DIR="${ROOT_DIR}/${BIN_DIR:-"build/bin"}"
WRAPPER_DEST_FILE="${WRAPPER_DEST_DIR}/yarn.sh"

if [[ ! -f "${DEST_FILE}" ]]; then
  echo "Downloading yarn v${YARN_VERSION} to build environment..."
  mkdir -p "${DEST_DIR}"
  curl -sSL -o "${DEST_FILE}" \
    "https://repo.yarnpkg.com/${YARN_VERSION}/packages/yarnpkg-cli/bin/yarn.js"
fi

# Install a wrapper script in build/ that executes yarn if it doesn't exist already.
WRAPPER_SCRIPT="#!/bin/bash\nnode \"${DEST_FILE}\" \"\$@\"\n"
if [[ ! -f "${WRAPPER_DEST_FILE}" || $(< "${WRAPPER_DEST_FILE}") != $(printf "%b" "${WRAPPER_SCRIPT}") ]]; then
  mkdir -p "${WRAPPER_DEST_DIR}"
  printf "%b" "${WRAPPER_SCRIPT}" > "${WRAPPER_DEST_FILE}"
  chmod +x "${WRAPPER_DEST_FILE}"
fi

#Link script to yarn config
cd "${ROOT_DEST_DIR}" || exit
"${WRAPPER_DEST_FILE}" config set yarnPath "${YARN_VERSION_FILE}"
