#!/bin/bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"

NAME=golangci-lint
RELEASE=v1.53.2
OSX_RELEASE_256=a4e83f5bfe52f42134c9783aa68ba31104c36e2ad4c221a3c77510dda66ae81c
LINUX_RELEASE_256=2298f73b9bc03b88b91fee06c5d519fc7f9d7f328e2c388615bbd7e85a9d6cae

ARCH=amd64

RELEASE_BINARY="${BUILD_BIN}/${NAME}-${RELEASE}"

main() {
  cd "${REPO_ROOT}/backend"
  ensure_binary

  "${RELEASE_BINARY}" "$@"
}

ensure_binary() {
  if [[ ! -f "${RELEASE_BINARY}" ]]; then
    echo "info: Downloading ${NAME} ${RELEASE} to build environment"

    mkdir -p "${BUILD_BIN}"

    case "${OSTYPE}" in
      "darwin"*) os_type="darwin"; sum="${OSX_RELEASE_256}" ;;
      "linux"*) os_type="linux"; sum="${LINUX_RELEASE_256}" ;;
      *) echo "error: Unsupported OS '${OSTYPE}' for shellcheck install, please install manually" && exit 1 ;;
    esac

    release_archive="/tmp/${NAME}-${RELEASE}.tar.gz"

    URL="https://github.com/golangci/golangci-lint/releases/download/${RELEASE}/golangci-lint-${RELEASE:1}-${os_type}-${ARCH}.tar.gz"
    curl -sSL -o "${release_archive}" "${URL}"
    echo ${sum} ${release_archive} | sha256sum --check --quiet -

    release_tmp_dir="/tmp/${NAME}-${RELEASE}"
    mkdir -p "${release_tmp_dir}"
    tar -xzf "${release_archive}" --strip=1 -C "${release_tmp_dir}"

    if [[ ! -f "${RELEASE_BINARY}" ]]; then
      find "${BUILD_BIN}" -maxdepth 0 -regex '.*/'${NAME}'-[A-Za-z0-9\.]+$' -exec rm {} \;  # cleanup older versions
      mv "${release_tmp_dir}/${NAME}" "${RELEASE_BINARY}"
    fi

    # Cleanup stale resources.
    rm "${release_archive}"
    rm -rf "${release_tmp_dir}"
  fi
}

main "$@"
