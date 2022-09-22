#!/usr/bin/env bash
GREEN='\033[0;32m'
RED='\033[0;31m'
SET='\033[0m'

echo "Lint shell scripts"
if ! shellcheck ./**/*.sh;
then
  echo -e "${RED}FAILED${SET}"
  exit 1
else
  echo -e "${GREEN}OK${SET}"
fi

echo "Vet"
if ! go vet ./...;
then
  echo -e "${RED}FAILED${SET}"
  exit 1
else
  echo -e "${GREEN}OK${SET}"
fi

echo "Format (fmt)"
if ! go fmt ./...;
then
  echo -e "${RED}FAILED${SET}"
  exit 1
else
  echo -e "${GREEN}OK${SET}"
fi
