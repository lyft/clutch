#!/bin/bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"


NAME=kind
RELEASE=v0.8.1
OSX_RELEASE_SUM=cdd8dfe7dff764429badcd636179b0e3eb937640cfe56749dd9b8f9c048cb7db
LINUX_RELEASE_SUM=781c3db479b805d161b7c2c7a31896d1a504b583ebfcce8fcd49538c684d96bc

ARCH=amd64

RELEASE_BINARY="${BUILD_BIN}/${NAME}-${RELEASE}"

main() {
  cd "${REPO_ROOT}"
  ensure_binary

  "${RELEASE_BINARY}" "$@"
}

ensure_binary() {
  if [[ ! -f "${RELEASE_BINARY}" ]]; then
    echo "info: Downloading ${NAME} ${RELEASE} to build environment"

    mkdir -p "${BUILD_BIN}"

    case "${OSTYPE}" in
      "darwin"*) os_type="darwin"; sum="${OSX_RELEASE_SUM}" ;;
      "linux"*) os_type="linux"; sum="${LINUX_RELEASE_SUM}" ;;
      *) echo "error: Unsupported OS '${OSTYPE}' for kind install, please install manually" && exit 1 ;;
    esac

    release_archive="/tmp/${NAME}-${RELEASE}"

    URL="https://github.com/kubernetes-sigs/kind/releases/download/${RELEASE}/kind-${os_type}-${ARCH}"
    curl -sSL -o "${release_archive}" "${URL}"
    echo ${sum} ${release_archive} | sha256sum --check --quiet -

    find "${BUILD_BIN}" -maxdepth 0 -regex '.*/'${NAME}'-[A-Za-z0-9\.]+$' -exec rm {} \;  # cleanup older versions
    mv "${release_archive}" "${RELEASE_BINARY}"
    chmod +x "${RELEASE_BINARY}"
  fi
}

main "$@"
