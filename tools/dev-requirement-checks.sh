#!/usr/bin/env bash
set -euo pipefail

# Will be set to false if any of the steps fail
# 
did_checks_pass=true

# If were on OSX lets check for brew and coreutils as they are requriments
# https://clutch.sh/docs/getting-started/local-build/#requirements
if [[ "$OSTYPE" == "darwin"* ]]; then
  # check brew is installed
  if command -v brew &> /dev/null; then
    # check if coreutils is installed
    if ! brew ls --versions coreutils > /dev/null; then
      echo "coreutils is not installed, this is a required dependency."
      echo "install by running [brew install coreutils]"
      did_checks_pass=false
    fi
  else
      echo "brew is not installed, unable to verify coreutils dependency."
      did_checks_pass=false
  fi
fi

if ! command -v go -v &> /dev/null; then
  echo "golang is not installed, this is a required dependency."
  did_checks_pass=false
fi

if ! command -v node -v &> /dev/null; then
  echo "nodejs is not installed, this is a required dependency."
  did_checks_pass=false
fi

if ! command -v yarn &> /dev/null; then
  echo "yarn is not installed, this is a required dependency."
  did_checks_pass=false
fi

if [ "$did_checks_pass" = false ] ; then
  printf "\nPlease refer to the development requirments https://clutch.sh/docs/getting-started/local-build/#requirements"
fi
