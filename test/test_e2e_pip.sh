#!/usr/bin/env bash
GREEN='\033[0;32m'
RED='\033[0;31m'
SET='\033[0m'
set -e

echo -e "Running pip resolution e2e tests..."

go run ./cmd/debricked/main.go files resolve test/testdata/pip/requirements.txt

if diff test/testdata/pip/expected.lock test/testdata/pip/requirements.txt.debricked.lock >/dev/null ; then
    echo "Pip resolution e2e tests completed ${GREEN}successfully${SET}"
else
    echo "Pip resolution e2e tests completed ${RED}unsuccessfully${SET}"
fi


