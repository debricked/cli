#!/usr/bin/env bash

type="$1"

case $type in
  "dev")
    docker build -f build/docker/Dockerfile -t debricked/cli-dev:latest --target dev .
    ;;
  "cli")
    docker build -f build/docker/Dockerfile -t debricked/cli:latest --target cli .
    ;;
  "scan")
    docker build -f build/docker/Dockerfile -t debricked/cli-scan:latest --target scan .
    ;;
  *)
    echo "${type} type is not supported!"
    exit 1
    ;;
esac


