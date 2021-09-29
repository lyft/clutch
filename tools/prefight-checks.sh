#!/usr/bin/env bash
set -euo pipefail

# Will be set to false if any of the steps fail
did_checks_pass=true

# Minimum versions
MIN_GO_VERSION="1.17.0"
MIN_NODE_VERSION="14.0.0"
MIN_YARN_VERSION="1.22.11"

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
  # If were on OSX lets check for brew and coreutils as they are requriments
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
        echo "brew is not installed, unable to verify coreutils dependency."
        did_checks_pass=false
    fi
  fi
}

backend() {
  echo "Checking backend"
  if ! command -v go -v &> /dev/null; then
    echo "golang is not installed, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$(go version | { read -r _ _ v _; echo "${v#go}"; })
    if ! is_version_ok $MIN_GO_VERSION "$current_version"; then
      echo "golang version must be >= $MIN_GO_VERSION"
      did_checks_pass=false
    fi
  fi
}

frontend() {
  echo "Checking frontend"
  if ! command -v node -v &> /dev/null; then
    echo "nodejs is not installed, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$(node --version)
    # remove the leading v from the version output
    nov=${current_version:1}
    if ! is_version_ok $MIN_NODE_VERSION "$nov"; then
      echo "node version must be >= $MIN_NODE_VERSION"
      did_checks_pass=false
    fi
  fi

  if ! command -v yarn &> /dev/null; then
    echo "yarn is not installed, this is a required dependency."
    did_checks_pass=false
  else
    current_version=$(yarn --version)
    if ! is_version_ok $MIN_YARN_VERSION "$current_version:1"; then
      echo "yarn version must be >= $MIN_YARN_VERSION"
      did_checks_pass=false
    fi
  fi
}

main() {
  # always check OS level requirments
  os

  if [ $# -ge 1 ] && [ -n "$1" ]; then
    if [ "$1" == "backend" ]; then
      backend
    elif [ "$1" == "frontend" ]; then
      frontend
    else
      backend
      frontend
    fi
  fi

  if [ "$did_checks_pass" = false ] ; then
    printf "\nPlease refer to the development requirments https://clutch.sh/docs/getting-started/local-build/#requirements"
    return 1
  fi
}

main "$@"
