package main

import (
	"os"

	"github.com/debricked/cli/pkg/cmd/root"
	"github.com/debricked/cli/pkg/wire"
)

var version string // Set at compile time

func main() {
	if err := root.NewRootCmd(version, wire.GetCliContainer()).Execute(); err != nil {
		os.Exit(1)
	}
}
