#!/bin/bash/env

type="$1"

case $type in
  "pip")
    go test ./test/resolve/pip_test.go
    ;;
  *)
    go test ./test/...
    ;;
esac
