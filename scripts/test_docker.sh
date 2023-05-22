#!/usr/bin/env bash

type="$1"

case $type in
  "dev")
    docker build -f build/docker/Dockerfile -t debricked/cli:dev --target dev .
    ;;
  "cli")
    docker build -f build/docker/Dockerfile -t debricked/cli:latest --target cli .
    ;;
  "scan")
    docker build -f build/docker/Dockerfile -t debricked/cli:scan --target scan .
    ;;
  "resolution")
    docker build -f build/docker/Dockerfile -t debricked/cli:resolution --target resolution .
    ;;
  *)
    echo -e "Please use the following type dev, cli, scan. For example ./test_docker.sh dev"
    exit 1
    ;;
esac


