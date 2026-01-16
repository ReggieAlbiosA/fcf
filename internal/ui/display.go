package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"

	"github.com/ReggieAlbiosA/fcf/internal/platform"
)

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

// Opts holds the global command-line options
var Opts Options

// ColorFuncs holds color functions for output
type ColorFuncs struct {
	Red     func(format string, a ...interface{}) string
	Green   func(format string, a ...interface{}) string
	Yellow  func(format string, a ...interface{}) string
	Blue    func(format string, a ...interface{}) string
	Cyan    func(format string, a ...interface{}) string
	Magenta func(format string, a ...interface{}) string
	Bold    func(format string, a ...interface{}) string
	Dim     func(format string, a ...interface{}) string
}

// Colors is the global instance of ColorFuncs
var Colors ColorFuncs

// InitColors initializes color functions based on terminal support
func InitColors() {
	isTerm := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	if isTerm {
		Colors = ColorFuncs{
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
		Colors = ColorFuncs{
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
func ShowHeader() {
	clearScreen()
	cyan := Colors.Cyan
	bold := Colors.Bold
	fmt.Println(bold(cyan("╔════════════════════════════════════════╗")))
	fmt.Println(bold(cyan("║")) + "   " + bold("fcf") + " - Find File or Folder          " + bold(cyan("║")))
	fmt.Println(bold(cyan("╚════════════════════════════════════════╝")))
	fmt.Println()
}

// showResult displays a single search result with appropriate icon and color
func ShowResult(filePath string, count int) {
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
			Colors.Cyan(fmt.Sprintf("  [%d]", count)),
			Colors.Blue(fmt.Sprintf("📁 %s%c", filePath, filepath.Separator)),
			fileInfo)
	} else if info.Mode()&os.ModeSymlink != 0 {
		// Symlink
		fmt.Printf("%s %s%s\n",
			Colors.Cyan(fmt.Sprintf("  [%d]", count)),
			Colors.Magenta(fmt.Sprintf("🔗 %s", filePath)),
			fileInfo)
	} else if platform.IsExecutable(filePath) {
		// Executable
		fmt.Printf("%s %s%s\n",
			Colors.Cyan(fmt.Sprintf("  [%d]", count)),
			Colors.Green(fmt.Sprintf("⚡ %s", filePath)),
			fileInfo)
	} else {
		// Regular file
		fmt.Printf("%s 📄 %s%s\n",
			Colors.Cyan(fmt.Sprintf("  [%d]", count)),
			filePath,
			fileInfo)
	}
}

// getFileInfo returns formatted file size info if ShowSize is enabled
func getFileInfo(path string, info os.FileInfo) string {
	if !Opts.ShowSize || info.IsDir() {
		return ""
	}
	return Colors.Dim(fmt.Sprintf(" (%s)", FormatSize(info.Size())))
}

// FormatSize formats bytes into human-readable size
func FormatSize(bytes int64) string {
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
func ShowSearchInfo(searchPath, pattern string, usingFd bool) {
	fmt.Println(Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Printf("%s %s\n", Colors.Blue("Searching in:"), Colors.Cyan(searchPath))
	fmt.Printf("%s %s\n", Colors.Blue("Pattern:"), Colors.Yellow(pattern))

	if usingFd {
		fmt.Printf("%s %s\n", Colors.Blue("Method:"), Colors.Green("fd (parallel search)"))
	} else {
		fmt.Printf("%s %s\n", Colors.Blue("Method:"), Colors.Yellow("walk (sequential - install 'fd' for faster search)"))
	}
	fmt.Println(Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	fmt.Println()
}

// showSummary displays search results summary
func ShowSummary(count int, elapsed float64) {
	fmt.Println()
	fmt.Println(Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))

	if count == 0 {
		fmt.Printf("%s for pattern: %s\n", Colors.Yellow("No matches found"), Colors.Cyan(Opts.Pattern))
		fmt.Println()
		fmt.Println(Colors.Dim("Tips:"))
		fmt.Println("  - Try a different pattern")
		fmt.Printf("  - Use %s for case-insensitive search\n", Colors.Cyan("-i"))
	} else {
		fmt.Printf("%s in %s\n",
			Colors.Green(Colors.Bold(fmt.Sprintf("Found %d match(es)", count))),
			Colors.Cyan(fmt.Sprintf("%.2fs", elapsed)))

		if Opts.MaxDisplay > 0 && count > Opts.MaxDisplay {
			fmt.Printf("%s\n", Colors.Yellow(fmt.Sprintf("(Displayed first %d of %d)", Opts.MaxDisplay, count)))
		}
	}
	fmt.Println(Colors.Bold("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
}
