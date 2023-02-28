package main

import (
	"github.com/debricked/cli/pkg/cmd/root"
	"github.com/debricked/cli/pkg/wire"
	"os"
)

var version string // Set at compile time

func main() {
	if err := root.NewRootCmd(version, wire.GetCliContainer()).Execute(); err != nil {
		os.Exit(1)
	}
}
