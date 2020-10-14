#!/usr/bin/env bash
set -euo pipefail

# https://github.com/protocolbuffers/protobuf/releases
PROTOC_RELEASE=3.13.0
PROTO_ZIP_RELEASE_MD5_LINUX=20a5326fbc666e1fd069eaa80875fcac
PROTO_ZIP_RELEASE_MD5_OSX=a545b7fb818dba564aa161e4232e69f5

# https://github.com/protobufjs/protobuf.js/releases
# NOTE: should match frontend/package.json
PROTOBUFJS_RELEASE=6.10.1

# https://github.com/angular/clang-format/releases
ANGULAR_CLANG_FORMAT_RELEASE=1.4.0
ANGULAR_CLANG_FORMAT_RELEASE_MD5_LINUX=fee8c52e196e28ae5928d6ff8757f58c
ANGULAR_CLANG_FORMAT_RELEASE_MD5_OSX=c3ebe742599dcc38b9dc6544cacd69bb

PROTOS=()
PROTO_DIRS=()
CLUTCH_PROTOS=()

SCRIPT_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"

# parse options
ACTION="compile"
LINT_FIX=false
CLUTCH_API_ROOT=""
while getopts "lfc:" opt; do
  case $opt in
    l) ACTION="lint" ;;
    f) LINT_FIX=true ;;
    c) CLUTCH_API_ROOT="${OPTARG}" ;;
    *) echo "usage: $0 [-l]" >&2
     exit 1 ;;
  esac
done
shift "$((OPTIND-1))" # shift so that $@, $1, etc. refer to the non-option arguments

