package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

const Version = "2.0.0"

// Options holds the command-line options
type Options struct {
	Pattern    string
	Path       string
	IgnoreCase bool
	Type       string
	ShowSize   bool
	MaxDisplay int
	Help       bool
}

var opts Options

// timeNow is a wrapper for time.Now (for testing)
var timeNow = time.Now

func main() {
	// Check for install subcommand first (before flag parsing)
	if len(os.Args) > 1 && os.Args[1] == "install" {
		runInstall()
		return
	}

	// Clean up any stale navigation path on start
	cleanupNavFile()

	// Initialize colors
	initColors()

	// Parse command-line arguments
	parseArgs()

	// Show help if requested
	if opts.Help {
		showHelp()
		os.Exit(0)
	}

	// If pattern provided, run single search; otherwise interactive mode
	if opts.Pattern != "" {
		runSingleSearch()
	} else {
		runInteractiveMode()
	}
}

func parseArgs() {
	flag.BoolVar(&opts.Help, "h", false, "Show help message")
	flag.BoolVar(&opts.Help, "help", false, "Show help message")
	flag.BoolVar(&opts.IgnoreCase, "i", false, "Case-insensitive pattern matching")
	flag.StringVar(&opts.Type, "t", "", "Filter by type: 'f' for files, 'd' for directories")
	flag.BoolVar(&opts.ShowSize, "show-size", false, "Display file sizes")
	flag.IntVar(&opts.MaxDisplay, "max-display", 0, "Maximum results to display (0 = unlimited)")

	flag.Parse()

	// Get positional arguments
	args := flag.Args()
	if len(args) >= 1 {
		opts.Pattern = args[0]
	}
	if len(args) >= 2 {
		opts.Path = args[1]
	} else {
		opts.Path = "."
	}
}

func runSingleSearch() {
	showHeader()

	startTime := getTime()
	results, _ := search(opts.Pattern, opts.Path)
	elapsed := getTime() - startTime

	showSummary(len(results), elapsed)

	// If results found, offer navigation
	if len(results) > 0 {
		targetPath := selectResult(results)
		if targetPath != "" {
			fmt.Println()
			navigateToPath(targetPath)
		}
	}
}

func showHelp() {
	fmt.Printf("%s v%s\n", colors.Bold("fcf - Find File or Folder"), Version)
	fmt.Println()
	fmt.Println(colors.Bold("USAGE:"))
	fmt.Println("    fcf [OPTIONS] [PATTERN] [PATH]")
	fmt.Println("    fcf                          # Interactive mode")
	fmt.Println("    fcf install                  # Install fcf system-wide")
	fmt.Println()
	fmt.Println(colors.Bold("DESCRIPTION:"))
	fmt.Println("    Interactive tool to find files and folders with pattern matching")
	fmt.Println("    and real-time streaming results. Uses parallel search for speed.")
	fmt.Println()
	fmt.Println(colors.Bold("OPTIONS:"))
	fmt.Printf("    %s               Show this help message\n", colors.Cyan("-h, --help"))
	fmt.Printf("    %s                  Case-insensitive pattern matching\n", colors.Cyan("-i"))
	fmt.Printf("    %s           Filter by type: %s(file) or %s(directory)\n",
		colors.Cyan("-t TYPE"), colors.Yellow("f"), colors.Yellow("d"))
	fmt.Printf("    %s           Display file sizes\n", colors.Cyan("--show-size"))
	fmt.Printf("    %s    Maximum results to display (default: unlimited)\n", colors.Cyan("--max-display NUM"))
	fmt.Println()
	fmt.Println(colors.Bold("COMMANDS:"))
	fmt.Printf("    %s              Install fcf to system (requires sudo/admin)\n", colors.Cyan("install"))
	fmt.Println()
	fmt.Println(colors.Bold("SHELL INTEGRATION (for navigation to work):"))
	showShellIntegrationHelp()
	fmt.Println()
	fmt.Println(colors.Bold("EXAMPLES:"))
	fmt.Printf("    %s\n", colors.Green("# Interactive mode"))
	fmt.Println("    fcf")
	fmt.Println()
	fmt.Printf("    %s\n", colors.Green("# Find all .log files"))
	fmt.Println("    fcf \"*.log\"")
	fmt.Println()
	fmt.Printf("    %s\n", colors.Green("# Find in specific directory"))
	showExamplePath()
	fmt.Println()
	fmt.Printf("    %s\n", colors.Green("# Case-insensitive search for PNG files"))
	fmt.Println("    fcf -i \"*.PNG\"")
	fmt.Println()
	fmt.Printf("    %s\n", colors.Green("# Find only directories named 'src'"))
	fmt.Println("    fcf -t d src")
	fmt.Println()
	fmt.Printf("    %s\n", colors.Green("# Find with file sizes"))
	fmt.Println("    fcf --show-size \"*.mp4\"")
	fmt.Println()
	fmt.Println(colors.Bold("INTERACTIVE WORKFLOW:"))
	fmt.Println("    Step 1: Enter path to search")
	fmt.Println("    Step 2: Enter pattern to find")
	fmt.Println("    Step 3: Navigate to a result path")
	fmt.Println()
	fmt.Println(colors.Bold("NAVIGATION OPTIONS:"))
	fmt.Println("    After navigation, choose:")
	fmt.Printf("    %s - Find again (restart from Step 1)\n", colors.Cyan("f"))
	fmt.Printf("    %s - Repeat search (go to Step 2, same path)\n", colors.Cyan("r"))
	fmt.Printf("    %s - Exit\n", colors.Cyan("n"))
	fmt.Println()
	fmt.Println(colors.Bold("PERFORMANCE:"))
	fmt.Println("    - Uses 'fd' for fast parallel searching (if installed)")
	fmt.Println("    - Falls back to Go's filepath.WalkDir if fd is not available")
	fmt.Printf("    - Install fd: %s\n", colors.Cyan(getFdInstallHint()))
	fmt.Println()
}

