#!/usr/bin/env bash
set -euo pipefail

# Will be set to false if any of the steps fail
did_checks_pass=true

# Minimum versions
MIN_GO_VERSION="1.17"
MIN_NODE_VERSION="18.0.0"
MIN_YARN_VERSION="4.1.1"

SCRIPT_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"
REPO_ROOT="${SCRIPT_ROOT}"
# Use alternate root if provided as command line argument.
if [[ -n "${1-}" ]] && [[ "$1" == *"/"* ]]; then
  REPO_ROOT="${1}"
  shift
fi

echo "${REPO_ROOT}"

BUILD_ROOT="${REPO_ROOT}/build"

# param 1 - required version
# param 2 - current version
# returns true or false if the version is ok
is_version_ok() {
  required_version=$1
  current_version=$2
  if [ "$(printf '%s\n' "$required_version" "$current_version" | sort -V | head -n1)" = "$required_version" ]; then 
    return 0
  else
    return 1
  fi
}

os() {
  # If we're on OSX lets check for brew and coreutils as they are requirements
  # https://clutch.sh/docs/getting-started/local-build/#requirements
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # check brew is installed
    if command -v brew --version &> /dev/null; then
      # check if coreutils is installed
      if ! brew ls --versions coreutils > /dev/null; then
        echo "coreutils is not installed, this is a required dependency."
        echo "install by running [brew install coreutils]"
        did_checks_pass=false
      fi
    else
        echo "brew is not installed (visit https://brew.sh/ to install), unable to verify coreutils dependency."
        did_checks_pass=false
    fi
  fi
}

backend() {
  if ! command -v go -v &> /dev/null; then
    echo "golang is not installed or cannot be found in the current PATH, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$(go version | { read -r _ _ v _; echo "${v#go}"; })
    if ! is_version_ok $MIN_GO_VERSION "$current_version"; then
      echo "golang version must be >= $MIN_GO_VERSION, current version $current_version"
      did_checks_pass=false
    fi
  fi
}

frontend() {
  if ! command -v node -v &> /dev/null; then
    echo "nodejs is not installed or cannot be found in the current PATH, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$(node --version)
    # remove the leading v from the version output
    nov=${current_version:1}
    if ! is_version_ok $MIN_NODE_VERSION "$nov"; then
      echo "node version must be >= $MIN_NODE_VERSION, current version $nov"
      did_checks_pass=false
    fi
  fi

  if ! command -v "${BUILD_ROOT}/bin/yarn.sh" &> /dev/null; then
    echo "yarn is not installed or cannot be found in the current PATH, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$("${BUILD_ROOT}/bin/yarn.sh" --version)
    if ! is_version_ok $MIN_YARN_VERSION "$current_version"; then
      echo "yarn version must be >= $MIN_YARN_VERSION, current version $current_version"
      did_checks_pass=false
    fi
  fi
}

main() {
  echo "Running pre-flight checks..."
  # always check OS level requirements
  os

  if [ $# -ge 1 ] && [ -n "$1" ]; then
    if [ "$1" == "backend" ]; then
      backend
    elif [ "$1" == "frontend" ]; then
      frontend
    fi
  else
    backend
    frontend
  fi

  if [ "$did_checks_pass" = false ] ; then
    printf "\nPlease refer to the development requirements https://clutch.sh/docs/getting-started/local-build/#requirements"
    return 1
  fi
  echo "Pre-flight checks satisfied!"
}

main "$@"