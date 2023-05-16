#!/usr/bin/env bash
GREEN='\033[0;32m'
RED='\033[0;31m'
SET='\033[0m'
set -e

go test -cover -coverprofile=coverage.out ./internal/...

echo -e "\nChecking test coverage threshold..."
regex='[0-9]+\.*[0-9]*'
if ! [[ $TEST_COVERAGE_THRESHOLD =~ $regex ]]; then
  echo "Failed to find test coverage threshold. Defaults to 95%"
  TEST_COVERAGE_THRESHOLD=95
fi
echo "Test coverage threshold     : $TEST_COVERAGE_THRESHOLD %"
if [ ! -f "./coverage.out" ]; then
  echo "Failed to find coverage.out. Make sure coverage.out is created"
  echo -e "${RED}FAILED${SET}"
  exit 1
fi

# Find test coverage
totalTestCoverage=$(go tool cover -func=coverage.out | grep total | grep -Eo "$regex")
# Store coverage report
go tool cover -html=coverage.out -o=coverage.html

echo "Current test coverage       : $totalTestCoverage %"
if (( $(echo "$totalTestCoverage $TEST_COVERAGE_THRESHOLD" | awk '{print ($1 > $2)}') )); then
  echo -e "${GREEN}OK${SET}"
else
  echo "Current test coverage in below threshold. Please extend your unit tests"
  echo -e "${RED}FAILED${SET}"
  exit 1
fi
