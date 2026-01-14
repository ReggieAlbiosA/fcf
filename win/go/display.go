package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// Colors holds color functions for output
type Colors struct {
	Red     func(format string, a ...interface{}) string
	Green   func(format string, a ...interface{}) string
	Yellow  func(format string, a ...interface{}) string
	Blue    func(format string, a ...interface{}) string
	Cyan    func(format string, a ...interface{}) string
	Magenta func(format string, a ...interface{}) string
	Bold    func(format string, a ...interface{}) string
	Dim     func(format string, a ...interface{}) string
}

var colors Colors

// initColors initializes color functions based on terminal support
func initColors() {
	isTerm := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	if isTerm {
		colors = Colors{
			Red:     color.New(color.FgRed).SprintfFunc(),
			Green:   color.New(color.FgGreen).SprintfFunc(),
			Yellow:  color.New(color.FgYellow, color.Bold).SprintfFunc(),
			Blue:    color.New(color.FgBlue).SprintfFunc(),
			Cyan:    color.New(color.FgCyan).SprintfFunc(),
			Magenta: color.New(color.FgMagenta).SprintfFunc(),
			Bold:    color.New(color.Bold).SprintfFunc(),
			Dim:     color.New(color.Faint).SprintfFunc(),
		}
	} else {
		// No colors when not a terminal
		noColor := func(format string, a ...interface{}) string {
			return fmt.Sprintf(format, a...)
		}
		colors = Colors{
			Red:     noColor,
			Green:   noColor,
			Yellow:  noColor,
			Blue:    noColor,
			Cyan:    noColor,
			Magenta: noColor,
			Bold:    noColor,
			Dim:     noColor,
		}
	}
}

// clearScreen clears the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// showHeader displays the FCF header
func showHeader() {
	clearScreen()
	cyan := colors.Cyan
	bold := colors.Bold
	fmt.Println(bold(cyan("╔════════════════════════════════════════╗")))
	fmt.Println(bold(cyan("║")) + "   " + bold("fcf") + " - Find File or Folder          " + bold(cyan("║")))
	fmt.Println(bold(cyan("╚════════════════════════════════════════╝")))
	fmt.Println()
}

// showResult displays a single search result with appropriate icon and color
func showResult(filePath string, count int) {
	info, err := os.Lstat(filePath)
	if err != nil {
		fmt.Printf("  [%d] %s\n", count, filePath)
		return
	}

	// Get file info string (size if applicable)
	fileInfo := getFileInfo(filePath, info)

	// Determine file type and display accordingly
	if info.IsDir() {
		// Directory
		fmt.Printf("%s %s%s\n",
			colors.Cyan(fmt.Sprintf("  [%d]", count)),
			colors.Blue(fmt.Sprintf("📁 %s%c", filePath, filepath.Separator)),
			fileInfo)
	} else if info.Mode()&os.ModeSymlink != 0 {
		// Symlink
		fmt.Printf("%s %s%s\n",
			colors.Cyan(fmt.Sprintf("  [%d]", count)),
			colors.Magenta(fmt.Sprintf("🔗 %s", filePath)),
			fileInfo)
	} else if isExecutable(filePath) {
		// Executable
		fmt.Printf("%s %s%s\n",
			colors.Cyan(fmt.Sprintf("  [%d]", count)),
			colors.Green(fmt.Sprintf("⚡ %s", filePath)),
			fileInfo)
	} else {
		// Regular file
		fmt.Printf("%s 📄 %s%s\n",
			colors.Cyan(fmt.Sprintf("  [%d]", count)),
			filePath,
			fileInfo)
	}
}

// getFileInfo returns formatted file size info if ShowSize is enabled
func getFileInfo(path string, info os.FileInfo) string {
	if !opts.ShowSize || info.IsDir() {
		return ""
	}
	return colors.Dim(fmt.Sprintf(" (%s)", formatSize(info.Size())))
}

// formatSize formats bytes into human-readable size
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1fG", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1fM", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1fK", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

// showSearchInfo displays search parameters
func showSearchInfo(searchPath, pattern string, usingFd bool) {
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Printf("%s %s\n", colors.Blue("Searching in:"), colors.Cyan(searchPath))
	fmt.Printf("%s %s\n", colors.Blue("Pattern:"), colors.Yellow(pattern))

	if usingFd {
		fmt.Printf("%s %s\n", colors.Blue("Method:"), colors.Green("fd (parallel search)"))
	} else {
		fmt.Printf("%s %s\n", colors.Blue("Method:"), colors.Yellow("walk (sequential - install 'fd' for faster search)"))
	}
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
}

// showSummary displays search results summary
func showSummary(count int, elapsed float64) {
	fmt.Println()
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))

	if count == 0 {
		fmt.Printf("%s for pattern: %s\n", colors.Yellow("No matches found"), colors.Cyan(opts.Pattern))
		fmt.Println()
		fmt.Println(colors.Dim("Tips:"))
		fmt.Println("  - Try a different pattern")
		fmt.Printf("  - Use %s for case-insensitive search\n", colors.Cyan("-i"))
	} else {
		fmt.Printf("%s in %s\n",
			colors.Green(colors.Bold(fmt.Sprintf("Found %d match(es)", count))),
			colors.Cyan(fmt.Sprintf("%.2fs", elapsed)))

		if opts.MaxDisplay > 0 && count > opts.MaxDisplay {
			fmt.Printf("%s\n", colors.Yellow(fmt.Sprintf("(Displayed first %d of %d)", opts.MaxDisplay, count)))
		}
	}
	fmt.Println(colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
}
