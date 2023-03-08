package main

import (
	"fmt"
	"os"

	"cli_example/cli"
)

// Version and CommitHash set in compile time through ldflags
// Will be passed down to cli package
var (
	Version    = ""
	CommitHash = ""
)

func main() {
	// This main function is only responsible for calling the cli package
	// and handling errors returned by the cli package

	// Pass down version and commit hash
	cli.Version = Version
	cli.CommitHash = CommitHash

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
