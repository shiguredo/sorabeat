package main

import (
	"os"

	"github.com/shiguredo/sorabeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
