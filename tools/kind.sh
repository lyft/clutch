#!/bin/bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"
BUILD_BIN="${BUILD_ROOT}/bin"
KUBECONFIG=$BUILD_ROOT/kubeconfig-clutch


NAME=kind
RELEASE=v0.10.0
OSX_RELEASE_SUM=a934e573621917a2785f3ddfa7b6187d18fa1c20c94c013919736b3256d37f57
LINUX_RELEASE_SUM=74767776488508d847b0bb941212c1cb76ace90d9439f4dee256d8a04f1309c6

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
    KUBECONFIG=$KUBECONFIG kubectl create ns "stateful-${env}" || true

    # Creating resources in `envoy-*` namespace
    KUBECONFIG=$KUBECONFIG kubectl create deployment envoy --image envoyproxy/envoy:v1.14-latest -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl autoscale deployment envoy --cpu-percent=50 --min=1 --max=2 -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl expose deployment envoy --port=8080 -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl create configmap "configmap-${env}-test-1" --from-literal=environment="${env}" -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl create configmap "configmap-${env}-test-2" --from-literal=environment="${env}" -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl create job "job-${env}-test-1" --image busybox -n "envoy-${env}" || true
    KUBECONFIG=$KUBECONFIG kubectl create job "job-${env}-test-2" --image busybox -n "envoy-${env}" || true


    # Creating resources in `cron-*` namespace
    KUBECONFIG=$KUBECONFIG kubectl create cronjob cron-test --schedule "*/1 * * * *" --image busybox -n "cron-${env}" || true

    # Creating resources in `stateful-*` namespace
    KUBECONFIG=$KUBECONFIG kubectl apply -f tools/kind-stateful-set.yaml  || true

    # Adding labels to resources
    KUBECONFIG=$KUBECONFIG kubectl label configmap "configmap-${env}-test-1" -n "envoy-${env}" app=envoy || true
    KUBECONFIG=$KUBECONFIG kubectl label job "job-${env}-test-1" -n "envoy-${env}" app=envoy || true
    KUBECONFIG=$KUBECONFIG kubectl label cronjob cron-test -n "cron-${env}" app=cron || true
    KUBECONFIG=$KUBECONFIG kubectl annotate job "job-${env}-test-1" -n "envoy-${env}" url=foo@example.com || true
    KUBECONFIG=$KUBECONFIG kubectl annotate cronjob cron-test -n "cron-${env}" url=foo@example.com || true
  done
}

main "$@"
