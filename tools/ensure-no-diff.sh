#!/usr/bin/env bash
set -euo pipefail

modified=$(git status --porcelain "$@")
if [[ -n "${modified}" ]]; then
  git --no-pager diff HEAD "$@"
  untracked=$(echo "${modified}" | grep '??')
  if [[ -n "${untracked}" ]]; then
    echo -e "\n\nUNTRACKED FILES:"
    echo "${untracked}"  | awk '{print "+++ " $2}'
  fi
  echo -e "\nerror: commit changes to the generated files above"
  exit 1
fi
