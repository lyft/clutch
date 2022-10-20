#!/bin/bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"

NAME=golangci-lint
RELEASE=v1.47.3
RELEASE_BINARY="${BUILD_BIN}/${NAME}"

main() {
  cd "${REPO_ROOT}/backend"
  ensure_binary

  "${RELEASE_BINARY}" "$@"
}

ensure_binary() {
  if [[ ! -f "${RELEASE_BINARY}" ]]; then
    echo "info: Downloading ${NAME} ${RELEASE} to build environment"

    mkdir -p "${BUILD_BIN}"
    env GOBIN="${BUILD_BIN}" go install github.com/golangci/golangci-lint/cmd/golangci-lint@"${RELEASE}"
  fi
}

main "$@"
