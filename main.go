package main

import (
	"os"

	"github.com/copyleftdev/specgrade/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
