#!/usr/bin/env bash

readarray -t packages < <(go list ./... )
go test -cover -coverprofile=coverage.out "${packages[@]}"