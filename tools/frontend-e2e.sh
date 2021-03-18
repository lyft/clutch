#!/bin/bash
# Test comment
trap "exit" INT TERM
trap 'kill $(jobs -p)' EXIT
set -euo pipefail

make backend-dev-mock &
yarn --cwd frontend workspace @clutch-sh/app start &
until curl --output /dev/null --silent --head --fail http://localhost:3000; do
    sleep 1
done
(
  yarn --cwd frontend test:e2e
)