main() {
  check_prereqs

  REPO_ROOT="${SCRIPT_ROOT}"
  # Use alternate root if provided as command line argument.
  if [[ -n "${1-}" ]]; then
    REPO_ROOT="${1}"
  fi

  if [[ -z "${CLUTCH_API_ROOT}" ]]; then
    # if core is not provided then we need to use a downloaded version.
    CORE_VERSION=$(cd "${REPO_ROOT}/backend" && go list -f "{{ .Version }}" -m github.com/lyft/clutch/backend)
    if [[ "${CORE_VERSION}" == *-*-* ]]; then
      # if a pseudo-version, figure out just the SHA
      CORE_VERSION=$(echo "${CORE_VERSION}" | awk -F"-" '{print $NF}')
    fi

    core_out="${REPO_ROOT}/build/bin/clutch-api-${CORE_VERSION}"
    if [[ ! -d "${core_out}" ]]; then
      echo "info: downloading core APIs ${CORE_VERSION} to build environment..."

      core_zip_out="/tmp/clutch-${CORE_VERSION}.tar.gz"
      core_tmp_out="/tmp/clutch-${CORE_VERSION}"
      curl -sSL -o "${core_zip_out}" \
        "https://github.com/lyft/clutch/archive/${CORE_VERSION}.zip"

      mkdir -p "${core_tmp_out}"
      unzip -q -o "${core_zip_out}" -d "${core_tmp_out}"

      mkdir -p "${core_out}"
      mv "${core_tmp_out}"/clutch-*/api "${core_out}"
    fi

    CLUTCH_API_ROOT="${core_out}/api"
  fi

  API_ROOT="${REPO_ROOT}/api"
  BUILD_ROOT="${REPO_ROOT}/build"

  cd "${REPO_ROOT}/backend"

  prepare_build_environment
  discover_protos

  grpc_gateway_include_path="$(modpath github.com/grpc-ecosystem/grpc-gateway)/third_party/googleapis"
  pg_validate_include_path="$(modpath github.com/envoyproxy/protoc-gen-validate)"

  # Lint (fix) and exit if requested.
  if [[ "${ACTION}" == "lint" ]]; then
    cd "${API_ROOT}"

    buf_lint_config=$(cat "${API_ROOT}/buf.json")

    LINT_OK=true
    if [[ ${LINT_FIX} == true ]]; then
      for proto in "${PROTOS[@]}"; do
        "${CLANG_FORMAT_BIN}" --style=file -i "${proto}"
      done
    else
      for proto in "${PROTOS[@]}"; do
        if ! output=$("${CLANG_FORMAT_BIN}" --style=file "${proto}" | diff -u "${proto}" -); then
          echo "${output}"
          LINT_OK=false
        fi
      done

      for proto in "${PROTOS[@]}"; do
        if ! output=$("${PROTOC_BIN}" \
          -I"${PROTOC_INCLUDE_DIR}" -I"${API_ROOT}" \
          -I"${grpc_gateway_include_path}" -I"${pg_validate_include_path}" \
          --buf-check-lint_out=. \
          "--buf-check-lint_opt={\"input_config\": ${buf_lint_config}}" \
          --plugin=protoc-gen-buf-check-lint="${GOBIN}/protoc-gen-buf-check-lint" \
          "${proto}" 2>&1)
        then
          echo "--- ${proto}"
          echo "${output}" | sed 's/--buf-check-lint_out: //' | cut -d":" -f2-
          LINT_OK=false
        fi
      done
    fi
    ${LINT_OK} && exit 0 || exit 1
  fi

  # Compile.
  proto_out_dir="${REPO_ROOT}/backend/api"
  mkdir -p "${proto_out_dir}"

  MFLAGS=""
  package_dir="$(grep -m1 module go.mod | cut -d' ' -f2)/api"

  for proto in "${PROTOS[@]}"; do
    relative_path="${proto/#${API_ROOT}\/}"
    MFLAGS+="M${relative_path}=$(dirname "${package_dir}/${relative_path}"),"
  done

  # Add MFLAGS for Clutch protos when this is running for a private gateway.
  if [[ "${CLUTCH_API_ROOT}" != "${API_ROOT}" ]]; then
    discover_core_protos
    readonly CLUTCH_PREFIX="github.com/lyft/clutch/backend/api"
    for proto in "${CLUTCH_PROTOS[@]}"; do
      relative_path="${proto/#${CLUTCH_API_ROOT}\/}"
      MFLAGS+="M${relative_path}=${CLUTCH_PREFIX}/$(dirname "${relative_path}"),"
    done
  fi

  echo "info: compiling go"
  for proto_dir in "${PROTO_DIRS[@]}"; do
    echo "${proto_dir}"
    "${PROTOC_BIN}" \
      -I"${PROTOC_INCLUDE_DIR}" -I"${API_ROOT}" -I"${CLUTCH_API_ROOT}" \
      -I"${grpc_gateway_include_path}" -I"${pg_validate_include_path}" \
      --go_out="${MFLAGS}"plugins=grpc:"${proto_out_dir}" \
      --validate_out="${MFLAGS}"lang=go:"${proto_out_dir}" \
      --grpc-gateway_out="${proto_out_dir}" \
      --grpc-gateway_opt=warn_on_unbound_methods=true \
      --plugin=protoc-gen-go="${GOBIN}/protoc-gen-go" \
      --plugin=protoc-gen-grpc-gateway="${GOBIN}/protoc-gen-grpc-gateway" \
      --plugin=protoc-gen-validate="${GOBIN}/protoc-gen-validate" \
      "${proto_dir}"/*.proto
  done

  echo "info: compiling javascript bundle"
  cd ..
  mkdir -p "${REPO_ROOT}/frontend/api/src"
  js_out="frontend/api/src/index.js"
  "${PROTOBUFJS_DIR}/node_modules/.bin/pbjs" \
    -p "${PROTOC_INCLUDE_DIR}" -p "${API_ROOT}" -p"${CLUTCH_API_ROOT}" \
    -p "${grpc_gateway_include_path}" -p "${pg_validate_include_path}" \
    -t static-module \
    --no-create --no-encode --no-decode --no-delimited \
    -w es6 --es6 \
    -o "${js_out}" \
    "${PROTOS[@]}"
  echo -e "// Code generated by protobuf.js in compile-protos.sh. DO NOT EDIT.\n\n$(cat "${js_out}")" > "${js_out}"

  ts_out="frontend/api/src/index.d.ts"
  "${PROTOBUFJS_DIR}/node_modules/.bin/pbts" \
    -o "${ts_out}" \
    "frontend/api/src/index.js"
  echo -e "// Code generated by protobuf.js in compile-protos.sh. DO NOT EDIT.\n\n$(cat "${ts_out}")" > "${ts_out}"

  echo "OK"
}

discover_protos() {
  while IFS= read -r -d '' proto; do
    PROTOS+=("${proto}")
  done <  <(find "${API_ROOT}" -name '*.proto' -print0 | sort -sdzu)

  while IFS= read -r -d '' proto_dirs; do
    PROTO_DIRS+=("${proto_dirs}")
  done <  <(find "${API_ROOT}" -name '*.proto' -exec dirname {} \; | tr '\n' '\0' | sort -sdzu)
}

discover_core_protos() {
  while IFS= read -r -d '' proto; do
    CLUTCH_PROTOS+=("${proto}")
  done <  <(find "${CLUTCH_API_ROOT}" -name '*.proto' -print0 | sort -sdzu)
}


# Get the directory that the go module is stored in and ensure that it's the correct version.
modpath() {
  set -e
  go mod download "${1}"
  go list -f "{{ .Dir }}" -m "${1}"
}

prepare_build_environment() {
  export GOBIN="${BUILD_ROOT}/bin"
  mkdir -p "${GOBIN}"

  install_protoc

  if [[ "${ACTION}" == "compile" ]]; then
    install_protobufjs
  fi

  if [[ "${ACTION}" == "lint" ]]; then
    install_clang_format
  fi
}

check_prereqs() {
  if ! command -v "npm" &> /dev/null; then
    echo "ERROR: npm not found, see https://github.com/lyft/clutch/wiki/Requirements#nodejs for more information."
    exit 1
  fi
}

install_protobufjs() {
  export PROTOBUFJS_DIR="${BUILD_ROOT}/bin/protobufjs-${PROTOBUFJS_RELEASE}"
  if [[ ! -f "${PROTOBUFJS_DIR}/node_modules/.bin/pbjs" ]]; then
    echo "info: Downloading protobufjs to build environment"
    mkdir -p "${PROTOBUFJS_DIR}"
    "${BUILD_ROOT}/bin/yarn.sh" --cwd "${PROTOBUFJS_DIR}" add --frozen-lockfile "protobufjs@${PROTOBUFJS_RELEASE}"
  fi
}

install_clang_format() {
  CLANG_FORMAT_BIN="${GOBIN}/clang-format-${ANGULAR_CLANG_FORMAT_RELEASE}"
  if [[ ! -f "${CLANG_FORMAT_BIN}" ]]; then
    echo "info: Downloading clang-format to build environment"

    case "${OSTYPE}" in
      "darwin"*) clang_format_os="darwin"; clang_format_md5=${ANGULAR_CLANG_FORMAT_RELEASE_MD5_OSX} ;;
      "linux"*) clang_format_os="linux"; clang_format_md5=${ANGULAR_CLANG_FORMAT_RELEASE_MD5_LINUX} ;;
      *) echo "error: Unsupported OS '${OSTYPE}' for clang-format install, please install manually" && exit 1 ;;
    esac

    curl -sSL -o "/tmp/clang-format" \
      "https://github.com/angular/clang-format/raw/v${ANGULAR_CLANG_FORMAT_RELEASE}/bin/${clang_format_os}_x64/clang-format"
    echo ${clang_format_md5} "/tmp/clang-format" | md5sum --check --quiet -
    chmod a+x "/tmp/clang-format"
    mv "/tmp/clang-format" "${CLANG_FORMAT_BIN}"
  fi
}

install_protoc() {
  export PROTOC_BIN="${GOBIN}/protoc-v${PROTOC_RELEASE}"
  export PROTOC_INCLUDE_DIR="${GOBIN}/protoc-v${PROTOC_RELEASE}-include"

  go install \
    github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    github.com/golang/protobuf/protoc-gen-go \
    github.com/envoyproxy/protoc-gen-validate

  if [[ ! -f "${PROTOC_BIN}" || ! -d "${PROTOC_INCLUDE_DIR}" ]]; then
    echo "info: Downloading protoc-v${PROTOC_RELEASE} to build environment"

    proto_arch=x86_64
    case "${OSTYPE}" in
      "darwin"*) proto_os="osx"; proto_md5="${PROTO_ZIP_RELEASE_MD5_OSX}" ;;
      "linux"*) proto_os="linux"; proto_md5="${PROTO_ZIP_RELEASE_MD5_LINUX}" ;;
      *) echo "error: Unsupported OS '${OSTYPE}' for protoc install, please install manually" && exit 1 ;;
    esac

    proto_zip_out="/tmp/protoc-${PROTOC_RELEASE}.zip"
    curl -sSL -o "${proto_zip_out}" \
      "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_RELEASE}/protoc-${PROTOC_RELEASE}-${proto_os}-${proto_arch}.zip"
    echo ${proto_md5} ${proto_zip_out} | md5sum --check --quiet -

    proto_dir_out="/tmp/proto-${PROTOC_RELEASE}"
    mkdir -p "${proto_dir_out}"
    unzip -q -o "${proto_zip_out}" -d "${proto_dir_out}"

    if [[ ! -f ${PROTOC_BIN} ]]; then
      find "${GOBIN}" -maxdepth 0 -regex '.*/protoc-v[0-9\.]+$' -exec rm {} \;  # cleanup older versions
      mv "${proto_dir_out}"/bin/protoc "${PROTOC_BIN}"
    fi

    if [[ ! -d "${PROTOC_INCLUDE_DIR}" ]]; then
      find "${GOBIN}" -maxdepth 0 -regex '.*/protoc-v.*?-include$' -type d -exec rm -r {} \;  # cleanup older versions
      mv "${proto_dir_out}"/include "${PROTOC_INCLUDE_DIR}"
    fi
  fi
}

main "$@"
