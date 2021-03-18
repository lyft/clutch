#!/bin/bash
trap "exit" INT TERM
trap 'kill $(jobs -p)' EXIT
set -euo pipefail

BE_STARTUP_COUNT=0
FE_STARTUP_COUNT=0
STARTUP_WAIT=30

make backend-dev-mock &
yarn --cwd frontend workspace @clutch-sh/app start &
until curl --output /dev/null --silent --head --fail http://localhost:3000; do
    if [ "$FE_STARTUP_COUNT" -ge "$STARTUP_WAIT" ]; then
        echo "Error: could not start frontend dev server"
        exit 1
    fi;
    ((FE_STARTUP_COUNT++))
    sleep 1
done

until curl --output /dev/null --silent --fail http://localhost:8080/healthcheck; do
    if [ "$BE_STARTUP_COUNT" -ge "$STARTUP_WAIT" ]; then
        echo "Error: could not start backend mock server"
        exit 1
    fi;
    ((BE_STARTUP_COUNT++))
    sleep 1
done
(
  yarn --cwd frontend test:e2e
)
