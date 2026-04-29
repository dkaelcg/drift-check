// main.go is the entry point for the drift-check CLI tool.
// It delegates execution to the cmd package root command.
package main

import (
	"os"

	"github.com/yourorg/drift-check/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
