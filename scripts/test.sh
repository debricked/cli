#!/bin/bash

excludedPackages="debricked/cmd/debricked|debricked/pkg/cmd/login|debricked/pkg/cmd/check"
readarray -t packages < <(go list ./... | grep -Ev "$excludedPackages")
go test -cover -coverprofile=coverage.out -v "${packages[@]}"