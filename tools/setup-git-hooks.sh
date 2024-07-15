#!/bin/bash

# This script is used to setup git hooks for the repository.
SCRIPT_ROOT="$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")"
REPO_ROOT="${SCRIPT_ROOT}"

# Use alternate root if provided as command line argument.
if [[ -n "${1-}" ]] && [[ "$1" == *"/"* ]]; then
  REPO_ROOT="${1}"
  shift
fi

GITHUB_ROOT="${SCRIPT_ROOT}/.github"
GIT_REPO_ROOT="${REPO_ROOT}/.git"

if [[ -f "${GITHUB_ROOT}/hooks/pre-commit" && -f "${GIT_REPO_ROOT}/hooks/pre-commit" ]]; then
  echo "Setting up git pre-commit hooks for ${REPO_ROOT}"
  ln -s -f "${GITHUB_ROOT}/hooks/pre-commit" "${GIT_REPO_ROOT}/hooks/pre-commit"
fi
