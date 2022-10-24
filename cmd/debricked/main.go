package main

import (
	"fmt"
	"github.com/debricked/cli/pkg/cmd/root"
	"os"
)

func main() {
	if err := root.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
