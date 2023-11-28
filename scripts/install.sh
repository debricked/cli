#!/usr/bin/env bash
# test if git is installed
if ! command -v git &> /dev/null
then
    echo -e "Failed to find git, thus also the version. Version will be set to v0.0.0"
fi

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
PROJ_DIR="$(realpath "$SCRIPT_DIR/..")"
REMOTE_JSON_URL=https://debricked.com/api/1.0/open/files/supported-formats
LOCAL_JSON_DIR=$PROJ_DIR/internal/file/embedded
LOCAL_JSON_FILE=$PROJ_DIR/internal/file/embedded/supported_formats.json


if [ ! -f "$LOCAL_JSON_FILE" ]; then
    echo "Supported-formats is downloaded from remote for offline backup"
    mkdir -p $LOCAL_JSON_DIR && wget -O $LOCAL_JSON_FILE $REMOTE_JSON_URL
fi

version=$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
ldFlags="-X main.version=${version}"
go install -ldflags "${ldFlags}" $PROJ_DIR/cmd/debricked
