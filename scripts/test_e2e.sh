#!/bin/bash/env

type="$1"

case $type in
  "resolver")
    go test -timeout 120s ./test/resolve/resolver_test.go  
    ;;
  *)
    go test -timeout 120s ./test/...
    ;;
esac
