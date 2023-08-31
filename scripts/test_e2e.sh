#!/bin/bash/env

type="$1"

case $type in
  "pip")
    go test ./test/resolve/pip_test.go
    ;;
  "maven")
    go test -v ./test/callgraph/maven_test.go
    ;;
  *)
    go test ./test/...
    ;;
esac
