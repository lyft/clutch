#!/bin/bash

YARN_VERSION="4.1.1"

if ! command -v corepack *> /dev/null; then
  echo "Corepack must be installed, please upgrade your node version to >18"
  exit 1
fi

echo "Installing yarn@${YARN_VERSION} with corepack"
corepack enable
corepack prepare yarn@${YARN_VERSION} --activate
