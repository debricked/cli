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
    echo -e "Please use the following type dev, cli, scan. For example ./test_docker.sh dev"
    exit 1
    ;;
esac


