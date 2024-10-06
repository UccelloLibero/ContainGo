package main

import (
	"fmt"
	"os"
	"ContainGo/cmd"
)

func main() {
	// Execute the root command which includes the 'run' and 'stop' subcommands
	if err := cmd.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1) // Exit with status 1 on failure
	}
}