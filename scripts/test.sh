#!/bin/sh

excludedPackages="debricked/cmd/debricked|debricked/pkg/cmd/login|debricked/pkg/cmd/check"
go test -cover -coverprofile=coverage.out -v $(go list ./... | grep -Ev $excludedPackages)