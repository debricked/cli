package main

import (
	"github.com/debricked/cli/pkg/cmd/root"
	"os"
)

func main() {
	if err := root.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
