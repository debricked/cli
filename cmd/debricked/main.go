package main

import (
	"github.com/debricked/cli/pkg/cmd/root"
	"os"
)

var version string // Set at compile time

func main() {
	if err := root.NewRootCmd(version).Execute(); err != nil {
		os.Exit(1)
	}
}
