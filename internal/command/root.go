package command

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ReggieAlbiosA/fcf/internal/install"
	"github.com/ReggieAlbiosA/fcf/internal/navigation"
	"github.com/ReggieAlbiosA/fcf/internal/search"
	"github.com/ReggieAlbiosA/fcf/internal/ui"
)

// Options definition and Opts var removed (moved to ui package)

// timeNow is a wrapper for time.Now (for testing)
var timeNow = time.Now

// Execute parses arguments and runs the appropriate command
func Execute() {
	parseArgs()

	// Show help if requested
	if ui.Opts.Help {
		ShowHelp()
		os.Exit(0)
	}

	// If pattern provided, run single search; otherwise interactive mode
	if ui.Opts.Pattern != "" {
		runSingleSearch()
	} else {
		RunInteractiveMode()
	}
}

func parseArgs() {
	flag.BoolVar(&ui.Opts.Help, "h", false, "Show help message")
	flag.BoolVar(&ui.Opts.Help, "help", false, "Show help message")
	flag.BoolVar(&ui.Opts.IgnoreCase, "i", false, "Case-insensitive pattern matching")
	flag.StringVar(&ui.Opts.Type, "t", "", "Filter by type: 'f' for files, 'd' for directories")
	flag.BoolVar(&ui.Opts.ShowSize, "show-size", false, "Display file sizes")
	flag.IntVar(&ui.Opts.MaxDisplay, "max-display", 0, "Maximum results to display (0 = unlimited)")

	flag.Parse()

	// Get positional arguments
	args := flag.Args()
	if len(args) >= 1 {
		ui.Opts.Pattern = args[0]
	}
	if len(args) >= 2 {
		ui.Opts.Path = args[1]
	} else {
		ui.Opts.Path = "."
	}
}

func runSingleSearch() {
	ui.ShowHeader()

	startTime := getTime()
	result, _ := search.SearchWithStop(ui.Opts.Pattern, ui.Opts.Path)
	elapsed := getTime() - startTime

	ui.ShowSummaryWithStatus(len(result.Results), elapsed, result.Stopped)

	// If results found, offer navigation
	if len(result.Results) > 0 {
		targetPath := SelectResult(result.Results)
		if targetPath != "" {
			fmt.Println()
			navigation.NavigateToPath(targetPath)
		}
	}
}



// RunInstall is called from main for the install command
func RunInstall() {
	install.RunInstall()
}

// RunUninstall is called from main for the uninstall command
func RunUninstall() {
	install.RunUninstall()
}

// RunUpdate is called from main for the update command
func RunUpdate() {
	install.RunUpdate(Version)
}
