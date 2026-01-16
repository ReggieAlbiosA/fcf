package main

import (
	"os"

	"github.com/ReggieAlbiosA/fcf/internal/command"
	"github.com/ReggieAlbiosA/fcf/internal/navigation"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

func main() {
	// Check for install subcommand first (before flag parsing)
	if len(os.Args) > 1 && os.Args[1] == "install" {
		command.RunInstall()
		return
	}

	// Clean up any stale navigation path on start
	navigation.CleanupNavFile()

	// Initialize colors
	ui.InitColors()

	// Parse command-line arguments and run
	command.Execute()
}