// showShellIntegrationHelp displays platform-specific shell integration instructions
func showShellIntegrationHelp() {
	if runtime.GOOS == "windows" {
		fmt.Printf("    Add this to your %s:\n", colors.Cyan("$PROFILE"))
		fmt.Println()
		fmt.Println(colors.Yellow("    function fcf {"))
		fmt.Println(colors.Yellow("        $navFile = Join-Path $env:TEMP \"fcf_nav_path\""))
		fmt.Println(colors.Yellow("        if (Test-Path $navFile) { Remove-Item $navFile -Force }"))
		fmt.Println(colors.Yellow("        & \"C:\\Program Files\\fcf\\fcf.exe\" @args"))
		fmt.Println(colors.Yellow("        if (Test-Path $navFile) {"))
		fmt.Println(colors.Yellow("            $target = Get-Content $navFile -Raw"))
		fmt.Println(colors.Yellow("            Remove-Item $navFile -Force"))
		fmt.Println(colors.Yellow("            if (Test-Path $target -PathType Container) {"))
		fmt.Println(colors.Yellow("                Set-Location $target"))
		fmt.Println(colors.Yellow("            }"))
		fmt.Println(colors.Yellow("        }"))
		fmt.Println(colors.Yellow("    }"))
	} else {
		fmt.Printf("    Add this to your %s or %s:\n", colors.Cyan("~/.bashrc"), colors.Cyan("~/.zshrc"))
		fmt.Println()
		fmt.Println(colors.Yellow("    fcf() {"))
		fmt.Println(colors.Yellow("        local nav_file=\"/tmp/fcf_nav_path\""))
		fmt.Println(colors.Yellow("        rm -f \"$nav_file\""))
		fmt.Println(colors.Yellow("        command fcf \"$@\""))
		fmt.Println(colors.Yellow("        if [[ -f \"$nav_file\" ]]; then"))
		fmt.Println(colors.Yellow("            local target"))
		fmt.Println(colors.Yellow("            target=$(cat \"$nav_file\")"))
		fmt.Println(colors.Yellow("            rm -f \"$nav_file\""))
		fmt.Println(colors.Yellow("            if [[ -d \"$target\" ]]; then"))
		fmt.Println(colors.Yellow("                cd \"$target\" || return"))
		fmt.Println(colors.Yellow("            fi"))
		fmt.Println(colors.Yellow("        fi"))
		fmt.Println(colors.Yellow("    }"))
	}
}

// showExamplePath displays a platform-appropriate example path
func showExamplePath() {
	if runtime.GOOS == "windows" {
		fmt.Println("    fcf \"*.js\" C:\\Projects")
	} else {
		fmt.Println("    fcf \"*.js\" ~/projects")
	}
}
