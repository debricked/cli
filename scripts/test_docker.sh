#!/usr/bin/env bash

type="$1"

build_command()
{
    docker build -f build/docker/Dockerfile -t debricked/cli:$1 --target $1 .
    docker build -f build/docker/debian.Dockerfile -t debricked/cli:$1-debian --target $1 .
}

case $type in
  "dev")
    build_command $type
    ;;
  "cli")
    build_command $type
    ;;
  "scan")
    build_command $type
    ;;
  "resolution")
    build_command $type
    ;;
  *)
    echo -e "Please use the following type dev, cli, scan. For example ./test_docker.sh dev"
    exit 1
    ;;
esac


