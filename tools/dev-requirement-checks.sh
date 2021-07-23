#!/usr/bin/env bash
set -euo pipefail

# If were on OSX lets check for brew and coreutils as they are requriments
# https://clutch.sh/docs/getting-started/local-build/#requirements
if [[ "$OSTYPE" == "darwin"* ]]; then
  # check brew is installed
  if command -v brew &> /dev/null; then
    # check if coreutils is installed
    if ! brew ls --versions coreutils > /dev/null; then
      echo "coreutils is not installed, this is a required dependency."
      echo "brew install coreutils"
      exit
    fi
  else
      echo "brew is not installed, unable to verify coreutils dependency."
      echo "Please refer to the development requirments https://clutch.sh/docs/getting-started/local-build/#requirements"
      exit
  fi
fi
