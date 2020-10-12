#!/bin/bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-"$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"}"
BUILD_ROOT="${REPO_ROOT}/build"

printf "\n\nCreating fake resources in local k8s cluster\n\n"


for env in staging production; do
    # Creating namespaces
    KUBECONFIG=$BUILD_ROOT/kubeconfig-clutch kubectl create ns "envoy-${env}" || true

    # Creating resources in `envoy-*` namespace
    KUBECONFIG=$BUILD_ROOT/kubeconfig-clutch kubectl create deployment envoy --image envoyproxy/envoy:v1.14-latest -n "envoy-${env}" || true
    KUBECONFIG=$BUILD_ROOT/kubeconfig-clutch kubectl autoscale deployment envoy --cpu-percent=50 --min=1 --max=2 -n "envoy-${env}" || true
done
