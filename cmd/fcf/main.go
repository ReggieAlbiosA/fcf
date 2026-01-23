package main

import (
	"os"

	"github.com/ReggieAlbiosA/fcf/internal/command"
	"github.com/ReggieAlbiosA/fcf/internal/navigation"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

func main() {
	// Check for subcommands first (before flag parsing)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			command.RunInstall()
			return
		case "uninstall":
			command.RunUninstall()
			return
		case "update":
			command.RunUpdate()
			return
		}
	}

	// Clean up any stale navigation path on start
	navigation.CleanupNavFile()

	// Initialize colors
	ui.InitColors()

	// Parse command-line arguments and run
	command.Execute()
}
