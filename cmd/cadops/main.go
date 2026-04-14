package main

import (
	"os"

	"github.com/cadops/cadops/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
