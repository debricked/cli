#!/usr/bin/env bash
# test if git is installed
if ! command -v git &> /dev/null
then
    echo -e "Failed to find git, thus also the version. Version will be set to v0.0.0"
fi
version=$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
ldFlags="-s -w -X main.version=${version}"
go install -ldflags "${ldFlags}" ./cmd/debricked
go generate -v -x ./cmd/debricked
go build -ldflags "${ldFlags}" ./cmd/debricked
