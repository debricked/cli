#!/usr/bin/env bash
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
SET='\033[0m'
echo "Lint"
if ! command -v golangci-lint &> /dev/null
then
    echo -e "${YELLOW}golangci-lint${SET} could not be found. Make sure it is installed"
    echo -e "${RED}FAILED${SET}"
    exit
fi
if ! golangci-lint run ./...;
then
  echo -e "${RED}FAILED${SET}"
  exit
else
  echo -e "${GREEN}OK${SET}"
fi