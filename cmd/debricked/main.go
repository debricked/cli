package main

import (
	"debricked/pkg/cmd/root"
	"fmt"
	"os"
)

func main() {
	if err := root.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
