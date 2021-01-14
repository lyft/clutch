#!/bin/bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"
KUBECONFIG=$BUILD_ROOT/kubeconfig-clutch


NAME=kind
RELEASE=v0.9.0
OSX_RELEASE_SUM=849034ffaea8a0e50f9153078890318d5863bafe01495418ea0ad037b518de90
LINUX_RELEASE_SUM=35a640e0ca479192d86a51b6fd31c657403d2cf7338368d62223938771500dc8

ARCH=amd64

RELEASE_BINARY="${BUILD_BIN}/${NAME}-${RELEASE}"

main() {
  cd "${REPO_ROOT}"
  ensure_binary

  if [ "$1" == "seed" ]; then
    seed
    exit 0
  fi

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

seed() {
  printf "\n\nCreating fake resources in clutch-local k8s cluster\n\n"
  for env in staging production; do
    # Creating namespaces
    KUBECONFIG=$KUBECONFIG kubectl create ns "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl create ns "cron-${env}" || true

    # Creating resources in `envoy-*` namespace
    KUBECONFIG=$KUBECONFIG kubectl create deployment envoy --image envoyproxy/envoy:v1.14-latest -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl autoscale deployment envoy --cpu-percent=50 --min=1 --max=2 -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl expose deployment envoy --port=8080 -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl create cronjob cron-test --schedule "*/1 * * * *" --image busybox -n "cron-${env}" || true
  done
}

main "$@"
